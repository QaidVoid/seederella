package config_test

import (
	"testing"
	"os"

	"github.com/QaidVoid/seederella/internal/config"
)

func TestLoadConfig(t *testing.T) {
	yamlContent := `
tables:
  users:
    count: 2
    fields:
      name:
        faker: name
      email:
        faker: email
`

	tmp := "test_config.yaml"
	if err := os.WriteFile(tmp, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp)

	cfg, err := config.LoadConfig(tmp)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(cfg.Tables) != 1 {
		t.Errorf("Expected 1 table, got %d", len(cfg.Tables))
	}

	users, ok := cfg.Tables["users"]
	if !ok {
		t.Fatal("Expected users table")
	}
	if users.Count != 2 {
		t.Errorf("Expected count=2, got %d", users.Count)
	}
}
