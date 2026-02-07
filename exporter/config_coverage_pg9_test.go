package exporter

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// Ensure the legacy config/pg9/ covers PG9.0..PG9.6 without version gaps for
// collectors that are supposed to work on legacy PG9.x.
func TestConfigCoveragePG9(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	configDir := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "config", "pg9"))
	if _, err := os.Stat(configDir); err != nil {
		t.Skipf("legacy config dir not found: %s: %v", configDir, err)
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

	versions := []int{90000, 90100, 90200, 90300, 90400, 90500, 90600} // PG9.0..PG9.6
	for name, qs := range byName {
		minMin := 0
		for i, q := range qs {
			if i == 0 || q.MinVersion < minMin {
				minMin = q.MinVersion
			}
		}

		for _, v := range versions {
			// Collectors introduced after v are allowed to have gaps for older versions.
			if minMin != 0 && v < minMin {
				continue
			}

			var appl []*Query
			for _, q := range qs {
				if q.MinVersion != 0 && v < q.MinVersion {
					continue
				}
				if q.MaxVersion != 0 && v >= q.MaxVersion {
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
