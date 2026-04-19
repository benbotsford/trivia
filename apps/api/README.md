# api — database layer

Schema, migrations, and sqlc query sources for the trivia backend.

## Layout

```
apps/api/
  migrations/       # goose migrations, numbered 0001..NNNN
  queries/          # sqlc query sources (one file per entity)
  internal/store/   # sqlc-generated Go code (do not edit by hand)
  sqlc.yaml
```

## Prerequisites

Install goose and sqlc once per machine:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

The Go app depends on (to be added when `go.mod` lands):

- `github.com/jackc/pgx/v5`
- `github.com/google/uuid` (v1.6+ for `uuid.NewV7`)

## Primary keys

IDs are UUIDv7, **generated app-side** with `uuid.NewV7()` from `github.com/google/uuid`. No
database extension required. Every `INSERT` takes an explicit `id` parameter — the database
does not generate it. Rationale: keeps the schema portable across vanilla Postgres / Neon /
any Postgres-compatible host, and makes ID creation observable in application code.

## Running migrations

### Local Postgres

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/trivia?sslmode=disable"

# apply
goose -dir migrations postgres "$DATABASE_URL" up

# roll back one step
goose -dir migrations postgres "$DATABASE_URL" down

# status
goose -dir migrations postgres "$DATABASE_URL" status
```

### Neon

Set `DATABASE_URL` to the connection string for the **dev** branch (pooled or direct —
either works for migrations; direct is simpler). Neon requires `sslmode=require`:

```bash
export DATABASE_URL="postgres://<user>:<pass>@<host>/<db>?sslmode=require"
goose -dir migrations postgres "$DATABASE_URL" up
```

Promote the dev branch's schema to `main` with Neon's branch reset/merge flow once the
dev schema is stable.

## Generating Go code from sqlc

From `apps/api/`:

```bash
sqlc generate
```

Output lands in `internal/store/`. Commit the generated files.

## Schema notes

- `citext` is the only required extension; `CREATE EXTENSION IF NOT EXISTS citext` runs
  in migration `0001_init.sql`. Neon ships with citext available.
- An `updated_at` trigger function (`trigger_set_updated_at()`) is installed in `0001`
  and attached to tables that track mutation time.
- `games.code` is a 6-char `[A-Z0-9]` short code. Uniqueness is enforced only among
  **active** games (`lobby` / `in_progress`) via a partial unique index, so codes can
  be recycled after a game ends.
- `questions.choices` is JSONB and only required when `type = 'multiple_choice'`.
  A CHECK constraint enforces the shape (array, 2–10 elements).
- `subscriptions` is intentionally empty. The billing system ships later; until then
  `EntitlementChecker` should always return `true` (see `PROJECT.md`).
- Nullable `citext`/`text`/`jsonb` columns come through sqlc's pgx/v5 driver as
  `pgtype.Text` / `pgtype.JSONB`. Wrap these at the service boundary if you want
  ergonomic `*string` / typed struct APIs.

## Adding a new migration

```bash
goose -dir migrations create <name> sql
```

Then edit the generated file to include both `-- +goose Up` and `-- +goose Down`
sections. Prefer one logical change per migration (new table, new index, new column).
After the migration is in, regenerate with `sqlc generate`.
