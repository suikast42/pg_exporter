package exporter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// GetConfig will try load config from target path
func GetConfig() (res string) {
	// priority: cli-args > env  > default settings (check exist)
	if res = *configPath; res != "" {
		logInfof("retrieve config path %s from command line", res)
		return res
	}
	if res = os.Getenv("PG_EXPORTER_CONFIG"); res != "" {
		logInfof("retrieve config path %s from PG_EXPORTER_CONFIG", res)
		return res
	}

	candidate := []string{"pg_exporter.yml", "/etc/pg_exporter.yml", "/etc/pg_exporter"}
	for _, res = range candidate {
		if _, err := os.Stat(res); err == nil { // default1 exist
			logInfof("fallback on default config path: %s", res)
			return res
		}
	}
	return ""
}

// ParseConfig turn config content into Query struct
func ParseConfig(content []byte) (queries map[string]*Query, err error) {
	queries = make(map[string]*Query)
	if err = yaml.Unmarshal(content, &queries); err != nil {
		return nil, fmt.Errorf("malformed config: %w", err)
	}

	// parse additional fields
	for branch, query := range queries {
		if query == nil {
			return nil, fmt.Errorf("query %q is null", branch)
		}
		query.Branch = branch
		if query.Name == "" {
			query.Name = branch
		}
		if strings.TrimSpace(query.SQL) == "" {
			return nil, fmt.Errorf("query %q has empty SQL", branch)
		}
		if query.TTL < 0 {
			return nil, fmt.Errorf("query %q has negative ttl: %v", branch, query.TTL)
		}
		for i, pq := range query.PredicateQueries {
			if strings.TrimSpace(pq.SQL) == "" {
				return nil, fmt.Errorf("query %q has empty predicate_query at index %d", branch, i)
			}
			if pq.TTL < 0 {
				return nil, fmt.Errorf("query %q has negative predicate_queries[%d].ttl: %v", branch, i, pq.TTL)
			}
		}
		if len(query.Metrics) == 0 {
			return nil, fmt.Errorf("query %q has no metrics definition", branch)
		}
		// parse query column info
		columns := make(map[string]*Column, len(query.Metrics))
		var allColumns, labelColumns, metricColumns []string
		for _, colMap := range query.Metrics {
			if len(colMap) == 0 {
				return nil, fmt.Errorf("query %q has an empty metrics entry", branch)
			}
			if len(colMap) != 1 {
				return nil, fmt.Errorf("query %q has invalid metrics entry with %d columns, expect exactly 1", branch, len(colMap))
			}
			for colName, column := range colMap { // one-entry map
				if column == nil {
					return nil, fmt.Errorf("query %q has null column definition for %q", branch, colName)
				}
				if column.Name == "" {
					column.Name = colName
				}
				usage := strings.ToUpper(strings.TrimSpace(column.Usage))
				if usage == "" {
					return nil, fmt.Errorf("query %q column %q has empty usage", branch, colName)
				}
				if _, isValid := ColumnUsage[usage]; !isValid {
					return nil, fmt.Errorf("query %q column %q has unsupported usage: %s", branch, colName, column.Usage)
				}
				column.Usage = usage
				if err := column.parseNumbers(); err != nil {
					return nil, fmt.Errorf("query %q column %q: %w", branch, colName, err)
				}
				switch column.Usage {
				case LABEL:
					labelColumns = append(labelColumns, column.Name)
				case GAUGE, COUNTER:
					metricColumns = append(metricColumns, column.Name)
				}
				allColumns = append(allColumns, column.Name)
				if _, exists := columns[column.Name]; exists {
					return nil, fmt.Errorf("query %q has duplicate column name %q", branch, column.Name)
				}
				columns[column.Name] = column
			}
		}
		if len(metricColumns) == 0 {
			return nil, fmt.Errorf("query %q defines no GAUGE/COUNTER columns", branch)
		}
		query.Columns, query.ColumnNames, query.LabelNames, query.MetricNames = columns, allColumns, labelColumns, metricColumns

		// Validate prometheus label names and metric names. This prevents panics at scrape time.
		seenLabels := make(map[string]bool, len(query.LabelNames))
		for _, labelColName := range query.LabelNames {
			c := query.Columns[labelColName]
			if c == nil {
				return nil, fmt.Errorf("query %q missing label column %q", branch, labelColName)
			}
			lbl := c.Name
			if c.Rename != "" {
				lbl = c.Rename
			}
			if err := validatePromLabelName(lbl); err != nil {
				return nil, fmt.Errorf("query %q label %q: %w", branch, lbl, err)
			}
			if seenLabels[lbl] {
				return nil, fmt.Errorf("query %q has duplicate label name %q", branch, lbl)
			}
			seenLabels[lbl] = true
		}

		seenMetrics := make(map[string]bool, len(query.MetricNames))
		for _, metricColName := range query.MetricNames {
			c := query.Columns[metricColName]
			if c == nil {
				return nil, fmt.Errorf("query %q missing metric column %q", branch, metricColName)
			}
			suffix := c.Name
			if c.Rename != "" {
				suffix = c.Rename
			}
			metricName := fmt.Sprintf("%s_%s", query.Name, suffix)
			if err := validatePromMetricName(metricName); err != nil {
				return nil, fmt.Errorf("query %q metric %q: %w", branch, metricName, err)
			}
			if seenMetrics[metricName] {
				return nil, fmt.Errorf("query %q has duplicate metric name %q", branch, metricName)
			}
			seenMetrics[metricName] = true
		}
	}
	return
}

func FinalizeQueries(queries map[string]*Query, source string) error {
	for branch, q := range queries {
		if q == nil {
			return fmt.Errorf("query %q is null", branch)
		}
		q.Path = source
		// If timeout is not set, set to 100ms by default.
		// If timeout is set to a negative number, set to 0 (disabled).
		if q.Timeout == 0 {
			q.Timeout = 0.1
		}
		if q.Timeout < 0 {
			q.Timeout = 0
		}
	}
	return nil
}

// ParseQuery generate a single query from config string
func ParseQuery(config string) (*Query, error) {
	queries, err := ParseConfig([]byte(config))
	if err != nil {
		return nil, err
	}
	if len(queries) == 0 {
		return nil, fmt.Errorf("no query definition found")
	}
	if len(queries) > 1 {
		return nil, fmt.Errorf("multiple query definition found")
	}
	if err := FinalizeQueries(queries, "<inline>"); err != nil {
		return nil, err
	}
	for _, q := range queries {
		return q, nil // return the only query instance
	}
	return nil, fmt.Errorf("no query definition found")
}

// LoadConfig will read single conf file or read multiple conf file if a dir is given
// conf file in a dir will be load in alphabetic order, query with same name will overwrite predecessor
func LoadConfig(configPath string) (queries map[string]*Query, err error) {
	stat, err := os.Stat(configPath)
	if err != nil {
		return nil, fmt.Errorf("invalid config path: %s: %w", configPath, err)
	}
	if stat.IsDir() { // iterate conf files (non-recursive) if a dir is given
		files, err := os.ReadDir(configPath)
		if err != nil {
			return nil, fmt.Errorf("fail reading config dir: %s: %w", configPath, err)
		}

		logDebugf("load config from dir: %s", configPath)
		confFiles := make([]string, 0)
		for _, conf := range files {
			if conf.IsDir() {
				continue // skip subdirectories
			}
			if !(strings.HasSuffix(conf.Name(), ".yaml") || strings.HasSuffix(conf.Name(), ".yml")) {
				continue // skip non-yaml files
			}
			confFiles = append(confFiles, filepath.Join(configPath, conf.Name()))
		}

		// make global config map and assign priority according to config file alphabetic orders
		// priority is an integer range from 1 to 999, where 1 - 99 is reserved for user
		queries = make(map[string]*Query)
		var queryCount, configCount int
		var firstErr error
		for _, confPath := range confFiles {
			if singleQueries, err := LoadConfig(confPath); err != nil {
				logWarnf("skip config %s due to error: %s", confPath, err.Error())
				if firstErr == nil {
					firstErr = err
				}
			} else {
				configCount++
				for name, query := range singleQueries {
					queryCount++
					if query.Priority == 0 { // set to config rank if not manually set
						query.Priority = 100 + configCount
					}
					queries[name] = query // so the later one will overwrite former one
				}
			}
		}
		if len(confFiles) > 0 && len(queries) == 0 {
			if firstErr != nil {
				return nil, fmt.Errorf("no valid queries loaded from config dir %s (%d yaml files), first error: %w", configPath, len(confFiles), firstErr)
			}
			return nil, fmt.Errorf("no queries loaded from config dir %s (%d yaml files)", configPath, len(confFiles))
		}
		logDebugf("load %d of %d queries from %d config files", len(queries), queryCount, configCount)
		return queries, nil
	}

	// single file case: recursive exit condition
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("fail reading config file %s: %w", configPath, err)
	}
	queries, err = ParseConfig(content)
	if err != nil {
		return nil, err
	}
	if err := FinalizeQueries(queries, stat.Name()); err != nil {
		return nil, err
	}
	logDebugf("load %d queries from %s", len(queries), configPath)
	return queries, nil

}
