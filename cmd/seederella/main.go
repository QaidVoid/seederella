package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/QaidVoid/seederella/internal/config"
	"github.com/QaidVoid/seederella/internal/db"
	"github.com/QaidVoid/seederella/internal/faker"
	"github.com/spf13/cobra"
)

var configPath string

var rootCmd = &cobra.Command{
	Use:   "seederella",
	Short: "seederella - a database seeding tool that brings your schema to life effortlessly and efficiently",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
		fmt.Printf("✅ Loaded config for tables: %v\n", keys(cfg.Tables))

		sqlDB, err := db.Connect(cfg.Driver, cfg.DSN, cfg.Schema)
		if err != nil {
			log.Fatalf("❌ Failed to connect to DB: %v", err)
		}
		defer sqlDB.Close()
		fmt.Println("✅ Connected to DB")

		for tableName, table := range cfg.Tables {
			fmt.Printf("🚀 Seeding %s (%d rows)\n", tableName, table.Count)

			for i := 0; i < table.Count; i++ {
				fields := []string{}
				values := []any{}
				placeholders := []string{}

				idx := 1
				for colName, colCfg := range table.Fields {
					var val any
					if colCfg.Value != nil {
						val = colCfg.Value
					} else if colCfg.Faker != "" {
						val, err = faker.Generate(colCfg.Faker)
						if err != nil {
							log.Fatalf("faker error: %v", err)
						}
					}
					fields = append(fields, `"`+colName+`"`)
					values = append(values, val)
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
					log.Fatalf("failed to insert into %s: %v", tableName, err)
				}
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
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "config.yaml", "Path to config file")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
