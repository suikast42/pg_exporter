---
title: Getting Started
weight: 1
breadcrumbs: false
comments: false
---

PG Exporter is an advanced PostgreSQL and pgBouncer metrics exporter for Prometheus. This guide will help you get up and running quickly.

## Prerequisites

Before you begin, ensure you have:

- PostgreSQL 10+ or pgBouncer 1.8+ instance to monitor
- A user account with appropriate permissions for monitoring
- Prometheus Compatible System (for metrics scraping)
- Basic understanding of PostgreSQL connection strings

## Quick Start

The fastest way to get started with PG Exporter:

```bash
# Download and install the latest release
curl -L https://github.com/pgsty/pg_exporter/releases/latest/download/pg_exporter-$(uname -s)-$(uname -m).tar.gz | tar xz
sudo install pg_exporter /usr/bin/

# Run with PostgreSQL connection URL
PG_EXPORTER_URL='postgres://user:pass@localhost:5432/postgres' pg_exporter

# Verify metrics are available
curl http://localhost:9630/metrics
```

## Understanding the Basics

### Connection String

PG Exporter uses standard PostgreSQL connection URLs:

```
postgres://[user][:password]@[host][:port]/[database][?param=value]
```

Examples:
- Local PostgreSQL: `postgres:///postgres`
- Remote with auth: `postgres://monitor:password@db.example.com:5432/postgres`
- With SSL: `postgres://user:pass@host/db?sslmode=require`
- pgBouncer: `postgres://pgbouncer:password@localhost:6432/pgbouncer`

### Built-in Metrics

PG Exporter provides 4 core built-in metrics out of the box:

| Metric | Type | Description |
|--------|------|-------------|
| `pg_up` | Gauge | 1 if exporter can connect to PostgreSQL, 0 otherwise |
| `pg_version` | Gauge | PostgreSQL server version number |
| `pg_in_recovery` | Gauge | 1 if server is in recovery mode (replica), 0 if primary |
| `pg_exporter_build_info` | Gauge | Exporter version and build information |

### Configuration File

All other metrics (600+) are defined in the `pg_exporter.yml` configuration file. By default, PG Exporter looks for this file in:

1. Path specified by `--config` flag
2. Path in `PG_EXPORTER_CONFIG` environment variable
3. Current directory (`./pg_exporter.yml`)
4. System config (`/etc/pg_exporter.yml` or `/etc/pg_exporter/`)

## Your First Monitoring Setup

### Step 1: Create a Monitoring User

Create a dedicated PostgreSQL user for monitoring:

```sql
-- Create monitoring user
CREATE USER pg_monitor WITH PASSWORD 'secure_password';

-- Grant necessary permissions
GRANT pg_monitor TO pg_monitor;
GRANT CONNECT ON DATABASE postgres TO pg_monitor;

-- For PostgreSQL 10+, pg_monitor role provides read access to monitoring views
-- For older versions, you may need additional grants
```

### Step 2: Test Connection

Verify the exporter can connect to your database:

```bash
# Set connection URL
export PG_EXPORTER_URL='postgres://pg_monitor:secure_password@localhost:5432/postgres'

# Run in dry-run mode to test configuration
pg_exporter --dry-run
```

### Step 3: Run the Exporter

Start PG Exporter:

```bash
# Run with default settings
pg_exporter

# Or with custom flags
pg_exporter \
  --url='postgres://pg_monitor:secure_password@localhost:5432/postgres' \
  --web.listen-address=':9630' \
  --log.level=info
```

### Step 4: Configure Prometheus

Add PG Exporter as a target in your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'postgresql'
    static_configs:
      - targets: ['localhost:9630']
        labels:
          instance: 'postgres-primary'
```

### Step 5: Verify Metrics

Check that metrics are being collected:

```bash
# View raw metrics
curl http://localhost:9630/metrics | grep pg_

# Check exporter statistics
curl http://localhost:9630/stat

# Verify server detection
curl http://localhost:9630/explain
```

## Auto-Discovery Mode

PG Exporter can automatically discover and monitor all databases in a PostgreSQL instance:

```bash
# Enable auto-discovery (default behavior)
pg_exporter --auto-discovery

# Exclude specific databases
pg_exporter --auto-discovery \
  --exclude-database="template0,template1,postgres"

# Include only specific databases
pg_exporter --auto-discovery \
  --include-database="app_db,analytics_db"
```

When auto-discovery is enabled:
- Cluster-level metrics (1xx-5xx) are collected once per instance
- Database-level metrics (6xx-8xx) are collected for each discovered database
- Metrics are labeled with `datname` to distinguish between databases

## Monitoring pgBouncer

To monitor pgBouncer instead of PostgreSQL:

```bash
# Connect to pgBouncer admin database
PG_EXPORTER_URL='postgres://pgbouncer:password@localhost:6432/pgbouncer' \
pg_exporter --config=/etc/pg_exporter.yml
```

Note: if you see `unsupported startup parameter: extra_float_digits` in the exporter logs while connecting to pgBouncer, add the following to your `pgbouncer.ini`:

```ini
[pgbouncer]
ignore_startup_parameters = extra_float_digits
```

The exporter automatically detects pgBouncer and:
- Uses `pgbouncer` namespace for metrics
- Executes pgBouncer-specific collectors (9xx series)
- Provides pgBouncer-specific health checks

## Using Docker

Run PG Exporter in a container:

```bash
docker run -d \
  --name pg_exporter \
  -p 9630:9630 \
  -e PG_EXPORTER_URL="postgres://user:pass@host.docker.internal:5432/postgres" \
  pgsty/pg_exporter:latest
```

With custom configuration:

```bash
docker run -d \
  --name pg_exporter \
  -p 9630:9630 \
  -v /path/to/pg_exporter.yml:/etc/pg_exporter.yml \
  -e PG_EXPORTER_URL="postgres://user:pass@db:5432/postgres" \
  pgsty/pg_exporter:latest
```

## Health Checks

PG Exporter provides health check endpoints for load balancers and orchestrators:

```bash
# Basic health check
curl http://localhost:9630/up
# Returns: 200 if connected, 503 if not

# Primary detection
curl http://localhost:9630/primary
# Returns: 200 if primary, 404 if replica, 503 if unknown

# Replica detection
curl http://localhost:9630/replica
# Returns: 200 if replica, 404 if primary, 503 if unknown
```

## Troubleshooting

### Connection Issues

```bash
# Test with detailed logging
pg_exporter --log.level=debug --dry-run

# Check server planning
pg_exporter --explain
```

### Permission Errors

Ensure the monitoring user has necessary permissions:

```sql
-- Check current permissions
SELECT * FROM pg_roles WHERE rolname = 'pg_monitor';

-- Grant additional permissions if needed
GRANT USAGE ON SCHEMA pg_catalog TO pg_monitor;
GRANT SELECT ON ALL TABLES IN SCHEMA pg_catalog TO pg_monitor;
```

### Slow Scrapes

If scrapes are timing out:

1. Check slow queries: `curl http://localhost:9630/stat`
2. Adjust collector timeouts in configuration
3. Use caching for expensive queries (set `ttl` in collector config)
4. Disable expensive collectors if not needed

## Next Steps

- [Installation Guide](../installation) - Detailed installation instructions for all platforms
- [Configuration Reference](../configuration) - Complete configuration documentation
- [Deployment Guide](../deployment) - Production deployment best practices
- [API Reference](../api) - Full API endpoint documentation
- [Collectors Documentation](../collectors) - Details on all available metrics collectors
