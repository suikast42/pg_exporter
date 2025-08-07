---
title: Installation
description: How to download and install the pg_exporter
weight: 2
breadcrumbs: false
comments: false
---

PG Exporter provides multiple installation methods to suit different deployment scenarios. 
This guide covers all available installation options with detailed instructions for each platform.




## Pigsty

The easiest way to get started with `pg_exporter` is to use [Pigsty](https://pigsty.io), 
which is a complete PostgreSQL distribution with built-in Observability best practices based on `pg_exporter`, Prometheus, and Grafana.
You don't even need to know any details about `pg_exporter`, it just gives you all the metrics and dashboard panels

```bash
curl -fsSL https://repo.pigsty.io/get | bash; cd ~/pigsty;
```


## Release

You can also download `pg_exporter` package (`RPM`/`DEB`/ Tarball) directly from the [Latest GitHub Release Page](https://github.com/pgsty/pg_exporter/releases/latest):

{{< filetree/container >}}
{{< filetree/folder name="v1.0.1" state="open" >}}
{{< filetree/file name="https://github.com/pgsty/pg_exporter/releases/download/v1.0.1/pg-exporter_1.0.1_amd64.deb" >}}
{{< filetree/file name="https://github.com/pgsty/pg_exporter/releases/download/v1.0.1/pg-exporter_1.0.1_arm64.deb" >}}
{{< filetree/file name="https://github.com/pgsty/pg_exporter/releases/download/v1.0.1/pg_exporter-1.0.1-1.aarch64.rpm" >}}
{{< filetree/file name="https://github.com/pgsty/pg_exporter/releases/download/v1.0.1/pg_exporter-1.0.1-1.x86_64.rpm" >}}
{{< filetree/file name="https://github.com/pgsty/pg_exporter/releases/download/v1.0.1/pg_exporter-v1.0.1.linux-amd64.tar.gz" >}}
{{< filetree/file name="https://github.com/pgsty/pg_exporter/releases/download/v1.0.1/pg_exporter-v1.0.1.linux-arm64.tar.gz" >}}
{{< filetree/file name="https://github.com/pgsty/pg_exporter/releases/download/v1.0.1/pg_exporter-v1.0.1.darwin-amd64.tar.gz" >}}
{{< filetree/file name="https://github.com/pgsty/pg_exporter/releases/download/v1.0.1/pg_exporter-v1.0.1.darwin-arm64.tar.gz" >}}
{{< /filetree/folder >}}
{{< /filetree/container >}}

You can install it directly with your OS package manager (`rpm`/`dpkg`), or even put the binary in your `$PATH`.


## Repository

The pig package is also available in the [`pigsty-infra`](https://ext.pgsty.com/repo/infra) repo,
You can add the repo to your system, and install it with OS package manager:

### YUM

For EL distribution such as  RHEL，RockyLinux，CentOS，Alma Linux，OracleLinux,...:

```bash tab="yum"
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

### APT

For Debian, Ubuntu and compatible Linux Distributions:

```bash tab="apt"
sudo tee /etc/apt/sources.list.d/pigsty-infra.list > /dev/null <<EOF
deb [trusted=yes] https://repo.pigsty.io/apt/infra generic main
EOF

sudo apt update;
sudo apt install -y pg-exporter
```


## Docker

We have prebuilt docker images for `amd64` and `arm64` architectures on docker hub: [pgsty/pg_exporter](https://hub.docker.com/r/pgsty/pg_exporter).

```bash
# Basic usage
docker run -d \
  --name pg_exporter \
  -p 9630:9630 \
  -e PG_EXPORTER_URL="postgres://user:password@host:5432/postgres" \
  pgsty/pg_exporter:latest

# With custom configuration
docker run -d \
  --name pg_exporter \
  -p 9630:9630 \
  -v /path/to/pg_exporter.yml:/etc/pg_exporter.yml:ro \
  -e PG_EXPORTER_CONFIG="/etc/pg_exporter.yml" \
  -e PG_EXPORTER_URL="postgres://user:password@host:5432/postgres" \
  pgsty/pg_exporter:latest

# With auto-discovery enabled
docker run -d \
  --name pg_exporter \
  -p 9630:9630 \
  -e PG_EXPORTER_URL="postgres://user:password@host:5432/postgres" \
  -e PG_EXPORTER_AUTO_DISCOVERY="true" \
  -e PG_EXPORTER_EXCLUDE_DATABASE="template0,template1" \
  pgsty/pg_exporter:latest
```


## Binary

The `pg_exporter` can be installed as a standalone binary.


## Compatibility

The current pg_exporter support PostgreSQL version 10 and above.
While it is designed to work with any PostgreSQL major version (back to 9.x).

The only problem to use with legacy version (9.6 and below) is that 
we removed older metrics collector branches definition due to EOL.

You can always retrieve these legacy version of config files and use against historic versions of PostgreSQL

| PostgreSQL Version | Support Status |
|--------------------|----------------|
| 10 ~ 17            | ✅ Full Support |
| 9.6-               | ⚠️ Legacy Conf |

pg_exporter works with pgbouncer 1.8+, Since `v1.8` is the first version with `SHOW` command support.

| pgBouncer Version | Support Status |
|-------------------|----------------|
| 1.8.x ~ 1.24.x    | ✅ Full Support |
| before 1.8.x      | ⚠️ No Metrics  |
