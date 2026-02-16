package exporter

import (
	"reflect"
	"testing"

	"github.com/alecthomas/kingpin/v2"
)

func TestNormalizeKingpinBoolEqualsArgs_Long(t *testing.T) {
	app := kingpin.New("test", "")
	app.Flag("auto-discovery", "").Short('a').Default("true").Bool()
	app.Flag("disable-cache", "").Short('C').Default("false").Bool()
	app.Flag("log.level", "").Default("info").String()

	tests := []struct {
		in   []string
		want []string
	}{
		{[]string{"--auto-discovery=false"}, []string{"--no-auto-discovery"}},
		{[]string{"--auto-discovery=true"}, []string{"--auto-discovery"}},
		{[]string{"--disable-cache=false"}, []string{"--no-disable-cache"}},
		{[]string{"--disable-cache=true"}, []string{"--disable-cache"}},
		{[]string{"--log.level=debug"}, []string{"--log.level=debug"}},
		{[]string{"--unknown=false"}, []string{"--unknown=false"}},
	}

	for _, tt := range tests {
		got := normalizeKingpinBoolEqualsArgs(tt.in, app.Model())
		if !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("normalize(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestNormalizeKingpinBoolEqualsArgs_Short(t *testing.T) {
	app := kingpin.New("test", "")
	app.Flag("auto-discovery", "").Short('a').Default("true").Bool()
	app.Flag("disable-cache", "").Short('C').Default("false").Bool()
	app.Flag("dry-run", "").Short('D').Default("false").Bool()

	tests := []struct {
		in   []string
		want []string
	}{
		{[]string{"-a=false"}, []string{"--no-auto-discovery"}},
		{[]string{"-a=true"}, []string{"-a"}},
		{[]string{"-C=false"}, []string{"--no-disable-cache"}},
		{[]string{"-C=true"}, []string{"-C"}},
		{[]string{"-D=false"}, []string{"--no-dry-run"}},
		{[]string{"-D=true"}, []string{"-D"}},
		{[]string{"-x=false"}, []string{"-x=false"}}, // unknown short
	}

	for _, tt := range tests {
		got := normalizeKingpinBoolEqualsArgs(tt.in, app.Model())
		if !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("normalize(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

