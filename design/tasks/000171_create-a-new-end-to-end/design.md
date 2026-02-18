# Design: End-to-End Test with Kind Cluster in GitHub Actions

## Architecture Overview

The e2e test will run as a new GitHub Actions workflow that creates a kind cluster, builds Keel, and runs the existing acceptance tests against the cluster.

```
┌─────────────────────────────────────────────────────────┐
│                   GitHub Actions                        │
│  ┌───────────────────────────────────────────────────┐  │
│  │  1. Checkout code                                 │  │
│  │  2. Set up Go                                     │  │
│  │  3. Run start-local-cluster.sh (creates kind)    │  │
│  │  4. Build keel binary (make install)             │  │
│  │  5. Run e2e tests (make e2e)                     │  │
│  │  6. Cleanup (kind deletes on job completion)     │  │
│  └───────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

## Key Decisions

### Decision 1: Reuse Existing Test Infrastructure
**Choice:** Use the existing `tests/` Go tests rather than writing new tests.

**Rationale:** The project already has comprehensive acceptance tests in `tests/acceptance_test.go` and `tests/acceptance_polling_test.go` that:
- Create test deployments with Keel policy labels
- Trigger updates via webhooks or polling
- Verify deployments are updated correctly

### Decision 2: Use Webhook-Based Tests
**Choice:** Run the webhook-based tests (not polling tests) for faster CI.

**Rationale:** Polling tests require waiting for poll intervals (2+ seconds per test). Webhook tests are faster and more deterministic.

### Decision 3: Separate Workflow File
**Choice:** Create a new `e2e.yml` workflow instead of adding to `ci.yml`.

**Rationale:** 
- E2E tests are slower and more resource-intensive
- Allows running e2e tests independently
- Keeps CI workflow fast for quick feedback

## Existing Patterns Discovered

From `keel/tests/helpers.go`:
- `KeelCmd` struct manages starting/stopping keel with `--no-incluster` flag
- Tests use `getKubeConfig()` to connect via `~/.kube/config`
- `createNamespaceForTest()` creates isolated test namespaces
- `waitFor()` polls deployments to verify image updates

From `keel/Makefile`:
- `make install` builds and installs the keel binary
- `make e2e` runs: `cd tests && go test`

## Workflow Structure

```yaml
# .github/workflows/e2e.yml
name: E2E Tests
on: [push, pull_request]
jobs:
  e2e:
    runs-on: ubuntu-latest
    steps:
      - Checkout
      - Setup Go 1.23
      - Run start-local-cluster.sh
      - make install
      - make e2e (with timeout)
```

## Constraints

1. **Docker required** - Kind runs Kubernetes nodes as Docker containers
2. **Kubectl context** - Must use `kind-keel-dev` context set by the script
3. **Test timeout** - E2E tests need sufficient time (~5-10 minutes)
4. **No sudo in GHA** - The script uses sudo; GHA runners allow this by default