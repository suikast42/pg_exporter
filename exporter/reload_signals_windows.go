//go:build windows

package exporter

import (
	"os"
	"syscall"
)

var reloadSignals = []os.Signal{syscall.SIGHUP}
