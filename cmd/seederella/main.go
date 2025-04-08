package main

import (
	"fmt"
	"log"

	"github.com/QaidVoid/seederella/internal/config"
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

