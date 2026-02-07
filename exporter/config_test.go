package exporter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseConfigUsageCaseInsensitive(t *testing.T) {
	config := `
test_query:
  query: SELECT 1 AS metric, 'db' AS datname
  metrics:
    - metric:
        usage: gauge
        description: metric value
    - datname:
        usage: label
        description: database name
`

	queries, err := ParseConfig([]byte(config))
	if err != nil {
		t.Fatalf("ParseConfig returned error: %v", err)
	}

	query, ok := queries["test_query"]
	if !ok {
		t.Fatalf("query test_query not found")
	}

	if got := query.Columns["metric"].Usage; got != GAUGE {
		t.Fatalf("metric usage = %s, want %s", got, GAUGE)
	}
	if got := query.Columns["datname"].Usage; got != LABEL {
		t.Fatalf("datname usage = %s, want %s", got, LABEL)
	}
}

func TestParseConfigInvalidUsage(t *testing.T) {
	config := `
bad_query:
  query: SELECT 1 AS metric
  metrics:
    - metric:
        usage: bad_usage
        description: metric value
`

	if _, err := ParseConfig([]byte(config)); err == nil {
		t.Fatal("ParseConfig should fail on unsupported usage")
	}
}

func TestParseQueryErrors(t *testing.T) {
	if _, err := ParseQuery(`{}`); err == nil {
		t.Fatal("ParseQuery should fail when no query is defined")
	}

	multi := `
q1:
  query: SELECT 1 AS metric
  metrics:
    - metric:
        usage: gauge
q2:
  query: SELECT 2 AS metric
  metrics:
    - metric:
        usage: gauge
`
	if _, err := ParseQuery(multi); err == nil {
		t.Fatal("ParseQuery should fail when multiple queries are defined")
	}
}

func TestLoadConfigDirectoryPriorityAndOverride(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "0100-a.yml")
	f2 := filepath.Join(dir, "0200-b.yml")

	cfg1 := `
q_common:
  query: SELECT 1 AS metric
  metrics:
    - metric:
        usage: gauge
`
	cfg2 := `
q_common:
  query: SELECT 2 AS metric
  metrics:
    - metric:
        usage: gauge
q_extra:
  query: SELECT 3 AS metric
  metrics:
    - metric:
        usage: gauge
`
	if err := os.WriteFile(f1, []byte(cfg1), 0o644); err != nil {
		t.Fatalf("write config 1 failed: %v", err)
	}
	if err := os.WriteFile(f2, []byte(cfg2), 0o644); err != nil {
		t.Fatalf("write config 2 failed: %v", err)
	}

	queries, err := LoadConfig(dir)
	if err != nil {
		t.Fatalf("LoadConfig dir failed: %v", err)
	}
	if len(queries) != 2 {
		t.Fatalf("LoadConfig query count = %d, want 2", len(queries))
	}
	if queries["q_common"].SQL != "SELECT 2 AS metric" {
		t.Fatalf("q_common should be overridden by later file, got: %s", queries["q_common"].SQL)
	}
	// 2nd config file gets default priority 102.
	if queries["q_common"].Priority != 102 {
		t.Fatalf("q_common priority = %d, want 102", queries["q_common"].Priority)
	}
	if queries["q_extra"].Priority != 102 {
		t.Fatalf("q_extra priority = %d, want 102", queries["q_extra"].Priority)
	}
}

func TestGetConfigPrecedence(t *testing.T) {
	originConfigPath := *configPath
	t.Cleanup(func() { *configPath = originConfigPath })

	originEnv := os.Getenv("PG_EXPORTER_CONFIG")
	t.Cleanup(func() { _ = os.Setenv("PG_EXPORTER_CONFIG", originEnv) })

	*configPath = "/tmp/from-cli.yml"
	_ = os.Setenv("PG_EXPORTER_CONFIG", "/tmp/from-env.yml")
	if got := GetConfig(); got != "/tmp/from-cli.yml" {
		t.Fatalf("GetConfig CLI precedence failed: got %s", got)
	}

	*configPath = ""
	if got := GetConfig(); got != "/tmp/from-env.yml" {
		t.Fatalf("GetConfig env fallback failed: got %s", got)
	}
}
