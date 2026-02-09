---
title: Deployment
weight: 6
breadcrumbs: false
comments: false
---


This guide covers production deployment strategies, best practices, and real-world configurations for PG Exporter.

pg_exporter itself can be configured through:

1. **Command-line arguments** (higher priority)
2. **Environment variables** (lower priority)

The metrics collectors are configured with a YAML configuration file (dir/files):

- `/etc/pg_exporter.yml` (default)
- `/etc/pg_exporter/` (directory with multiple files)


The configuration file uses YAML format and consists of **collector definitions** that specify what metrics to collect and how to collect them.

### CLI Arg

All configuration options can be specified via command-line flags:

```bash
pg_exporter \
  --url="postgres://user:pass@localhost:5432/postgres" \
  --config="/etc/pg_exporter/pg_exporter.yml" \
  --web.listen-address=":9630" \
  --auto-discovery \
  --exclude-database="template0,template1" \
  --log.level="info"
```

Run `pg_exporter --help` for a complete list of available flags:

```bash
Flags:
  -h, --[no-]help                Show context-sensitive help (also try --help-long and --help-man).
  -u, --url=URL                  postgres target url
  -c, --config=CONFIG            path to config dir or file
      --[no-]web.systemd-socket  Use systemd socket activation listeners instead of port listeners (Linux only).
      --web.listen-address=:9630 ...
                                 Addresses on which to expose metrics and web interface. Repeatable for multiple addresses. Examples: `:9100` or `[::1]:9100` for http, `vsock://:9100` for vsock
      --web.config.file=""       Path to configuration file that can enable TLS or authentication. See: https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md
  -l, --label=""                 constant lables:comma separated list of label=value pair ($PG_EXPORTER_LABEL)
  -t, --tag=""                   tags,comma separated list of server tag ($PG_EXPORTER_TAG)
  -C, --[no-]disable-cache       force not using cache ($PG_EXPORTER_DISABLE_CACHE)
  -m, --[no-]disable-intro       disable internal/exporter self metrics ($PG_EXPORTER_DISABLE_INTRO)
  -a, --[no-]auto-discovery      automatically scrape all database for given server ($PG_EXPORTER_AUTO_DISCOVERY)
  -x, --exclude-database="template0,template1,postgres"
                                 excluded databases when enabling auto-discovery ($PG_EXPORTER_EXCLUDE_DATABASE)
  -i, --include-database=""      included databases when enabling auto-discovery ($PG_EXPORTER_INCLUDE_DATABASE)
  -n, --namespace=""             prefix of built-in metrics, (pg|pgbouncer) by default ($PG_EXPORTER_NAMESPACE)
  -f, --[no-]fail-fast           fail fast instead of waiting during start-up ($PG_EXPORTER_FAIL_FAST)
  -T, --connect-timeout=100      connect timeout in ms, 100 by default ($PG_EXPORTER_CONNECT_TIMEOUT)
  -P, --web.telemetry-path="/metrics"
                                 URL path under which to expose metrics. ($PG_EXPORTER_TELEMETRY_PATH)
  -D, --[no-]dry-run             dry run and print raw configs
  -E, --[no-]explain             explain server planned queries
      --log.level="info"         log level: debug|info|warn|error]
      --log.format="logfmt"      log format: logfmt|json
      --[no-]version             Show application version.
```

### Environment Variables

All command-line arguments have corresponding environment variables:

```bash
PG_EXPORTER_URL='postgres://:5432/?sslmode=disable'
PG_EXPORTER_CONFIG=/etc/pg_exporter.yml
PG_EXPORTER_LABEL=""
PG_EXPORTER_TAG=""
PG_EXPORTER_DISABLE_CACHE=false
PG_EXPORTER_AUTO_DISCOVERY=true
PG_EXPORTER_EXCLUDE_DATABASE="template0,template1,postgres"
PG_EXPORTER_INCLUDE_DATABASE=""
PG_EXPORTER_NAMESPACE="pg"
PG_EXPORTER_FAIL_FAST=false
PG_EXPORTER_CONNECT_TIMEOUT=100
PG_EXPORTER_TELEMETRY_PATH="/metrics"
PG_EXPORTER_OPTS='--log.level=info'

pg_exporter
```



## Deployment Architecture

The simplest deployment with one exporter per PostgreSQL instance:

```
┌─────────────┐     ┌──────────────┐     ┌────────────┐
│ Prometheus  │────▶│ PG Exporter  │────▶│ PostgreSQL │
└─────────────┘     └──────────────┘     └────────────┘
                         :9630                :5432
```

### Multi-Database Environment

Using auto-discovery to monitor multiple databases, which is enabled by default

```
┌─────────────┐     ┌────────────────┐     ┌────────────┐
│ Prometheus  │────▶│ PG Exporter    │────▶│ PostgreSQL │
└─────────────┘     │     with       │     │  ├─ db1    │
                    │ auto-discovery │     │  ├─ db2    │
                    └────────────────┘     │  └─ db3    │
                                           └────────────┘
```



## Production Configuration

### PostgreSQL User Setup

Create a dedicated monitoring user with minimal required permissions:

```sql
-- Create monitoring role
CREATE ROLE pg_monitor WITH LOGIN PASSWORD 'strong_password' CONNECTION LIMIT 5;

-- Grant necessary permissions
GRANT pg_monitor TO pg_monitor;  -- PostgreSQL 10+ built-in role
GRANT CONNECT ON DATABASE postgres TO pg_monitor;

-- For specific databases
GRANT CONNECT ON DATABASE app_db TO pg_monitor;
GRANT USAGE ON SCHEMA public TO pg_monitor;

-- Additional permissions for extended monitoring
GRANT SELECT ON ALL TABLES IN SCHEMA pg_catalog TO pg_monitor;
GRANT SELECT ON ALL SEQUENCES IN SCHEMA pg_catalog TO pg_monitor;
```



### Connection Security

#### Using SSL/TLS

```bash
# Connection string with SSL
PG_EXPORTER_URL='postgres://pg_monitor:password@db.example.com:5432/postgres?sslmode=require&sslcert=/path/to/client.crt&sslkey=/path/to/client.key&sslrootcert=/path/to/ca.crt'
```

#### Using .pgpass File

```bash
# Create .pgpass file
echo "db.example.com:5432:*:pg_monitor:password" > ~/.pgpass
chmod 600 ~/.pgpass

# Use without password in URL
PG_EXPORTER_URL='postgres://pg_monitor@db.example.com:5432/postgres'
```




## Systemd Service Configuration

Complete production systemd setup:

```ini
[Unit]
Description=Prometheus exporter for PostgreSQL/Pgbouncer server metrics
Documentation=https://github.com/pgsty/pg_exporter
After=network.target

[Service]
EnvironmentFile=-/etc/default/pg_exporter
User=prometheus
ExecStart=/usr/bin/pg_exporter $PG_EXPORTER_OPTS
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

Environment file `/etc/default/pg_exporter`:

```bash
PG_EXPORTER_URL='postgres://:5432/?sslmode=disable'
PG_EXPORTER_CONFIG=/etc/pg_exporter.yml
PG_EXPORTER_LABEL=""
PG_EXPORTER_TAG=""
PG_EXPORTER_DISABLE_CACHE=false
PG_EXPORTER_AUTO_DISCOVERY=true
PG_EXPORTER_EXCLUDE_DATABASE="template0,template1,postgres"
PG_EXPORTER_INCLUDE_DATABASE=""
PG_EXPORTER_NAMESPACE="pg"
PG_EXPORTER_FAIL_FAST=false
PG_EXPORTER_CONNECT_TIMEOUT=100
PG_EXPORTER_TELEMETRY_PATH="/metrics"
PG_EXPORTER_OPTS='--log.level=info'
```
