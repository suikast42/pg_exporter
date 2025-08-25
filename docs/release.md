---
title: Release Notes
weight: 9
breadcrumbs: false
comments: false
---

The latest stable version of `pg_exporter` is [v1.0.2](https://github.com/pgsty/pg_exporter/releases/tag/v1.0.2)

|     Version     |    Date    | Summary                                                 |                               GitHub                               |
|:---------------:|:----------:|---------------------------------------------------------|:------------------------------------------------------------------:|
| [v1.0.2](#v102) | 2023-08-14 | Build for more os arch with goreleaser                  | [v1.0.2](https://github.com/pgsty/pg_exporter/releases/tag/v1.0.2) |
| [v1.0.1](#v101) | 2023-07-17 | DockerHub images, Go 1.24.5, disable pg_tsdb_hypertable | [v1.0.1](https://github.com/pgsty/pg_exporter/releases/tag/v1.0.1) |
| [v1.0.0](#v100) | 2023-05-06 | PostgreSQL 18 support, new WAL/checkpointer/I/O metrics | [v1.0.0](https://github.com/pgsty/pg_exporter/releases/tag/v1.0.0) |
| [v0.9.0](#v090) | 2023-04-26 | TimescaleDB, Citus, pg_wait_sampling collectors         | [v0.9.0](https://github.com/pgsty/pg_exporter/releases/tag/v0.9.0) |
| [v0.8.1](#v081) | 2024-02-14 | Dependencies update, docker image tags                  | [v0.8.1](https://github.com/pgsty/pg_exporter/releases/tag/v0.8.1) |
| [v0.8.0](#v080) | 2024-02-14 | PgBouncer 1.24 support, Go 1.24, logging refactor       | [v0.8.0](https://github.com/pgsty/pg_exporter/releases/tag/v0.8.0) |
| [v0.7.1](#v071) | 2022-08-13 | Routine update, configuration as Reader support         | [v0.7.1](https://github.com/pgsty/pg_exporter/releases/tag/v0.7.1) |
| [v0.7.0](#v070) | 2022-08-13 | PostgreSQL 17 support, predicate queries feature        | [v0.7.0](https://github.com/pgsty/pg_exporter/releases/tag/v0.7.0) |
| [v0.6.0](#v060) | 2022-10-18 | PostgreSQL 16 support, ARM64 packages, security fixes   | [v0.6.0](https://github.com/pgsty/pg_exporter/releases/tag/v0.6.0) |
| [v0.5.0](#v050) | 2022-04-27 | RPM/DEB builds, column scaling, metrics enhancements    | [v0.5.0](https://github.com/pgsty/pg_exporter/releases/tag/v0.5.0) |
| [v0.4.1](#v041) | 2022-03-08 | Collector updates, connect-timeout parameter            | [v0.4.1](https://github.com/pgsty/pg_exporter/releases/tag/v0.4.1) |
| [v0.4.0](#v040) | 2021-07-12 | PostgreSQL 14 support, auto-discovery feature           | [v0.4.0](https://github.com/pgsty/pg_exporter/releases/tag/v0.4.0) |
| [v0.3.2](#v032) | 2021-02-01 | Shadow DSN fixes, documentation updates                 | [v0.3.2](https://github.com/pgsty/pg_exporter/releases/tag/v0.3.2) |
| [v0.3.1](#v031) | 2020-12-04 | Configuration fixes for older PostgreSQL versions       | [v0.3.1](https://github.com/pgsty/pg_exporter/releases/tag/v0.3.1) |
| [v0.3.0](#v030) | 2020-10-29 | PostgreSQL 13 support, REST APIs, dummy server          | [v0.3.0](https://github.com/pgsty/pg_exporter/releases/tag/v0.3.0) |
| [v0.2.0](#v020) | 2020-03-21 | YUM packages, configuration reload support              | [v0.2.0](https://github.com/pgsty/pg_exporter/releases/tag/v0.2.0) |
| [v0.1.2](#v012) | 2020-02-20 | Dynamic configuration reload, bulky mode                | [v0.1.2](https://github.com/pgsty/pg_exporter/releases/tag/v0.1.2) |
| [v0.1.1](#v011) | 2020-01-10 | Startup hang bug fix                                    | [v0.1.1](https://github.com/pgsty/pg_exporter/releases/tag/v0.1.1) |
| [v0.1.0](#v010) | 2020-01-08 | Initial stable release                                  | [v0.1.0](https://github.com/pgsty/pg_exporter/releases/tag/v0.1.0) |
| [v0.0.4](#v004) | 2019-12-20 | Production tested release                               | [v0.0.4](https://github.com/pgsty/pg_exporter/releases/tag/v0.0.4) |
| [v0.0.3](#v003) | 2019-12-14 | Production environment testing                          | [v0.0.3](https://github.com/pgsty/pg_exporter/releases/tag/v0.0.3) |
| [v0.0.2](#v002) | 2019-12-09 | Early testing release                                   | [v0.0.2](https://github.com/pgsty/pg_exporter/releases/tag/v0.0.2) |
| [v0.0.1](#v001) | 2019-12-06 | Initial release with PgBouncer mode                     | [v0.0.1](https://github.com/pgsty/pg_exporter/releases/tag/v0.0.1) |



## v1.0.2

- Build with Go 1.25.0 and latest dependencies
- Dedicate website and homepage: https://exp.pgsty.com
- Release with goreleaser for more os/arch:
  - add windows amd64 support
  - add linux ppc64le support

**Checksums**

```bash
```

https://github.com/pgsty/pg_exporter/releases/tag/v1.0.2



## v1.0.1

- Add dockerhub images: [pgsty/pg_exporter](https://hub.docker.com/r/pgsty/pg_exporter)
- Bump go dependencies to the latest version, build with go 1.24.5
- Disable `pg_tsdb_hypertable` collector by default, since `timescaledb` catalog is changed.

**Checksums**

```bash
d5e2d6a656eef0ae1b29cd49695f9773  pg_exporter-1.0.1-1.aarch64.rpm
cb01bb78d7b216a235363e9342803cb3  pg_exporter-1.0.1-1.x86_64.rpm
67093a756b04845f69ad333b6d458e81  pg_exporter-v1.0.1.darwin-amd64.tar.gz
2d3fdc10045d1cf494b9c1ee7f94f127  pg_exporter-v1.0.1.darwin-arm64.tar.gz
e242314461becfa99c3978ae72838ab0  pg_exporter-v1.0.1.linux-amd64.tar.gz
63de91da9ef711a53718bc60b89c82a6  pg_exporter-v1.0.1.linux-arm64.tar.gz
718f6afc004089f12c1ca6553f9b9ba5  pg-exporter_1.0.1_amd64.deb
57da7a8005cdf91ba8c1fb348e0d7367  pg-exporter_1.0.1_arm64.deb
```

https://github.com/pgsty/pg_exporter/releases/tag/v1.0.1




## v1.0.0

Add PostgreSQL 18 metrics support

- new collector branch `pg_wal_18`:
- remove `write`, `sync`, `write_time`, `sync_time` metrics
- move to `pg_stat_io`
- new collector branch `pg_checkpointer_18`:
- new metric `num_done`
- new metric `slru_written`
- new collector branch `pg_db_18`:
- new metric `parallel_workers_to_launch`
- new metric `parallel_workers_launched`
- new collector branch `pg_table_18`:
- `table_parallel_workers_to_launch`
- `table_parallel_workers_launched`
- new collector branch `pg_io_18`:
- new series about WAL statistics
- new metric `read_bytes`
- new metric `write_bytes`
- new metric `extend_bytes`
- remove `op_bytes` due to fixed value
- new collector branch `pg_vacuuming_18`
- new metric `delay_time`

```bash
8637bc1a05b93eedfbfd3816cca468dd  pg_exporter-1.0.0-1.aarch64.rpm
a28c4c0dcdd3bf412268a2dbff79f5b9  pg_exporter-1.0.0-1.x86_64.rpm
229129209b8e6bc356c28043c7c22359  pg_exporter-v1.0.0.darwin-amd64.tar.gz
d941c2c28301269e62a8853c93facf12  pg_exporter-v1.0.0.darwin-arm64.tar.gz
5bbb94db46cacca4075d4c341c54db37  pg_exporter-v1.0.0.linux-amd64.tar.gz
da9ad428a50546a507a542d808f1c0fa  pg_exporter-v1.0.0.linux-arm64.tar.gz
0fa2395d9d7a43ab87e5c87e5b06ffcc  pg-exporter_1.0.0_amd64.deb
fed56f8a37e30cc59e85f03c81fce3f5  pg-exporter_1.0.0_arm64.deb
```

https://github.com/pgsty/pg_exporter/releases/tag/v1.0.0




## v0.9.0

**Default Collectors**

* new metrics collector for `timescaledb` hypertable
* new metrics collector for `citus` dist node
* new metrics collector for `pg_wait_sampling` wait event profile
* `pg_slot` overhaul: Add 16/17 pg_replication_slot metrics
* allow `pg_slot` collector run on replica since 16/17
* refactor `pg_wait` collector to agg from all processes
* restrict pg_clustering, pg_indexing, pg_vacuuming run on primary
* mark all `reset_time` as `GAUGE` rather than `COUNTER`
* fix `pg_recovery_prefetch_skip_fpw` type from `GAUGE` to `COUNTER`
* fix `pg_recv.state` type from `LABEL` to `GAUGE`
* Format collector in compact mode
* new default metric `pg_exporter_build_info` / `pgbouncer_exporter_build_info`
* add `server_encoding` to `pg_meta` collector
* add 12 new setting metrics to `pg_setting` collector
- wal_block_size
- segment_size
- wal_segment_size
- wal_level
- wal_log_hints
- work_mem
- hugepage_count
- hugepage_status
- max_wal_size
- min_wal_size
- max_slot_wal_keep_size

**Exporter Codebase**

* normalize collector branch name with min pg ver suffix
* Add license file to binary packages
* move `pgsty/pg_exporter` repo to `pgsty/pg_exporter`
* refactor `server.go` to reduce `Compatible` and `PostgresPrecheck` complexity
* rename metrics collector with extra number prefix for better sorting
* bump dependencies to the latest version
* execute fatal collectors ahead of all non-fatal collectors, and fail fast

https://github.com/pgsty/pg_exporter/releases/tag/v0.9.0





## v0.8.1

* Bump dependencies to the latest version
* [Bump golang.org/x/net from 0.35.0 to 0.36.0 #67](https://github.com/pgsty/pg_exporter/pull/67)
* Update docker images building tags

https://github.com/pgsty/pg_exporter/releases/tag/v0.8.1






## v0.8.0

* Add PgBouncer 1.24 new metrics support (stat, pool, database)
* Fix: `310-pg_size.yml` fails if log dir not set properly [#64](https://github.com/pgsty/pg_exporter/issues/64) by [@Süleyman Vurucu](https://github.com/suikast42)
* Build with the latest Go 1.24 and bump all dependencies
* Refactor logging with the standard `log/slog` instead of `go-kit`
* **Full Changelog**: https://github.com/pgsty/pg_exporter/compare/v0.7.1...v0.8.0

https://github.com/pgsty/pg_exporter/releases/tag/v0.8.0





## v0.7.1

Routine update with dependabot

* Feat: support specifying configuration as Reader by [@ringerc](https://github.com/ringerc) in [#62](https://github.com/pgsty/pg_exporter/pull/62)
* Bump golang.org/x/crypto from 0.21.0 to 0.31.0 by [@dependabot](https://github.com/dependabot) in [#63](https://github.com/pgsty/pg_exporter/pull/63)
* Fix some typos
* **Full Changelog**: https://github.com/pgsty/pg_exporter/compare/v0.7.0...v0.7.1

https://github.com/pgsty/pg_exporter/releases/tag/v0.7.1





## v0.7.0

Refactor codebase for the latest go version.

- [PostgreSQL 17 Metrics Support](https://github.com/pgsty/pg_exporter/issues/53) by @Vonng
- [pg_exporter: predicate queries feature](https://github.com/pgsty/pg_exporter/pull/47) by [@ringerc](https://github.com/ringerc)
- [Do a clean build in the dockerfile](https://github.com/pgsty/pg_exporter/pull/54) by [@ringerc](https://github.com/ringerc) by [@ringerc](https://github.com/ringerc)
- [pg_exporter: don't panic after "bind: address already in use"](https://github.com/pgsty/pg_exporter/pull/46) by [@ringerc](https://github.com/ringerc)
- [pg_exporter: fix /stat endpoint formatting](https://github.com/pgsty/pg_exporter/pull/48) by [@ringerc](https://github.com/ringerc)
- [pg_exporter: omit default query properties on yaml export](https://github.com/pgsty/pg_exporter/pull/49) by [@ringerc](https://github.com/ringerc)
- [Exclude template DBs from discovery and schema-qualify discovery query](https://github.com/pgsty/pg_exporter/pull/50) by [@ringerc](https://github.com/ringerc)
- [Fix some typos and some metric description mistakes](https://github.com/pgsty/pg_exporter/pull/51) by [@ringerc](https://github.com/ringerc)
- [Switch from unmaintained lib/pq driver to pgx with stdlib wrapper](https://github.com/pgsty/pg_exporter/pull/52) by [@ringerc](https://github.com/ringerc)

https://github.com/pgsty/pg_exporter/releases/tag/v0.7.0




## v0.6.0

- Security Enhancement: Fix [security](https://github.com/pgsty/pg_exporter/security/dependabot?q=is%3Aclosed)
  dependent-bot issue
- Add pg16 collectors
- Add `arm64` & `aarch64` packages
- Remove the `monitor` schema requirement for `pg_query` collectors (you have to ensure it with search_path or just
  install `pg_stat_statements` in the default `public` schema)
- Fix pgbouncer version parsing message level from info to debug
- Fix `pg_table_10_12` collector missing `relid` issue.

- [Recognize the files with yml suffix in config directory](https://github.com/pgsty/pg_exporter/pull/34) by [@Japin Li](https://github.com/japinli)
- [Support PostgreSQL 15 and higher](https://github.com/pgsty/pg_exporter/pull/35) by [@Japin Li](https://github.com/japinli)
- [Fix connect-timeout propagation](https://github.com/pgsty/pg_exporter/pull/37/files) by [@mouchar](https://github.com/mouchar)

https://github.com/pgsty/pg_exporter/releases/tag/v0.6.0





## v0.5.0

**Exporter Enhancement**

- Build rpm & deb with `nfpm`
- Add `column.default`, replace when metric value is NULL
- Add `column.scale`, multiply scale factor when metric value is float/int (e.g µs to second)
- Fix `/stat` endpoint output
- Add docker container [`pgsty/pg_exporter`](https://hub.docker.com/r/pgsty/pg_exporter)

**Metrics Collector**

- scale bgwriter & pg_wal time unit to second
- remove pg_class collector and move it to pg_table & pg_inex
- add pg_class metrics to pg_table
- add pg_class metrics to pg_index
- enable pg_table_size by default
- scale pg_query pg_db pg_bgwriter pg_ssl pgbouncer_stat time metrics to second

https://github.com/pgsty/pg_exporter/releases/tag/v0.5.0








## v0.4.1

- update default collectors
    - omit citus & timescaledb schemas on object monitoring
    - avoid duplicate pg_statio tuples
    - support pgbouncer v1.16
    - bug fix: `pg_repl` collector overlap on pg 12
- new parameter: `-T` `connect-timeout` `PG_EXPORTER_CONNECT_TIMEOUT`
  this can be useful when monitoring remote Postgres instances.
- now `pg_exporter.yaml` are renamed as `pg_exporter.yml` in rpm package.

https://github.com/pgsty/pg_exporter/releases/tag/v0.4.1








## v0.4.0

- Add PG 14 support
- Default metrics configuration overhaul. (BUT you can still use the old configuration)
- add `auto-discovery` , `include-database` and `exclude-database` option
- Add multiple database monitoring implementations (with `auto-discovery = on`)

https://github.com/pgsty/pg_exporter/releases/tag/v0.4.0







## v0.3.2

- fix shadow DSN corner case
- fix typo & docs

https://github.com/pgsty/pg_exporter/releases/tag/v0.3.2







## v0.3.1

fix default configuration problems (especially for versions lower than 13)

- setting `primary_conninfo` not exists until PG13
- add `funcid` label to `pg_func` collector to avoid func name duplicate label
- fix version string to `pg_exporter`

https://github.com/pgsty/pg_exporter/releases/tag/v0.3.1








## v0.3.0

https://github.com/pgsty/pg_exporter/releases/tag/v0.3.0

- Change default configuration, Support PostgreSQL 13 new metrics (`pg_slru`, `pg_shmem`, `pg_query13`,`pg_backup`,
  etc...)
- Add a series of new REST APIs for health / recovery status check
- Add a dummy server with fake `pg_up 0` metric, which serves before PgExporter is initialized.
- Add `sslmode=disable` to URL if `sslmode` is not given
- fix typos and bugs



## v0.2.0

- add yum package and linux service definition
- add a 'skip' flag into query config
- fix `pgbouncer_up` metrics
- add conf reload support

https://github.com/pgsty/pg_exporter/releases/tag/v0.2.0








## v0.1.2

* fix pgbouncer_up metrics
* add dynamic configuration reload
* remove 'shard' related logic
* add a 'bulky' mode to default settings

https://github.com/pgsty/pg_exporter/releases/tag/v0.1.2






## v0.1.1

Fix the bug that pg_exporter will hang during start-up if any query is failed.

https://github.com/pgsty/pg_exporter/releases/tag/v0.1.1





## v0.1.0

It works, looks good to me.

https://github.com/pgsty/pg_exporter/releases/tag/v0.1.0




## v0.0.4

Tested in real world production environment with 200+ nodes for about 2 weeks. Looks good !

https://github.com/pgsty/pg_exporter/releases/tag/v0.0.4





## v0.0.3

v0.0.3 Release, Tested in Production Environment

This version is already tested in a production environment.

This project is still under rapid evolution, I would say if you want use it in production , try with caution.

https://github.com/pgsty/pg_exporter/releases/tag/v0.0.3





## v0.0.2

It's ok to try now

https://github.com/pgsty/pg_exporter/releases/tag/v0.0.2





## v0.0.1

Add pgbouncer mode

https://github.com/pgsty/pg_exporter/releases/tag/v0.0.1