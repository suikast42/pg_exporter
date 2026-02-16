package exporter

import (
	"errors"
	"fmt"
	"testing"

	"github.com/lib/pq"
)

func TestIsPostgresStartupError(t *testing.T) {
	if !isPostgresStartupError(&pq.Error{Code: pq.ErrorCode(pgSQLStateCannotConnectNow)}) {
		t.Fatal("expected SQLSTATE 57P03 to be recognized as startup error")
	}

	wrapped := fmt.Errorf("wrapped: %w", &pq.Error{Code: pq.ErrorCode(pgSQLStateCannotConnectNow)})
	if !isPostgresStartupError(wrapped) {
		t.Fatal("expected wrapped SQLSTATE 57P03 to be recognized as startup error")
	}

	if isPostgresStartupError(errors.New("plain error")) {
		t.Fatal("plain error should not be recognized as startup error")
	}
	if isPostgresStartupError(&pq.Error{Code: "08006"}) {
		t.Fatal("non-57P03 postgres error should not be recognized as startup error")
	}
}

func TestUpdateHealthStateWithStartup(t *testing.T) {
	e := &Exporter{}

	e.updateHealthStateWithStartup(false, false, true)
	if e.Up() {
		t.Fatal("startup state should not be considered up")
	}
	if e.Recovery() {
		t.Fatal("startup state should not expose recovery=true")
	}
	if got := e.Status(); got != "starting" {
		t.Fatalf("status = %s, want starting", got)
	}

	e.updateHealthStateWithStartup(false, false, false)
	if got := e.Status(); got != "down" {
		t.Fatalf("status = %s, want down", got)
	}

	e.updateHealthStateWithStartup(true, true, false)
	if got := e.Status(); got != "replica" {
		t.Fatalf("status = %s, want replica", got)
	}

	e.updateHealthStateWithStartup(true, false, false)
	if got := e.Status(); got != "primary" {
		t.Fatalf("status = %s, want primary", got)
	}
}
