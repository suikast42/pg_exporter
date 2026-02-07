package exporter

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// Ensure the repo-bundled config/ covers PG10..PG18 without version gaps for
// collectors that are supposed to work on PG10+. This is a cheap static check
// (no DB required) to catch off-by-one mistakes on min/max_version splits.
func TestConfigCoveragePG10To18(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	configDir := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "config"))
	if _, err := os.Stat(configDir); err != nil {
		t.Skipf("config dir not found: %s: %v", configDir, err)
	}

	queries, err := LoadConfig(configDir)
	if err != nil {
		t.Fatalf("LoadConfig(%s) failed: %v", configDir, err)
	}

	byName := make(map[string][]*Query)
	for _, q := range queries {
		if q.HasTag("pgbouncer") { // PG and pgbouncer versions are in different namespaces.
			continue
		}
		byName[q.Name] = append(byName[q.Name], q)
	}

	for name, qs := range byName {
		minMin := 0
		for i, q := range qs {
			if i == 0 || q.MinVersion < minMin {
				minMin = q.MinVersion
			}
		}

		// Collectors introduced after PG10 are allowed to have gaps for PG10-.
		if minMin > 100000 {
			continue
		}

		for v := 100000; v <= 180000; v += 10000 { // PG10..PG18
			var appl []*Query
			for _, q := range qs {
				if q.MinVersion != 0 && v < q.MinVersion {
					continue
				}
				if q.MaxVersion != 0 && v >= q.MaxVersion { // exclude
					continue
				}
				appl = append(appl, q)
			}

			if len(appl) == 0 {
				t.Errorf("collector %q has no branch for server_version_num=%d", name, v)
				continue
			}

			// Multiple branches for the same Name are only acceptable when they are
			// mutually exclusive via tags (e.g. primary vs replica).
			if len(appl) > 1 {
				if name == "pg" && len(appl) == 2 &&
					((appl[0].HasTag("primary") && appl[1].HasTag("replica")) ||
						(appl[0].HasTag("replica") && appl[1].HasTag("primary"))) {
					continue
				}
				t.Errorf("collector %q has %d overlapping branches for server_version_num=%d: %v", name, len(appl), v, func() []string {
					out := make([]string, 0, len(appl))
					for _, q := range appl {
						out = append(out, q.Branch)
					}
					return out
				}())
			}
		}
	}
}

