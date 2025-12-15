---
title: API Reference
weight: 4
---


PG Exporter provides a comprehensive REST API for metrics collection, health checking, traffic routing, and operational control. All endpoints are exposed via HTTP on the configured port (default: 9630).

## Endpoint Overview

| Endpoint                       | Method | Description                 |
|--------------------------------|--------|-----------------------------|
| [`/metrics`](#metrics)         | GET    | Prometheus metrics endpoint |
| [`/up`](#health-checks)        | GET    | Basic aliveness check       |
| [`/health`](#health-checks)    | GET    | Detailed health status      |
| [`/primary`](#traffic-routing) | GET    | Primary server check        |
| [`/replica`](#traffic-routing) | GET    | Replica server check        |
| [`/read`](#traffic-routing)    | GET    | Read traffic routing        |
| [`/reload`](#operational)      | GET    | Reload configuration        |
| [`/explain`](#operational)     | GET    | Explain query planning      |
| [`/stat`](#operational)        | GET    | Runtime statistics          |



## Metrics Endpoint

### GET /metrics

The primary endpoint that exposes all collected metrics in Prometheus format.

#### Request

```bash
curl http://localhost:9630/metrics
```

#### Response

```prometheus
# HELP pg_up PostgreSQL server is up and accepting connections
# TYPE pg_up gauge
pg_up 1

# HELP pg_version PostgreSQL server version number
# TYPE pg_version gauge
pg_version 140000

# HELP pg_in_recovery PostgreSQL server is in recovery mode
# TYPE pg_in_recovery gauge
pg_in_recovery 0

# HELP pg_exporter_build_info PG Exporter build information
# TYPE pg_exporter_build_info gauge
pg_exporter_build_info{version="1.1.0",branch="main",revision="abc123"} 1

# ... additional metrics
```

#### Response Format

Metrics follow the Prometheus exposition format:

```
# HELP <metric_name> <description>
# TYPE <metric_name> <type>
<metric_name>{<label_name>="<label_value>",...} <value> <timestamp>
```





## Health Check Endpoints

Health check endpoints provide various ways to monitor PG Exporter and the target database status.

### GET /up

Simple binary health check.

#### Response Codes

| Code | Status | Description |
|------|---------|-------------|
| 200 | OK | Exporter and database are up |
| 503 | Service Unavailable | Database is down or unreachable |

#### Example

```bash
# Check if service is up
curl -I http://localhost:9630/up

HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
```

### GET /health

Alias for `/up` with same behavior.

```bash
curl http://localhost:9630/health
```

### GET /liveness

Kubernetes liveness probe endpoint.

```bash
# Liveness probe configuration
livenessProbe:
  httpGet:
    path: /liveness
    port: 9630
  initialDelaySeconds: 30
  periodSeconds: 10
```

### GET /readiness

Kubernetes readiness probe endpoint.

```bash
# Readiness probe configuration
readinessProbe:
  httpGet:
    path: /readiness
    port: 9630
  initialDelaySeconds: 5
  periodSeconds: 5
```





## Traffic Routing Endpoints

These endpoints are designed for load balancers and proxies to route traffic based on server role.

### GET /primary

Check if the server is a primary (master) instance.

#### Response Codes

| Code | Status              | Description                            |
|------|---------------------|----------------------------------------|
| 200  | OK                  | Server is primary and accepting writes |
| 404  | Not Found           | Server is not primary (is replica)     |
| 503  | Service Unavailable | Server is down                         |

#### Aliases

- `/leader`
- `/master`
- `/read-write`
- `/rw`

#### Example

```bash
# Check if server is primary
curl -I http://localhost:9630/primary

# Use in HAProxy configuration
backend pg_primary
  option httpchk GET /primary
  server pg1 10.0.0.1:5432 check port 9630
  server pg2 10.0.0.2:5432 check port 9630
```

### GET /replica

Check if the server is a replica (standby) instance.

#### Response Codes

| Code | Status              | Description                        |
|------|---------------------|------------------------------------|
| 200  | OK                  | Server is replica and in recovery  |
| 404  | Not Found           | Server is not replica (is primary) |
| 503  | Service Unavailable | Server is down                     |

#### Aliases

- `/standby`
- `/slave`
- `/read-only`
- `/ro`

#### Example

```bash
# Check if server is replica
curl -I http://localhost:9630/replica

# Use in load balancer configuration
backend pg_replicas
  option httpchk GET /replica
  server pg2 10.0.0.2:5432 check port 9630
  server pg3 10.0.0.3:5432 check port 9630
```

### GET /read

Check if the server can handle read traffic (both primary and replica).

#### Response Codes

| Code | Status              | Description                       |
|------|---------------------|-----------------------------------|
| 200  | OK                  | Server is up and can handle reads |
| 503  | Service Unavailable | Server is down                    |

#### Example

```bash
# Check if server can handle reads
curl -I http://localhost:9630/read

# Route read traffic to any available server
backend pg_read
  option httpchk GET /read
  server pg1 10.0.0.1:5432 check port 9630
  server pg2 10.0.0.2:5432 check port 9630
  server pg3 10.0.0.3:5432 check port 9630
```

## Operational Endpoints

### POST /reload

Reload configuration without restarting the exporter.

#### Request

```bash
curl -X POST http://localhost:9630/reload
```

#### Response

```json
{
  "status": "success",
  "message": "Configuration reloaded successfully",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### Response Codes

| Code | Status                | Description                         |
|------|-----------------------|-------------------------------------|
| 200  | OK                    | Configuration reloaded successfully |
| 500  | Internal Server Error | Reload failed                       |

#### Use Cases

- Update collector definitions
- Change query parameters
- Modify cache TTL values
- Add or remove collectors

{{< callout type="info" >}}
Configuration reload does not affect the database connection. To change connection parameters, restart the exporter.
{{< /callout >}}

### GET /explain

Display query execution planning information for all configured collectors.

#### Request

```bash
curl http://localhost:9630/explain
```

#### Response

```text
Collector: pg_stat_database
  Query: SELECT datname, numbackends FROM pg_stat_database
  Tags: [cluster]
  TTL: 10s
  Timeout: 100ms
  Version: 100000-999999
  Status: Active

Collector: pg_stat_replication
  Query: SELECT client_addr, state FROM pg_stat_replication
  Tags: [primary]
  TTL: 5s
  Timeout: 100ms
  Version: 100000-999999
  Status: Active (primary only)

...
```
