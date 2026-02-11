package exporter

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestProcessPGURLKeepsEncodedQueryValues(t *testing.T) {
	input := "postgresql://user:pass@localhost:5432/postgres?application_name=a%26b&password=p%3Dq"
	output := ProcessPGURL(input)
	if output == "" {
		t.Fatalf("ProcessPGURL returned empty output")
	}

	parsed, err := url.Parse(output)
	if err != nil {
		t.Fatalf("failed to parse output URL: %v", err)
	}
	qs := parsed.Query()
	if got := qs.Get("application_name"); got != "a&b" {
		t.Fatalf("application_name = %q, want %q", got, "a&b")
	}
	if got := qs.Get("password"); got != "p=q" {
		t.Fatalf("password = %q, want %q", got, "p=q")
	}
	if got := qs.Get("sslmode"); got != "disable" {
		t.Fatalf("sslmode = %q, want %q", got, "disable")
	}
}

func TestShadowPGURLRedactsQueryPassword(t *testing.T) {
	input := "postgresql://user:pass@localhost:5432/postgres?password=p%26q%3D1&application_name=test"
	output := ShadowPGURL(input)

	parsed, err := url.Parse(output)
	if err != nil {
		t.Fatalf("failed to parse redacted URL: %v", err)
	}
	if got := parsed.Query().Get("password"); got != "xxxxx" {
		t.Fatalf("password = %q, want %q", got, "xxxxx")
	}
}

func TestParseDatnameAndReplaceDatname(t *testing.T) {
	src := "postgresql://user:pass@localhost:5432/postgres?sslmode=disable"
	if got := ParseDatname(src); got != "postgres" {
		t.Fatalf("ParseDatname = %q, want %q", got, "postgres")
	}

	replaced := ReplaceDatname(src, "otherdb")
	if got := ParseDatname(replaced); got != "otherdb" {
		t.Fatalf("ParseDatname(replaced) = %q, want %q", got, "otherdb")
	}

	srcWithDbname := "postgresql://user:pass@localhost:5432?sslmode=disable&dbname=pgbouncer"
	if got := ParseDatname(srcWithDbname); got != "pgbouncer" {
		t.Fatalf("ParseDatname(dbname=) = %q, want %q", got, "pgbouncer")
	}

	replacedDbname := ReplaceDatname(srcWithDbname, "postgres")
	if got := ParseDatname(replacedDbname); got != "postgres" {
		t.Fatalf("ParseDatname(replaced dbname=) = %q, want %q", got, "postgres")
	}
}

func TestRetrievePGURLPriority(t *testing.T) {
	originPGURL := *pgURL
	*pgURL = ""
	t.Cleanup(func() { *pgURL = originPGURL })

	originExporterURL := os.Getenv("PG_EXPORTER_URL")
	originPGURLenv := os.Getenv("PGURL")
	originFile := os.Getenv("PG_EXPORTER_URL_FILE")
	t.Cleanup(func() {
		_ = os.Setenv("PG_EXPORTER_URL", originExporterURL)
		_ = os.Setenv("PGURL", originPGURLenv)
		_ = os.Setenv("PG_EXPORTER_URL_FILE", originFile)
	})

	_ = os.Setenv("PG_EXPORTER_URL", "postgresql://env-user:env-pass@localhost:5432/envdb")
	_ = os.Setenv("PGURL", "postgresql://pgurl-user:pgurl-pass@localhost:5432/pgurldb")
	*pgURL = "postgresql://cli-user:cli-pass@localhost:5432/clidb"
	if got := RetrievePGURL(); got != *pgURL {
		t.Fatalf("RetrievePGURL CLI precedence failed: got %s", got)
	}

	*pgURL = ""
	if got := RetrievePGURL(); got != os.Getenv("PG_EXPORTER_URL") {
		t.Fatalf("RetrievePGURL env precedence failed: got %s", got)
	}

	_ = os.Unsetenv("PG_EXPORTER_URL")
	if got := RetrievePGURL(); got != os.Getenv("PGURL") {
		t.Fatalf("RetrievePGURL PGURL fallback failed: got %s", got)
	}

	_ = os.Unsetenv("PGURL")
	file := filepath.Join(t.TempDir(), "dsn.txt")
	fileURL := "postgresql://file-user:file-pass@localhost:5432/filedb"
	if err := os.WriteFile(file, []byte(fileURL), 0o644); err != nil {
		t.Fatalf("write dsn file failed: %v", err)
	}
	_ = os.Setenv("PG_EXPORTER_URL_FILE", file)
	if got := RetrievePGURL(); got != fileURL {
		t.Fatalf("RetrievePGURL file fallback failed: got %s", got)
	}

	_ = os.Unsetenv("PG_EXPORTER_URL_FILE")
	if got := RetrievePGURL(); got != defaultPGURL {
		t.Fatalf("RetrievePGURL default fallback failed: got %s", got)
	}
}
