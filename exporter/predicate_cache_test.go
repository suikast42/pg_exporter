package exporter

import (
	"context"
	"testing"
	"time"
)

func TestPredicateCacheHitSkipsDBQuery(t *testing.T) {
	q := &Query{
		Name:   "q",
		Branch: "q",
		PredicateQueries: []PredicateQuery{
			{Name: "p1", SQL: "SELECT true", TTL: 10},
		},
	}
	s := &Server{Database: "postgres"} // DB is nil; any QueryContext call would panic.
	c := NewCollector(q, s)

	now := time.Now()
	c.scrapeBegin = now

	// Cache hit: pass=true should continue and ultimately return true without touching DB.
	c.predicateCache[0] = predicateCacheEntry{at: now.Add(-time.Second), pass: true}
	if ok := c.executePredicateQueries(context.Background()); !ok {
		t.Fatal("expected cached pass=true to allow query execution")
	}

	// Cache hit: pass=false should return false without touching DB.
	c.scrapeBegin = now
	c.predicateCache[0] = predicateCacheEntry{at: now.Add(-time.Second), pass: false}
	if ok := c.executePredicateQueries(context.Background()); ok {
		t.Fatal("expected cached pass=false to skip query execution")
	}
}

func TestPredicateCacheMissTriggersDBQuery(t *testing.T) {
	q := &Query{
		Name:   "q",
		Branch: "q",
		PredicateQueries: []PredicateQuery{
			{Name: "p1", SQL: "SELECT true", TTL: 10},
		},
	}
	s := &Server{Database: "postgres"} // DB is nil; QueryContext must panic if called.
	c := NewCollector(q, s)
	c.scrapeBegin = time.Now()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic due to DB access on predicate cache miss")
		}
	}()
	_ = c.executePredicateQueries(context.Background())
}

func TestPredicateCacheDisabledByTTLZero(t *testing.T) {
	q := &Query{
		Name:   "q",
		Branch: "q",
		PredicateQueries: []PredicateQuery{
			{Name: "p1", SQL: "SELECT true", TTL: 0},
		},
	}
	s := &Server{Database: "postgres"} // DB is nil; QueryContext must panic if called.
	c := NewCollector(q, s)
	c.scrapeBegin = time.Now()
	c.predicateCache[0] = predicateCacheEntry{at: time.Now(), pass: true}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic due to DB access when predicate TTL is 0 (cache disabled)")
		}
	}()
	_ = c.executePredicateQueries(context.Background())
}
