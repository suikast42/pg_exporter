package exporter

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"testing"
)

func parseConfigDirLikeMerge(t *testing.T, dir string) map[string]*Query {
	t.Helper()

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir(%s) failed: %v", dir, err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := filepath.Ext(entry.Name())
		if ext != ".yml" && ext != ".yaml" {
			continue
		}
		names = append(names, entry.Name())
	}
	slices.Sort(names)

	queries := make(map[string]*Query)
	for _, name := range names {
		path := filepath.Join(dir, name)
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile(%s) failed: %v", path, err)
		}
		parsed, err := ParseConfig(content)
		if err != nil {
			t.Fatalf("ParseConfig(%s) failed: %v", path, err)
		}
		for branch, q := range parsed {
			queries[branch] = q
		}
	}
	return queries
}

func TestMergedConfigsMatchSplitDirectories(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), ".."))

	cases := []struct {
		name   string
		dir    string
		merged string
	}{
		{
			name:   "current",
			dir:    filepath.Join(repoRoot, "config"),
			merged: filepath.Join(repoRoot, "pg_exporter.yml"),
		},
		{
			name:   "legacy",
			dir:    filepath.Join(repoRoot, "legacy", "config"),
			merged: filepath.Join(repoRoot, "legacy", "pg_exporter.yml"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			splitQueries := parseConfigDirLikeMerge(t, tc.dir)

			mergedContent, err := os.ReadFile(tc.merged)
			if err != nil {
				t.Fatalf("ReadFile(%s) failed: %v", tc.merged, err)
			}
			mergedQueries, err := ParseConfig(mergedContent)
			if err != nil {
				t.Fatalf("ParseConfig(%s) failed: %v", tc.merged, err)
			}

			if len(splitQueries) != len(mergedQueries) {
				t.Fatalf("query count mismatch: split=%d merged=%d", len(splitQueries), len(mergedQueries))
			}

			for branch, splitQuery := range splitQueries {
				mergedQuery, ok := mergedQueries[branch]
				if !ok {
					t.Fatalf("branch %q missing from merged config %s", branch, tc.merged)
				}
				if !reflect.DeepEqual(splitQuery, mergedQuery) {
					t.Fatalf("branch %q differs between split dir %s and merged config %s", branch, tc.dir, tc.merged)
				}
			}
		})
	}
}
