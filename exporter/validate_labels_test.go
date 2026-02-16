package exporter

import (
	"strings"
	"testing"
)

func TestValidateConstLabelConflicts_QueryLabelOverlap(t *testing.T) {
	cfg := `
q1:
  query: SELECT 1 AS value, 'x' AS datname
  metrics:
    - datname:
        usage: LABEL
        rename: db
        description: database
    - value:
        usage: GAUGE
        description: value
`
	queries, err := ParseConfig([]byte(cfg))
	if err != nil {
		t.Fatalf("ParseConfig failed: %v", err)
	}

	labels := parseConstLabels("db=foo")
	if err := validateConstLabelConflicts(labels, queries, false); err == nil {
		t.Fatal("expected const label conflict error, got nil")
	}
}

func TestValidateConstLabelConflicts_InternalReservedLabels(t *testing.T) {
	labels := parseConstLabels("datname=foo")
	if err := validateConstLabelConflicts(labels, nil, false); err == nil {
		t.Fatal("expected reserved label conflict with intro metrics enabled, got nil")
	}
	// When intro metrics are disabled, internal dynamic series are not emitted.
	if err := validateConstLabelConflicts(labels, nil, true); err != nil {
		t.Fatalf("expected no error with intro metrics disabled, got %v", err)
	}
}

func TestNewExporterRejectsConstLabelConflict(t *testing.T) {
	cfg := `
q1:
  query: SELECT 1 AS value, 'x' AS datname
  metrics:
    - datname:
        usage: LABEL
        rename: db
    - value:
        usage: GAUGE
`
	_, err := NewExporter(
		"postgresql://u:p@localhost:5432/postgres?sslmode=disable",
		WithConfigReader(strings.NewReader(cfg)),
		WithConstLabels("db=foo"),
	)
	if err == nil {
		t.Fatal("expected NewExporter to fail on const label conflict, got nil")
	}
}
