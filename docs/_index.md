---
full: true
cascade:
  type: docs
breadcrumbs: false
comments: false
---

<div class="hx-mt-6 hx-mb-6">
{{< hextra/hero-headline >}}
  Advanced PostgreSQL & pgBouncer&nbsp;<br class="sm:hx-block hx-hidden" />Metrics Exporter for Prometheus
{{< /hextra/hero-headline >}}
</div>

<div class="hx-mb-12">
{{< hextra/hero-subtitle >}}
  The ultimate monitoring experience for PostgreSQL with **600+ metrics**,&nbsp;<br class="sm:hx-block hx-hidden" />
  **declarative configuration**, and **dynamic planning** capabilities
{{< /hextra/hero-subtitle >}}
</div>

<div class="hx-mb-6">
{{< hextra/hero-button text="Get Started" link="getting-started" >}}
</div>

<div class="hx-mt-6"></div>

{{< hextra/feature-grid >}}
  {{< hextra/feature-card
    title="Comprehensive Metrics"
    subtitle="Monitor PostgreSQL (10-18+) and pgBouncer (1.8-1.24+) with 600+ metrics and ~3K time series per instance"
    style="background: radial-gradient(ellipse at 50% 80%,rgba(142,53,74,0.15),hsla(0,0%,100%,0));"
  >}}
  {{< hextra/feature-card
  title="Declarative Configuration"
  subtitle="Define custom metrics through YAML configs with fine-grained control over timeout, caching, and skip conditions"
  style="background: radial-gradient(ellipse at 50% 80%,rgba(194,97,254,0.15),hsla(0,0%,100%,0));"
  >}}
  {{< hextra/feature-card
  title="Custom Collectors"
  subtitle="Define your own metrics with declarative YAML configuration and dynamic query planning"
  style="background: radial-gradient(ellipse at 50% 80%,rgba(255,119,168,0.15),hsla(0,0%,100%,0));"
  >}}
  {{< hextra/feature-card
    title="Auto-Discovery"
    subtitle="Automatically discover and monitor multiple databases within a PostgreSQL instance"
    style="background: radial-gradient(ellipse at 50% 80%,rgba(142,53,74,0.15),hsla(0,0%,100%,0));"
  >}}

  {{< hextra/feature-card
    title="Dynamic Planning"
    subtitle="Automatically adapt metric collection based on PostgreSQL version, extensions, and server characteristics"
    style="background: radial-gradient(ellipse at 50% 80%,rgba(221,210,59,0.15),hsla(0,0%,100%,0));"
  >}}
  {{< hextra/feature-card
  title="Production Ready"
  subtitle="Battle-tested in real-world environments across 12K+ cores for 6+ years with enterprise reliability"
  style="background: radial-gradient(ellipse at 50% 80%,rgba(120,119,198,0.15),hsla(0,0%,100%,0));"
  >}}
  {{< hextra/feature-card
    title="Health Check APIs"
    subtitle="Comprehensive HTTP endpoints for service health and traffic routing with primary/replica detection"
    style="background: radial-gradient(ellipse at 50% 80%,rgba(221,210,59,0.15),hsla(0,0%,100%,0));"
  >}}

  {{< hextra/feature-card
    title="Smart Caching"
    subtitle="Built-in caching mechanism with configurable TTL to reduce database load and improve performance"
    style="background: radial-gradient(ellipse at 50% 80%,rgba(194,97,254,0.15),hsla(0,0%,100%,0));"
  >}}

  {{< hextra/feature-card
    title="Extension Aware"
    subtitle="Native support for TimescaleDB, Citus, pg_stat_statements, pg_wait_sampling and automatic detection"
    style="background: radial-gradient(ellipse at 50% 80%,rgba(42,142,142,0.15),hsla(0,0%,100%,0));"
  >}}

{{< /hextra/feature-grid >}}


## Installation

PG Exporter provides multiple [install](/install) methods to fit your infrastructure:

{{< tabs items="Docker,YUM,APT,BIN,SRC" >}}
{{< tab >}}
```bash
docker run -d --name pg_exporter -p 9630:9630 -e PG_EXPORTER_URL="postgres://user:pass@host:5432/postgres" pgsty/pg_exporter:latest
```
{{< /tab >}}
{{< tab >}}
```bash
# RPM-based systems
sudo tee /etc/yum.repos.d/pigsty-infra.repo > /dev/null <<-'EOF'
[pigsty-infra]
name=Pigsty Infra for $basearch
baseurl=https://repo.pigsty.io/yum/infra/$basearch
enabled = 1
gpgcheck = 0
module_hotfixes=1
EOF

sudo yum makecache;
sudo yum install -y pg_exporter
```
{{< /tab >}}
{{< tab >}}
```bash
sudo tee /etc/apt/sources.list.d/pigsty-infra.list > /dev/null <<EOF
deb [trusted=yes] https://repo.pigsty.io/apt/infra generic main
EOF

sudo apt update;
sudo apt install -y pg-exporter
```
{{< /tab >}}
{{< tab >}}
```bash
wget https://github.com/pgsty/pg_exporter/releases/download/v1.1.1/pg_exporter-v1.1.1.linux-amd64.tar.gz
tar -xf pg_exporter-v1.1.1.linux-amd64.tar.gz
sudo install pg_exporter-v1.1.1.linux-amd64/pg_exporter /usr/bin/
sudo install pg_exporter-v1.1.1.linux-amd64/pg_exporter.yml /etc/pg_exporter.yml
```
{{< /tab >}}
{{< tab >}}
```bash
# Build from source
git clone https://github.com/pgsty/pg_exporter.git
cd pg_exporter
make build
```
{{< /tab >}}
{{< /tabs >}}


## Quick Start

Get PG Exporter up and running in minutes, [Getting Started](/getting-started) with:

```bash
# Run with PostgreSQL URL
PG_EXPORTER_URL='postgres://user:pass@localhost:5432/postgres' pg_exporter

# Access metrics
curl http://localhost:9630/metrics
```


## Documentation

<div class="hx-mt-6 hx-mb-6">
{{< cards >}}
{{< card link="start" title="Getting Started" subtitle="Quick start guide and basic concepts" >}}
{{< card link="install" title="Installation" subtitle="installation instructions for all platforms" icon="download" >}}
{{< card link="config" title="Configuration" subtitle="configuration reference and examples" >}}
{{< card link="deploy" title="Deployment" subtitle="Production deployment best practices" icon="server" >}}
{{< card link="https://github.com/pgsty/pg_exporter/tree/main/config" title="Collectors" subtitle="Detailed documentation for all metric collectors" icon="database" >}}
{{< /cards >}}
</div>

## Live Demo

Experience PG Exporter in action with our live demo environment:

<div class="hx-mt-4">
{{< hextra/hero-button text="View Live Demo" link="https://g.pgsty.com" >}}
</div>

The demo showcases real PostgreSQL clusters monitored by PG Exporter, featuring:
- Real-time metrics visualization with Grafana
- Multiple PostgreSQL versions and configurations
- Extension-specific metrics and monitoring
- Complete observability stack powered by [Pigsty](https://pgsty.com)

## Community & Support

<div class="hx-mt-6 hx-mb-6">
{{< cards >}}
{{< card link="https://github.com/pgsty/pg_exporter" title="GitHub" subtitle="Source code, issues, and contributions" icon="github" >}}
{{< card link="https://github.com/pgsty/pg_exporter/discussions" title="Discussions" subtitle="Ask questions and share experiences" icon="chat" >}}
{{< card link="https://pgsty.com" title="Pigsty" subtitle="Complete PostgreSQL Distro with PG Exporter" >}}
{{< /cards >}}
</div>

## License

PG Exporter is open-source software licensed under the [Apache License 2.0](https://github.com/pgsty/pg_exporter/blob/main/LICENSE).

Copyright 2018-2025 © [Ruohang Feng](https://vonng.com/en) / [rh@vonng.com](mailto:rh@vonng.com)