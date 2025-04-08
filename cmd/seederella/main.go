package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/QaidVoid/seederella/internal/config"
	"github.com/QaidVoid/seederella/internal/db"
	"github.com/QaidVoid/seederella/internal/faker"
	"github.com/QaidVoid/seederella/internal/resolve"
	"github.com/spf13/cobra"
)

var (
	configPath     string
	cleanDb        bool
	overrideDriver string
	overrideDSN    string
)

var rootCmd = &cobra.Command{
	Use:   "seederella",
	Short: "seederella - a database seeding tool that brings your schema to life effortlessly and efficiently",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		if overrideDriver != "" {
			cfg.Driver = overrideDriver
		}
		if overrideDSN != "" {
			cfg.DSN = overrideDSN
		}

		if cfg.Driver == "" || cfg.DSN == "" {
			log.Fatal("❌ Config must include 'driver' and 'dsn'")
		}
		fmt.Printf("✅ Loaded config for tables: %v\n", keys(cfg.Tables))

		sqlDB, err := db.Connect(cfg.Driver, cfg.DSN, cfg.Schema)
		if err != nil {
			log.Fatalf("❌ Failed to connect to DB: %v", err)
		}
		defer sqlDB.Close()
		fmt.Println("✅ Connected to DB")

		if cleanDb {
			if err := db.Clean(cfg.Driver, sqlDB, cfg.Schema); err != nil {
				log.Fatalf("failed to clean DB: %v", err)
			}
		}

		inserted := make(map[string][]map[string]any)

		for tableName, table := range cfg.Tables {
			fmt.Printf("🚀 Seeding %s (%d rows)\n", tableName, table.Count)

			colOrder, err := resolve.SortFields(table.Fields)
			if err != nil {
				log.Fatalf("❌ Field sort error for table %s: %v", tableName, err)
			}
			for i := 0; i < table.Count; i++ {
				fields := []string{}
				values := []any{}
				placeholders := []string{}
				rowData := make(map[string]any)

				idx := 1
				for _, colName := range colOrder {
					colCfg := table.Fields[colName]
					val, err := resolve.ResolveField(colName, colCfg, rowData, inserted)
					if err != nil {
						log.Fatalf("❌ Field resolution failed (%s.%s): %v", tableName, colName, err)
					}
					fields = append(fields, `"`+colName+`"`)
					values = append(values, val)
					rowData[colName] = val

					// PostgreSQL uses $1, $2, ... | MySQL and SQLite just use ?
					if cfg.Driver == "postgres" {
						placeholders = append(placeholders, fmt.Sprintf("$%d", idx))
					} else {
						placeholders = append(placeholders, "?")
					}
					idx++
				}

				query := fmt.Sprintf(
					"INSERT INTO \"%s\" (%s) VALUES (%s)",
					tableName,
					strings.Join(fields, ", "),
					strings.Join(placeholders, ", "),
				)

				_, err := sqlDB.Exec(query, values...)
				if err != nil {
					log.Fatalf("❌ Failed to insert into %s: %v", tableName, err)
				}

				inserted[tableName] = append(inserted[tableName], rowData)
			}

			fmt.Printf("✅ Done: %s\n", tableName)
		}
	},
}

func keys[T any](m map[string]T) []string {
	k := make([]string, 0, len(m))
	for key := range m {
		k = append(k, key)
	}
	return k
}

func init() {
	faker.Init()
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "Path to config file")
	rootCmd.PersistentFlags().BoolVar(&cleanDb, "clean", false, "Drop and recreate DB schema before seeding")
	rootCmd.PersistentFlags().StringVar(&overrideDriver, "driver", "", "Override DB driver")
	rootCmd.PersistentFlags().StringVar(&overrideDSN, "dsn", "", "Override DB DSN")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
