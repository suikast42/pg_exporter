package exporter

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func loadRepoConfigForTest(t *testing.T) map[string]*Query {
	t.Helper()

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
	return queries
}

func applicableAt(qs map[string]*Query, name string, version int) []*Query {
	var appl []*Query
	for _, q := range qs {
		if q.Name != name {
			continue
		}
		if q.MinVersion != 0 && version < q.MinVersion {
			continue
		}
		if q.MaxVersion != 0 && version >= q.MaxVersion {
			continue
		}
		appl = append(appl, q)
	}
	return appl
}

func hasCollector(qs map[string]*Query, name string) bool {
	for _, q := range qs {
		if q.Name == name {
			return true
		}
	}
	return false
}

func requireColumn(t *testing.T, q *Query, name, usage string) *Column {
	t.Helper()

	if q == nil {
		t.Fatal("query is nil")
	}
	col := q.Columns[name]
	if col == nil {
		t.Fatalf("%s should define column %q", q.Branch, name)
	}
	if col.Usage != usage {
		t.Fatalf("%s column %q has usage %s, want %s", q.Branch, name, col.Usage, usage)
	}
	return col
}

func requireLabelNames(t *testing.T, q *Query, want ...string) {
	t.Helper()

	if len(q.LabelNames) != len(want) {
		t.Fatalf("%s labels = %v, want %v", q.Branch, q.LabelNames, want)
	}
	for i, labelCol := range q.LabelNames {
		got := labelCol
		if col := q.Columns[labelCol]; col != nil && col.Rename != "" {
			got = col.Rename
		}
		if got != want[i] {
			t.Fatalf("%s labels = %v, want %v", q.Branch, q.LabelNames, want)
		}
	}
}

func requireTTL(t *testing.T, q *Query, want float64) {
	t.Helper()

	if q == nil {
		t.Fatal("query is nil")
	}
	if q.TTL != want {
		t.Fatalf("%s ttl = %v, want %v", q.Branch, q.TTL, want)
	}
}

func requireNoColumns(t *testing.T, q *Query, names ...string) {
	t.Helper()

	if q == nil {
		t.Fatal("query is nil")
	}
	for _, name := range names {
		if q.Columns[name] != nil {
			t.Fatalf("%s should not define column %q", q.Branch, name)
		}
	}
}

func metricNames(q *Query) map[string]bool {
	out := make(map[string]bool, len(q.MetricNames))
	for _, name := range q.MetricNames {
		col := q.Columns[name]
		if col == nil {
			continue
		}
		if col.Rename != "" {
			out[col.Rename] = true
		} else {
			out[col.Name] = true
		}
	}
	return out
}

func requireMetricSubset(t *testing.T, oldQ, newQ *Query) {
	t.Helper()

	newMetrics := metricNames(newQ)
	for name := range metricNames(oldQ) {
		if !newMetrics[name] {
			t.Fatalf("%s metric %q should remain available in %s", oldQ.Branch, name, newQ.Branch)
		}
	}
}

// Ensure the repo-bundled config/ covers PG10..PG19 without version gaps for
// collectors that are supposed to work on PG10+. This is a cheap static check
// (no DB required) to catch off-by-one mistakes on min/max_version splits.
func TestConfigCoveragePG10To19(t *testing.T) {
	queries := loadRepoConfigForTest(t)

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

		for v := 100000; v <= 190000; v += 10000 { // PG10..PG19
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

func TestConfigPG19BranchSelection(t *testing.T) {
	queries := loadRepoConfigForTest(t)

	for _, tc := range []struct {
		name   string
		branch string
	}{
		{name: "pg_sub", branch: "pg_sub_19"},
		// pg_stat_wal_receiver has no column change in PG19 (only a new status
		// value handled by the shared CASE mapping), so PG19 reuses pg_recv_13.
		{name: "pg_recv", branch: "pg_recv_13"},
		// pg_replication_slots only gains a new invalidation_reason value for
		// this collector's current metric surface, so PG19 reuses pg_slot_17.
		{name: "pg_slot", branch: "pg_slot_17"},
		{name: "pg_backup", branch: "pg_backup"},
		{name: "pg_wal", branch: "pg_wal_19"},
		// pg_stat_database_conflicts gains stats_reset in PG19, but this
		// collector intentionally keeps the core conflict counter surface and
		// avoids a PG19-only branch for that reset timestamp.
		{name: "pg_db_confl", branch: "pg_db_confl_16"},
		// pg_stat_progress_cluster exists across the whole v12+ support range,
		// so PG19 should not add a maintenance-only branch for the same surface.
		{name: "pg_clustering", branch: "pg_clustering"},
		// PG19 adds descriptive pg_stat_progress_vacuum fields, but no new
		// default historical progress surface; reuse the PG18 collector.
		{name: "pg_vacuuming", branch: "pg_vacuuming_18"},
		{name: "pg_recovery_state", branch: "pg_recovery_state"},
		{name: "pg_lock_stat", branch: "pg_lock_stat"},
		{name: "pg_vacuum_score", branch: "pg_vacuum_score"},
	} {
		if !hasCollector(queries, tc.name) {
			continue
		}
		t.Run(tc.name, func(t *testing.T) {
			appl := applicableAt(queries, tc.name, 190000)
			if len(appl) != 1 {
				t.Fatalf("server_version_num=190000 should select exactly one %s branch, got %d", tc.name, len(appl))
			}
			if appl[0].Branch != tc.branch {
				t.Fatalf("server_version_num=190000 selected %s, want %s", appl[0].Branch, tc.branch)
			}
		})
	}
}

func TestConfigPG19CompatibilityEssentials(t *testing.T) {
	queries := loadRepoConfigForTest(t)

	if pgSub := queries["pg_sub_19"]; pgSub != nil {
		if !strings.Contains(pgSub.SQL, "FROM pg_stat_subscription_stats s2") {
			t.Fatal("pg_sub_19 should anchor on pg_stat_subscription_stats so subscriptions without workers still emit stats")
		}
		if strings.Contains(pgSub.SQL, "s2.sync_error_count") {
			t.Fatal("pg_sub_19 should not reference removed PG19 sync_error_count field")
		}
		for _, col := range []struct {
			name  string
			usage string
		}{
			{name: "sync_error_count", usage: COUNTER},
			{name: "sync_table_error_count", usage: COUNTER},
			{name: "sync_seq_error_count", usage: COUNTER},
			{name: "has_worker", usage: GAUGE},
			{name: "reset_time", usage: GAUGE},
		} {
			requireColumn(t, pgSub, col.name, col.usage)
		}
		for _, col := range []string{
			"confl_insert_exists",
			"confl_update_origin_differs",
			"confl_update_exists",
			"confl_update_deleted",
			"confl_update_missing",
			"confl_delete_origin_differs",
			"confl_delete_missing",
			"confl_multiple_unique_conflicts",
		} {
			requireColumn(t, pgSub, col, COUNTER)
		}
	}

	if queries["pg_backup_19"] != nil {
		t.Fatal("PG19 should reuse pg_backup instead of adding a low-utility backup_type branch")
	}
	if pgBackup := queries["pg_backup"]; pgBackup != nil {
		if !strings.Contains(pgBackup.SQL, "pg_stat_get_progress_info") ||
			strings.Contains(pgBackup.SQL, "pg_stat_progress_basebackup") {
			t.Fatal("pg_backup should keep using pg_stat_get_progress_info for the stable core progress surface")
		}
		requireLabelNames(t, pgBackup, "pid")
		if pgBackup.Columns["backup_type"] != nil ||
			pgBackup.Columns["tablespaces_total"] != nil ||
			pgBackup.Columns["tablespaces_streamed"] != nil {
			t.Fatal("pg_backup should avoid low-utility PG19-only backup_type/tablespace fields")
		}
		for _, col := range []string{"total_bytes", "sent_bytes"} {
			requireColumn(t, pgBackup, col, GAUGE)
		}
		if pgBackup.Columns["backup_total"] != nil || pgBackup.Columns["backup_streamed"] != nil {
			t.Fatal("pg_backup should keep backward-compatible total_bytes/sent_bytes metric names")
		}
		if pgBackup.Columns["phase_name"] != nil {
			t.Fatal("pg_backup should keep phase as a gauge, not a churn-prone label")
		}
	}

	if pgRecv13 := queries["pg_recv_13"]; pgRecv13 != nil {
		for state, code := range map[string]string{"streaming": "0", "startup": "1", "catchup": "2", "backup": "3", "stopping": "4", "connecting": "5"} {
			if want := "WHEN '" + state + "' THEN " + code; !strings.Contains(pgRecv13.SQL, want) {
				t.Fatalf("pg_recv_13 should preserve status mapping with %q", want)
			}
		}
		for _, state := range []string{"starting", "waiting", "restarting"} {
			if strings.Contains(pgRecv13.SQL, state) {
				t.Fatalf("pg_recv_13 should not remap historical receiver transient state %q as part of PG19 support", state)
			}
		}
	}
	for _, branch := range []string{"pg_recv_10", "pg_recv_11"} {
		q := queries[branch]
		if q == nil {
			continue
		}
		if strings.Contains(q.SQL, "connecting") {
			t.Fatalf("%s should not include the PG19-only connecting status", branch)
		}
		for state, code := range map[string]string{"streaming": "0", "startup": "1", "catchup": "2", "backup": "3", "stopping": "4"} {
			if want := "WHEN '" + state + "' THEN " + code; !strings.Contains(q.SQL, want) {
				t.Fatalf("%s should preserve status mapping with %q", branch, want)
			}
		}
		for _, state := range []string{"starting", "waiting", "restarting"} {
			if strings.Contains(q.SQL, state) {
				t.Fatalf("%s should not remap historical receiver transient state %q as part of PG19 support", branch, state)
			}
		}
	}

	if queries["pg_vacuuming_19"] != nil {
		t.Fatal("PG19 should reuse pg_vacuuming_18 instead of adding a label-changing pg_vacuuming_19 branch")
	}
	if pgVacuuming := queries["pg_vacuuming_18"]; pgVacuuming != nil {
		requireLabelNames(t, pgVacuuming, "datname", "pid", "relname")
		requireColumn(t, pgVacuuming, "pid", LABEL)
		requireColumn(t, pgVacuuming, "relname", LABEL)
		for _, col := range []string{"progress", "indexes_total", "indexes_processed", "dead_tuple_bytes"} {
			requireColumn(t, pgVacuuming, col, GAUGE)
		}
		requireColumn(t, pgVacuuming, "delay_time", COUNTER)
		requireTTL(t, pgVacuuming, 10)
		requireNoColumns(t, pgVacuuming,
			"phase", "count", "mode", "started_by", "dead_tuple_mem_ratio",
			"index_vacuum_count", "max_dead_tuple_bytes", "num_dead_item_ids",
		)
		if strings.Contains(pgVacuuming.SQL, "GROUP BY") {
			t.Fatal("pg_vacuuming should stay row-compatible with the existing vacuum progress collectors")
		}
	}

	for _, branch := range []string{"pg_db_confl_15", "pg_db_confl_16"} {
		pgDBConfl := queries[branch]
		if pgDBConfl == nil {
			continue
		}
		if strings.Contains(pgDBConfl.SQL, "*") {
			t.Fatalf("%s should explicitly select conflict columns so future view columns are not exported accidentally", branch)
		}
	}
	if pgDBConfl := queries["pg_db_confl_16"]; pgDBConfl != nil {
		if strings.Contains(pgDBConfl.SQL, "stats_reset") || pgDBConfl.Columns["reset_time"] != nil {
			t.Fatal("pg_db_confl should not add a PG19 branch only to expose stats_reset")
		}
	}

	if queries["pg_clustering_19"] != nil {
		t.Fatal("pg_clustering should reuse the v12+ progress view branch instead of adding a PG19-only branch")
	}
	if pgClustering := queries["pg_clustering"]; pgClustering != nil {
		if !strings.Contains(pgClustering.SQL, "pg_stat_progress_cluster") ||
			strings.Contains(pgClustering.SQL, "pg_stat_get_progress_info") {
			t.Fatal("pg_clustering should use pg_stat_progress_cluster")
		}
		if !strings.Contains(pgClustering.SQL, "ELSE 0 END AS progress") {
			t.Fatal("pg_clustering should preserve the previous zero progress value before block totals are available")
		}
	}

	if pgLockStat := queries["pg_lock_stat"]; pgLockStat != nil {
		if requireColumn(t, pgLockStat, "wait_time", COUNTER).Scale != "1e-3" {
			t.Fatal("pg_lock_stat wait_time should scale PG milliseconds to seconds")
		}
	}
	if pgVacuumScore := queries["pg_vacuum_score"]; pgVacuumScore != nil {
		if pgVacuumScore.HasTag("cluster") || !pgVacuumScore.HasTag("primary") {
			t.Fatal("pg_vacuum_score should be current-database scoped and primary-only")
		}
		if pgVacuumScore.Timeout < 1 {
			t.Fatalf("pg_vacuum_score timeout = %v, want explicit table-scan timeout", pgVacuumScore.Timeout)
		}
		requireLabelNames(t, pgVacuumScore, "datname")
		for _, col := range []string{
			"max_score",
			"max_xid_score",
			"max_mxid_score",
			"max_vacuum_score",
			"max_vacuum_insert_score",
			"max_analyze_score",
			"table_count",
			"candidate_count",
			"vacuum_candidate_count",
			"analyze_candidate_count",
			"wraparound_candidate_count",
		} {
			requireColumn(t, pgVacuumScore, col, GAUGE)
		}
		if strings.Contains(pgVacuumScore.SQL, "ORDER BY") || strings.Contains(pgVacuumScore.SQL, "LIMIT") {
			t.Fatal("pg_vacuum_score should export stable aggregate rows, not volatile top-N table labels")
		}
		if !strings.Contains(pgVacuumScore.SQL, "count(*) FILTER (WHERE do_vacuum OR do_analyze)") {
			t.Fatal("pg_vacuum_score should expose the combined autovacuum candidate count")
		}
		if !strings.Contains(pgVacuumScore.SQL, "count(*) FILTER (WHERE for_wraparound)") {
			t.Fatal("pg_vacuum_score should expose wraparound candidate count")
		}
	}
}

func TestConfigPublishedVacuumingBranchesStayCompatible(t *testing.T) {
	queries := loadRepoConfigForTest(t)

	if pgVacuuming18 := queries["pg_vacuuming_18"]; pgVacuuming18 != nil {
		requireLabelNames(t, pgVacuuming18, "datname", "pid", "relname")
		for _, col := range []string{"progress", "indexes_total", "indexes_processed", "dead_tuple_bytes"} {
			requireColumn(t, pgVacuuming18, col, GAUGE)
		}
		requireColumn(t, pgVacuuming18, "delay_time", COUNTER)
	}

	if pgVacuuming17 := queries["pg_vacuuming_17"]; pgVacuuming17 != nil {
		requireLabelNames(t, pgVacuuming17, "datname", "pid", "relname")
		for _, col := range []string{"progress", "indexes_total", "indexes_processed", "dead_tuple_bytes"} {
			requireColumn(t, pgVacuuming17, col, GAUGE)
		}
	}

	if pgVacuuming12 := queries["pg_vacuuming_12"]; pgVacuuming12 != nil {
		requireLabelNames(t, pgVacuuming12, "datname", "pid", "relname")
		requireColumn(t, pgVacuuming12, "progress", GAUGE)
	}
}

func TestConfigPG19MetricCompatibility(t *testing.T) {
	queries := loadRepoConfigForTest(t)

	for _, name := range []string{
		"pg_sub",
		"pg_backup",
		"pg_recv",
		"pg_slot",
		"pg_wal",
		"pg_db_confl",
		"pg_clustering",
		"pg_vacuuming",
	} {
		if !hasCollector(queries, name) {
			continue
		}
		oldBranches := applicableAt(queries, name, 180000)
		newBranches := applicableAt(queries, name, 190000)
		if len(oldBranches) != 1 || len(newBranches) != 1 {
			t.Fatalf("%s should have exactly one PG18 branch and one PG19 branch, got %d/%d", name, len(oldBranches), len(newBranches))
		}
		requireMetricSubset(t, oldBranches[0], newBranches[0])
	}

	for _, name := range []string{"pg_sub", "pg_backup", "pg_recv", "pg_slot", "pg_wal", "pg_db_confl", "pg_vacuuming"} {
		if !hasCollector(queries, name) {
			continue
		}
		oldBranches := applicableAt(queries, name, 180000)
		newBranches := applicableAt(queries, name, 190000)
		requireLabelNames(t, newBranches[0], oldBranches[0].LabelNames...)
	}
}
