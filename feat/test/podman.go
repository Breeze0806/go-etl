// Copyright 2020 the go-etl Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package integration provides podman-based helpers for end-to-end integration
// tests of MySQL and PostgreSQL JSON support in go-etl. It exposes
// StartMySQLContainer, StartPostgresContainer and BuildGoETLBinary so that
// caller test files can spin up real database containers and a freshly built
// go-etl binary, and tear them down via the returned cleanup function.
//
// Build tag: every file in this package is gated behind the "integration"
// build tag, so it is excluded from the default `make test` and `make cover`
// builds and only compiled when explicitly requested, e.g.:
//
//	go test -tags integration ./...
//
// All helpers are safe to call from any working directory; callers that
// invoke the produced go-etl binary are responsible for os.Chdir-ing into a
// per-test temp directory so that the binary's CWD-relative log file
// (see cmd/datax/log.go) does not collide between tests.

//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// cmdTimeout is the maximum duration any single os/exec invocation in this
// file is allowed to run. Five minutes is generous enough for cold image
// pulls, database initialisation and full Go builds, while still failing
// tests in a bounded time when something goes wrong.
const cmdTimeout = 5 * time.Minute

// healthCheckTimeout is the maximum time we will wait for a container to
// become reachable after `podman run -d`. Thirty seconds matches typical
// MySQL/Postgres cold-start latencies in CI and local podman.
const healthCheckTimeout = 60 * time.Second

// healthCheckInterval is the poll cadence for the readiness probe above.
const healthCheckInterval = 2 * time.Second

// StartMySQLContainer launches a MySQL 8.0 container via podman, waits until
// it accepts connections, and returns a DSN suitable for database/sql along
// with a cleanup function that removes the container.
//
// The container name is unique per test (incorporates t.Name() and the
// current process PID) so parallel tests cannot collide. The host port is
// allocated dynamically with `podman run -P` and discovered with
// `podman port <name> 3306/tcp`, avoiding fixed port numbers.
//
// The returned cleanup is also registered with t.Cleanup so it still runs
// when the test fails. Callers may still invoke it directly if they need
// explicit teardown ordering.
func StartMySQLContainer(t testing.TB) (string, func()) {
	t.Helper()

	name := fmt.Sprintf("goetl-mysql-%s-%d", t.Name(), os.Getpid())

	runArgs := []string{
		"run", "-d", "--rm",
		"--name", name,
		"-P",
		"-e", "MYSQL_ROOT_PASSWORD=testpw",
		"-e", "MYSQL_DATABASE=goetl",
		"docker.io/library/mysql:8.0",
	}
	if err := runPodman(t, runArgs...); err != nil {
		t.Fatalf("StartMySQLContainer: podman run failed: %v", err)
	}

	port := discoverHostPort(t, name, "3306")
	dsn := fmt.Sprintf("root:testpw@tcp(127.0.0.1:%d)/goetl?parseTime=true&multiStatements=true", port)

	waitForMySQL(t, name)

	cleanup := func() { removeContainer(t, name) }
	t.Cleanup(cleanup)
	return dsn, cleanup
}

// StartPostgresContainer launches a PostgreSQL 14 container via podman,
// waits until pg_isready succeeds, and returns a DSN suitable for
// database/sql along with a cleanup function that removes the container.
//
// Like StartMySQLContainer, the container name is unique per test and the
// host port is discovered dynamically from `podman port`.
func StartPostgresContainer(t testing.TB) (string, func()) {
	t.Helper()

	name := fmt.Sprintf("goetl-pg-%s-%d", t.Name(), os.Getpid())

	runArgs := []string{
		"run", "-d", "--rm",
		"--name", name,
		"-P",
		"-e", "POSTGRES_PASSWORD=testpw",
		"-e", "POSTGRES_USER=goetl",
		"-e", "POSTGRES_DB=goetl",
		"docker.io/library/postgres:14",
	}
	if err := runPodman(t, runArgs...); err != nil {
		t.Fatalf("StartPostgresContainer: podman run failed: %v", err)
	}

	port := discoverHostPort(t, name, "5432")
	dsn := fmt.Sprintf("postgres://goetl:testpw@127.0.0.1:%d/goetl?sslmode=disable", port)

	waitForPostgres(t, name)

	cleanup := func() { removeContainer(t, name) }
	t.Cleanup(cleanup)
	return dsn, cleanup
}

// BuildGoETLBinary compiles ./cmd/datax into a per-test temporary location
// and returns the resulting binary path. The IGNORE_PACKAGES=db2 env var is
// injected into the build so the DB2 ODBC link step is skipped on machines
// where it has not been set up via `make dependencies`.
//
// The output path lives under t.TempDir() so it is automatically cleaned up
// by the testing framework at the end of the test.
func BuildGoETLBinary(t testing.TB) string {
	t.Helper()

	binPath := filepath.Join(t.TempDir(), "goetl")

	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "build", "-o", binPath, "./cmd/datax")
	// Append IGNORE_PACKAGES=db2 to the inherited environment so tools/datax/build/main.go
	// skips the DB2 plugin (and its CGO ODBC dependency) during `go generate` /
	// plugin registration. Without this the build fails on machines that have
	// not installed the IBM DB2 ODBC CLI driver.
	cmd.Env = append(os.Environ(), "IGNORE_PACKAGES=db2")

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("BuildGoETLBinary: go build failed: %v\n%s", err, string(out))
	}
	return binPath
}

// runPodman executes `podman <args...>` with a 5-minute timeout and returns
// a non-nil error if the command exits non-zero. Output (stdout+stderr) is
// included in the error so failure messages surface the underlying podman
// diagnostic.
func runPodman(t testing.TB, args ...string) error {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "podman", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// discoverHostPort runs `podman port <name> <port>/tcp` and returns the
// dynamically assigned host port. Output looks like `0.0.0.0:32789` or
// `[::]:32789`; we split on `:` and take the last field. The host part is
// discarded because integration tests always connect via 127.0.0.1.
func discoverHostPort(t testing.TB, name, containerPort string) int {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "podman", "port", name, containerPort+"/tcp")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("discoverHostPort(%s, %s): podman port failed: %v\n%s",
			name, containerPort, err, string(out))
	}
	line := strings.TrimSpace(strings.SplitN(string(out), "\n", 2)[0])
	parts := strings.Split(line, ":")
	if len(parts) < 2 {
		t.Fatalf("discoverHostPort(%s, %s): unexpected podman port output %q",
			name, containerPort, line)
	}
	var p int
	if _, err := fmt.Sscanf(parts[len(parts)-1], "%d", &p); err != nil {
		t.Fatalf("discoverHostPort(%s, %s): bad port in %q: %v",
			name, containerPort, line, err)
	}
	return p
}

// waitForMySQL polls `mysqladmin ping` inside the named container every
// 2s for up to 60s. It calls t.Fatalf if the container does not become
// ready within the budget; the last exec error is included so diagnosis is
// possible from the test log alone.
func waitForMySQL(t testing.TB, name string) {
	t.Helper()
	deadline := time.Now().Add(healthCheckTimeout)
	var lastErr error
	for time.Now().Before(deadline) {
		ctx, cancel := context.WithTimeout(context.Background(), healthCheckInterval)
		cmd := exec.CommandContext(ctx, "podman", "exec", name,
			"mysqladmin", "ping", "-h127.0.0.1", "-uroot", "-ptestpw", "--silent")
		err := cmd.Run()
		cancel()
		if err == nil {
			return
		}
		lastErr = err
		time.Sleep(healthCheckInterval)
	}
	t.Fatalf("waitForMySQL: container %s did not become ready within %s; last error: %v",
		name, healthCheckTimeout, lastErr)
}

// waitForPostgres polls `pg_isready` inside the named container every
// 2s for up to 60s. It mirrors waitForMySQL but uses the Postgres
// readiness probe.
func waitForPostgres(t testing.TB, name string) {
	t.Helper()
	deadline := time.Now().Add(healthCheckTimeout)
	var lastErr error
	for time.Now().Before(deadline) {
		ctx, cancel := context.WithTimeout(context.Background(), healthCheckInterval)
		cmd := exec.CommandContext(ctx, "podman", "exec", name,
			"pg_isready", "-U", "goetl")
		err := cmd.Run()
		cancel()
		if err == nil {
			return
		}
		lastErr = err
		time.Sleep(healthCheckInterval)
	}
	t.Fatalf("waitForPostgres: container %s did not become ready within %s; last error: %v",
		name, healthCheckTimeout, lastErr)
}

// removeContainer runs `podman rm -f <name>` to forcibly delete the
// container. Errors are reported via t.Logf rather than t.Fatalf because
// the test is already finishing and we do not want cleanup failures to
// mask the real test result.
func removeContainer(t testing.TB, name string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "podman", "rm", "-f", name)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Logf("removeContainer(%s): podman rm -f failed: %v\n%s", name, err, string(out))
	}
}
