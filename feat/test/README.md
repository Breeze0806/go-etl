# Integration Tests

## Purpose

This package contains real-database integration tests for the MySQL JSON
column type and the PostgreSQL JSON and JSONB column types. Each test spins
up a disposable MySQL or Postgres container, runs the `go-etl` binary
(`cmd/datax/main.go`) against it, and verifies that JSON-shaped values
round-trip cleanly through the Reader and Writer plugins.

These tests are not unit tests. They exist to catch the kind of bug that
only shows up against a real server, and they are deliberately gated behind
a build tag so they never run by accident on a developer machine or in
regular CI.

## Prerequisites

Before running the tests, make sure the following are in place.

- **podman** installed and on `PATH`. Verify with `podman --version`.
  - Fedora / RHEL: `sudo dnf install podman`
  - Debian / Ubuntu: `sudo apt install podman`
  - macOS: `brew install podman`
  - The podman machine must be running on macOS: `podman machine start`.
- **Go 1.20+** toolchain. Verify with `go version`.
- **CGO toolchain (gcc)**. The `go-etl` binary depends on the `godror`
  (Oracle) and `mattn/go-sqlite3` CGO drivers, so a working C compiler is
  required. `gcc` 4.8 or newer on Linux, or Xcode Command Line Tools on
  macOS (`xcode-select --install`).
- **Either** `make dependencies` has been run once, **or** the
  `IGNORE_PACKAGES=db2` environment variable is exported for the test
  build. **This is the single most common reason a fresh box fails to
  build the test binary**: without `IGNORE_PACKAGES=db2`, the build pulls
  in the DB2 ODBC driver and fails to link on a machine that has not run
  `make dependencies`.
- **make** is available for the `//license-header` check the lint target
  runs on generated files. The integration tests themselves do not need
  `make`, only the build prerequisites above do.

## How to run

From the repository root:

```bash
# From repo root
export IGNORE_PACKAGES=db2
go test -tags integration -v ./feat/test/...
```

The `integration` build tag is required because every test file in this
package starts with `//go:build integration`. The package itself is named
`integration` (not `test`), and the tag matches the package name on
purpose so that a single `-tags integration` flag is enough to flip the
whole directory on. Plain `go test ./...` will skip these files entirely.

The command above builds `cmd/datax/main.go` as a test helper, generates
a temporary `config.json` for each test, runs the binary in a child
process, and tails the log. There is no separate build step to run.

## What gets tested

- **MySQL JSON column**
  - simple object, JSON array, nested object
  - escaped characters inside strings
  - SQL `NULL` and the JSON `null` literal (the two must not be confused)
  - multi-column rows where the JSON column sits next to scalar columns
- **PostgreSQL JSON column**
  - same matrix as MySQL: simple object, array, nested, escaped chars,
    SQL `NULL`, JSON `null` literal, multi-column rows
- **PostgreSQL jsonb column**
  - the same matrix as the JSON column
  - plus a `writeMode=copyIn` round-trip, which exercises the path where
    the writer converts a `[]byte` payload to a `string` before handing it
    to libpq's `COPY FROM STDIN` protocol. This is the workaround for
    the `jsonb` driver refusing raw `[]byte` on the wire, and it is the
    one configuration that tends to break in subtle ways.

## Container lifecycle

Each test is self-contained.

- The test starts its own MySQL or Postgres container, named
  `goetl-mysql-<test-name>-<pid>` or `goetl-pg-<test-name>-<pid>`. The
  `<pid>` suffix makes it obvious which process owns the container if a
  previous run leaked one.
- The host port is assigned at random via `podman run -P` so two tests
  can run side by side without colliding. The test then calls
  `podman port` to discover the mapping it just got and points the
  generated `config.json` at it.
- A `t.Cleanup` hook runs at the end of the test, including on failure.
  It calls `podman rm -f` against the named container, so a failed test
  does not leave a zombie process behind.

You can list and inspect the running containers at any time with
`podman ps -a`.

## Log location

`cmd/datax/log.go` opens `go-etl.log` in the current working directory,
and the current working directory is whatever the test set before
launching the binary. Each test does an `os.Chdir` into its
`t.TempDir()` before invoking the binary, so `go-etl.log` ends up under
the temp dir, not in the repo root.

The easiest way to read it after a run is:

```bash
cat "$(go env TMPDIR)/goetl-*/go-etl.log"
```

The directory name has a `goetl-` prefix because `t.TempDir()` builds
its path from the test name. If you want the log from a specific
container instead, `podman logs <container-name>` works too and includes
any output the database server wrote while the test was running.

## CI

These tests are **not** run in CI today. They need a real container
runtime and outbound image pulls, which the existing `Build.yml`
workflow does not provide, and the `mysql` and `postgres` containers
would push the CI runtime past its budget. The plan is to add a
separate `Integration.yml` workflow that uses the `setup-podman` action
and runs `go test -tags integration -v ./feat/test/...` on demand, but
that is out of scope for v1. Until then, integration coverage depends
on contributors running the suite locally before opening a PR that
touches the MySQL or PostgreSQL JSON code paths.

## Troubleshooting

- `podman: command not found`. Install podman using the command in the
  Prerequisites section, or use your distro's package manager. On
  macOS, also run `podman machine start`.
- `build failed: db2 ODBC` (or any link error mentioning `clidriver` or
  `sqlcli`). You forgot to export `IGNORE_PACKAGES=db2`. Set it and try
  again, or run `make dependencies` once to install the IBM DB2 ODBC
  CLI driver and let the default build succeed.
- `health check timeout` while the test waits for the database to accept
  connections. The container image is not in the local podman cache, so
  podman is trying to pull it on the fly. Pull both images manually
  first to warm the cache:

  ```bash
  podman pull docker.io/library/mysql:8.0
  podman pull docker.io/library/postgres:14
  ```

  Then re-run the test. On a slow link the first pull can take a few
  minutes; subsequent runs are instant.
- `port already in use`. The tests use `podman run -P`, which asks
  podman for a random free host port, so a port collision means a
  previous container leaked. Run `podman ps -a` to see the leftovers
  and `podman rm -f` to clean them up. The test's `t.Cleanup` should
  have done this for you, so a leak usually points at an interrupted
  previous run (Ctrl-C, IDE abort, OOM kill).
