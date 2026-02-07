package exporter

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func makeCachedCollectorForServer(s *Server, name string, val float64) *Collector {
	q := makeGaugeQuery(name, 1)
	c := NewCollector(q, s)
	c.TTL = 3600
	c.lastScrape = time.Now()
	metric := prometheus.MustNewConstMetric(c.descriptors["value"], prometheus.GaugeValue, val, "db")
	c.result = []prometheus.Metric{metric}
	c.err = nil
	return c
}

func TestExporterCollectAndInternalMetrics(t *testing.T) {
	primary := NewServer("postgresql://u:p@localhost:5432/postgres")
	primary.beforeScrape = func(s *Server) error {
		s.UP = true
		s.Version = 160000
		s.Recovery = false
		return nil
	}
	primary.Planned = true
	primary.Collectors = []*Collector{makeCachedCollectorForServer(primary, "q_primary", 1)}
	primary.ResetStats()

	extra := NewServer("postgresql://u:p@localhost:5432/otherdb")
	extra.Forked = true
	extra.beforeScrape = func(s *Server) error {
		s.UP = true
		s.Version = 160000
		s.Recovery = false
		return nil
	}
	extra.Planned = true
	extra.Collectors = []*Collector{makeCachedCollectorForServer(extra, "q_extra", 2)}
	extra.ResetStats()

	e := &Exporter{
		server:  primary,
		servers: map[string]*Server{"otherdb": extra},
	}
	e.setupInternalMetrics()

	ch := make(chan prometheus.Metric, 256)
	e.Collect(ch)

	if !e.Up() {
		t.Fatal("Exporter should be UP after successful collect")
	}
	if e.Status() != "primary" {
		t.Fatalf("Exporter status = %s, want primary", e.Status())
	}
}

func TestExporterDescribeAndCloseNoPanic(t *testing.T) {
	s := NewServer("postgresql://u:p@localhost:5432/postgres")
	s.beforeScrape = func(s *Server) error {
		s.UP = true
		return nil
	}
	s.Planned = true
	s.Collectors = []*Collector{makeCachedCollectorForServer(s, "q", 1)}
	s.ResetStats()

	e := &Exporter{
		server:  s,
		servers: map[string]*Server{},
	}
	e.setupInternalMetrics()

	descCh := make(chan *prometheus.Desc, 32)
	e.Describe(descCh)
	if len(descCh) == 0 {
		t.Fatal("Describe should emit at least one descriptor")
	}

	// server DB pointers are nil in this synthetic test; Close should not panic.
	e.Close()
}

func TestServerIntrospectionHelpers(t *testing.T) {
	s := NewServer("postgresql://u:p@localhost:5432/postgres")
	c := makeCachedCollectorForServer(s, "q", 1)
	s.Collectors = []*Collector{c}
	s.ResetStats()

	if s.Error() != nil {
		t.Fatalf("new server Error should be nil, got %v", s.Error())
	}
	if got := s.Duration(); got != 0 {
		t.Fatalf("new server Duration = %v, want 0", got)
	}
	if got := s.Uptime(); got < 0 {
		t.Fatalf("Uptime should be non-negative, got %v", got)
	}

	if got := c.ResultSize(); got != 1 {
		t.Fatalf("collector ResultSize = %d, want 1", got)
	}
	if skip, _ := c.PredicateSkip(); skip {
		t.Fatal("collector PredicateSkip should be false by default")
	}
	if got := c.Duration(); got != 0 {
		t.Fatalf("collector Duration = %v, want 0", got)
	}

	if exp := s.Explain(); exp == "" {
		t.Fatal("Explain should not be empty")
	}
	if html := s.ExplainHTML(); html == "" {
		t.Fatal("ExplainHTML should not be empty")
	}
}
