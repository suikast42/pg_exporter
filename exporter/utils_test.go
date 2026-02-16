package exporter

import (
	"math"
	"reflect"
	"testing"
	"time"
)

func TestParseCSV(t *testing.T) {
	if got := parseCSV(""); got != nil {
		t.Fatalf("parseCSV empty = %v, want nil", got)
	}

	got := parseCSV(" a, b,, c , ")
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseCSV result = %#v, want %#v", got, want)
	}
}

func TestParseConstLabels(t *testing.T) {
	if got := parseConstLabels(""); got != nil {
		t.Fatalf("parseConstLabels empty = %v, want nil", got)
	}

	labels := parseConstLabels("env=prod,region=us-east-1")
	if labels["env"] != "prod" || labels["region"] != "us-east-1" {
		t.Fatalf("parseConstLabels valid result = %v", labels)
	}

	labels = parseConstLabels("token=a=b=c")
	if labels["token"] != "a=b=c" {
		t.Fatalf("parseConstLabels should preserve '=' in value, got %q", labels["token"])
	}

	labels = parseConstLabels("bad,noeq=,=noval,ok=1")
	if len(labels) != 1 || labels["ok"] != "1" {
		t.Fatalf("parseConstLabels malformed handling = %v", labels)
	}
}

func TestCastFloat64(t *testing.T) {
	now := time.Unix(1700000000, 0)

	withScale := func(scale, def string) *Column {
		c := &Column{Scale: scale, Default: def}
		if err := c.parseNumbers(); err != nil {
			t.Fatalf("parseNumbers failed for scale=%q default=%q: %v", scale, def, err)
		}
		return c
	}

	if got := castFloat64(int64(3), withScale("2", "")); got != 6 {
		t.Fatalf("int64 scale cast = %v, want 6", got)
	}
	if got := castFloat64(float64(1.5), withScale("2", "")); got != 3 {
		t.Fatalf("float64 scale cast = %v, want 3", got)
	}
	if got := castFloat64(now, nil); got != float64(now.Unix()) {
		t.Fatalf("time cast = %v, want %v", got, now.Unix())
	}
	if got := castFloat64([]byte("3.25"), withScale("10", "")); got != 32.5 {
		t.Fatalf("[]byte cast = %v, want 32.5", got)
	}
	if got := castFloat64("2.5", withScale("4", "")); got != 10 {
		t.Fatalf("string cast = %v, want 10", got)
	}
	if got := castFloat64(true, nil); got != 1 {
		t.Fatalf("bool true cast = %v, want 1", got)
	}
	if got := castFloat64(false, nil); got != 0 {
		t.Fatalf("bool false cast = %v, want 0", got)
	}
	if got := castFloat64(nil, withScale("", "2.5")); got != 2.5 {
		t.Fatalf("nil default cast = %v, want 2.5", got)
	}
	if got := castFloat64(nil, withScale("10", "2.5")); got != 25 {
		t.Fatalf("nil default cast should apply scale: %v, want 25", got)
	}

	if got := castFloat64("abc", nil); !math.IsNaN(got) {
		t.Fatalf("invalid string cast = %v, want NaN", got)
	}
	if got := castFloat64(struct{}{}, nil); !math.IsNaN(got) {
		t.Fatalf("unknown type cast = %v, want NaN", got)
	}

	// Invalid numeric options should be rejected at parse time.
	if err := (&Column{Scale: "bad"}).parseNumbers(); err == nil {
		t.Fatal("parseNumbers should fail for invalid scale")
	}
	if err := (&Column{Default: "bad"}).parseNumbers(); err == nil {
		t.Fatal("parseNumbers should fail for invalid default")
	}
}

func TestCastString(t *testing.T) {
	now := time.Unix(1700000000, 0)
	if got := castString(int64(3)); got != "3" {
		t.Fatalf("int64 cast = %q, want 3", got)
	}
	if got := castString(float64(1.5)); got != "1.5" {
		t.Fatalf("float64 cast = %q, want 1.5", got)
	}
	if got := castString(now); got != "1700000000" {
		t.Fatalf("time cast = %q, want 1700000000", got)
	}
	if got := castString([]byte("abc")); got != "abc" {
		t.Fatalf("[]byte cast = %q, want abc", got)
	}
	if got := castString(true); got != "true" {
		t.Fatalf("bool true cast = %q, want true", got)
	}
	if got := castString(nil); got != "" {
		t.Fatalf("nil cast = %q, want empty", got)
	}
}

func TestConfigureLogger(t *testing.T) {
	if l := configureLogger("debug", "json"); l == nil {
		t.Fatal("configureLogger returned nil for valid json format")
	}
	if l := configureLogger("bad-level", "logfmt"); l == nil {
		t.Fatal("configureLogger returned nil for fallback level")
	}
	if l := configureLogger("info", "unknown-format"); l == nil {
		t.Fatal("configureLogger returned nil for unknown format fallback")
	}
}

func TestLogHelpersWithNilLogger(t *testing.T) {
	origin := Logger
	Logger = nil
	t.Cleanup(func() { Logger = origin })

	logDebugf("debug %d", 1)
	logInfof("info %d", 1)
	logWarnf("warn %d", 1)
	logErrorf("error %d", 1)
	logError("plain error")
}
