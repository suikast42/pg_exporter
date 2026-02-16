package exporter

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	Logger = configureLogger("error", "logfmt")
	os.Exit(m.Run())
}
