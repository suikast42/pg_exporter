package exporter

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"testing"
)

var errProbeHealthTestPingCalled = errors.New("ping called")

// probeHealthTestDriver is a tiny database/sql driver used to verify that
// ProbeHealth in pgbouncer mode does not use db.PingContext (lib/pq Ping uses
// a ";" query which PgBouncer rejects). If Ping is called we return a sentinel
// error to fail the test.
type probeHealthTestDriver struct{}

func (d probeHealthTestDriver) Open(name string) (driver.Conn, error) {
	return &probeHealthTestConn{}, nil
}

type probeHealthTestConn struct{}

func (c *probeHealthTestConn) Prepare(query string) (driver.Stmt, error) {
	return nil, errors.New("prepare not supported")
}
func (c *probeHealthTestConn) Close() error              { return nil }
func (c *probeHealthTestConn) Begin() (driver.Tx, error) { return nil, errors.New("tx not supported") }

func (c *probeHealthTestConn) Ping(ctx context.Context) error { return errProbeHealthTestPingCalled }

func (c *probeHealthTestConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	// Return an empty resultset; database/sql will surface this as sql.ErrNoRows
	// to QueryRowContext, which ProbeHealth should treat as a successful probe
	// for PgBouncer (SHOW VERSION may return via NOTICE only).
	return &probeHealthTestRows{}, nil
}

type probeHealthTestRows struct{}

func (r *probeHealthTestRows) Columns() []string { return []string{"version"} }
func (r *probeHealthTestRows) Close() error      { return nil }
func (r *probeHealthTestRows) Next(dest []driver.Value) error {
	return io.EOF
}

func init() {
	sql.Register("probehealth_test", probeHealthTestDriver{})
}

func TestProbeHealthPgbouncerDoesNotPingAndTreatsNoRowsAsUp(t *testing.T) {
	db, err := sql.Open("probehealth_test", "")
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	defer db.Close()

	s := &Server{
		DB:             db,
		PgbouncerMode:  true,
		ConnectTimeout: 500,
	}

	up, recovery, starting, err := s.ProbeHealth()
	if err != nil {
		t.Fatalf("ProbeHealth error = %v", err)
	}
	if !up {
		t.Fatalf("up = %v, want true", up)
	}
	if recovery {
		t.Fatalf("recovery = %v, want false", recovery)
	}
	if starting {
		t.Fatalf("starting = %v, want false", starting)
	}
}
