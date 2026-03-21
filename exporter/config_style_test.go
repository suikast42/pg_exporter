package exporter

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

var inlineMetricDescriptionRE = regexp.MustCompile(`^\s*-\s*[^:]+:\s*\{.*\bdescription:\s*(.+)\}\s*$`)

func TestInlineMetricDescriptionsUseDoubleQuotes(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), ".."))

	for _, rel := range []string{"config", filepath.Join("legacy", "config")} {
		dir := filepath.Join(repoRoot, rel)
		t.Run(rel, func(t *testing.T) {
			entries, err := os.ReadDir(dir)
			if err != nil {
				t.Fatalf("ReadDir(%s) failed: %v", dir, err)
			}
			for _, entry := range entries {
				if entry.IsDir() || filepath.Ext(entry.Name()) != ".yml" {
					continue
				}
				path := filepath.Join(dir, entry.Name())
				f, err := os.Open(path)
				if err != nil {
					t.Fatalf("Open(%s) failed: %v", path, err)
				}

				scanner := bufio.NewScanner(f)
				for lineNo := 1; scanner.Scan(); lineNo++ {
					line := scanner.Text()
					m := inlineMetricDescriptionRE.FindStringSubmatch(line)
					if m == nil {
						continue
					}
					desc := strings.TrimSpace(m[1])
					if len(desc) < 2 || desc[0] != '"' || desc[len(desc)-1] != '"' {
						t.Errorf("%s:%d inline metric description must use double quotes: %s", path, lineNo, desc)
					}
				}
				if err := scanner.Err(); err != nil {
					_ = f.Close()
					t.Fatalf("Scan(%s) failed: %v", path, err)
				}
				if err := f.Close(); err != nil {
					t.Fatalf("Close(%s) failed: %v", path, err)
				}
			}
		})
	}
}

func TestLegacySplitConfigsEndWithTwoBlankLines(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "legacy", "config"))

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir(%s) failed: %v", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yml" {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile(%s) failed: %v", path, err)
		}
		trailingNewlines := len(data) - len(bytes.TrimRight(data, "\n"))
		if trailingNewlines != 3 {
			t.Errorf("%s must end with exactly two blank lines (3 trailing newlines), got %d", path, trailingNewlines)
		}
	}
}
