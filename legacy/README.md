# Legacy Config Bundle (v1.19)

This directory contains the **legacy pg_exporter configuration bundle v1.19** for **PostgreSQL 9.1 - 9.6** (EOL).

- `pg_exporter.yml`: merged legacy config (ready to use)
- `config/`: separated collector definitions (source of truth)

## Generate / Update

From the repository root:

```bash
make conf9
```

## Usage

Use the merged legacy config file:

```bash
PG_EXPORTER_CONFIG=legacy/pg_exporter.yml \
PG_EXPORTER_URL='postgres://user:pass@host:5432/postgres' \
pg_exporter --auto-discovery --exclude-database=template0,template1
```

Or load separated collectors directly (directory mode):

```bash
PG_EXPORTER_CONFIG=legacy/config \
PG_EXPORTER_URL='postgres://user:pass@host:5432/postgres' \
pg_exporter --auto-discovery --exclude-database=template0,template1
```

## Notes

- PostgreSQL **9.0 is not supported**.
- For PostgreSQL **10+**, use the default config in repo root: `pg_exporter.yml`.
