package exporter

import (
	"log/slog"
	"runtime"
	"sync"
	"sync/atomic"
)

/* ================ Parameters ================ */

// Version is read by make build procedure
var Version = "1.2.0"

// Build information. Populated at build-time.
var (
	Branch    = "main"
	Revision  = "HEAD"
	BuildDate = "20250421212100" // will be overwritten during release
	GoVersion = runtime.Version()
	GOOS      = runtime.GOOS
	GOARCH    = runtime.GOARCH
)

var defaultPGURL = "postgresql:///?sslmode=disable"

/* ================ Global Vars ================ */

// PgExporter is the global singleton of Exporter
var (
	PgExporter        *Exporter
	currentExporterPt atomic.Pointer[Exporter]
	ReloadLock        sync.RWMutex
	Logger            *slog.Logger
)

func setCurrentExporter(e *Exporter) {
	PgExporter = e
	currentExporterPt.Store(e)
}
