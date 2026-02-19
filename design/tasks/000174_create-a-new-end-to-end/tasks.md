# Implementation Tasks

## Setup

- [x] Create directory `tests/e2e/manifests/` for e2e test manifests

## Keel Deployment Manifests

- [x] Create `tests/e2e/manifests/keel-rbac.yaml` with ServiceAccount, ClusterRole, and ClusterRoleBinding for keel
- [~] Create `tests/e2e/manifests/keel-deployment.yaml` with keel Deployment configured for polling (`POLL=1`, `POLL_DEFAULTSCHEDULE=@every 5s`)

## Test Workload Manifest

- [ ] Create `tests/e2e/manifests/test-deployment.yaml` with a deployment that has:
  - Label `keel.sh/policy: major`
  - Label `keel.sh/trigger: poll`
  - Annotation `keel.sh/pollSchedule: @every 5s`
  - Image `keelhq/push-workflow-example:0.1.0`

## Verification Script

- [ ] Create `tests/e2e/verify.sh` script that:
  - Waits for keel pod to be ready
  - Waits up to 120 seconds for test deployment image to update
  - Verifies image changed from `0.1.0` to `0.10.0`
  - Exits with code 0 on success, 1 on failure

## GitHub Actions Workflow

- [ ] Create `.github/workflows/e2e.yml` with:
  - Triggers: push to master, pull_request to master, workflow_dispatch
  - Job on `ubuntu-latest`
  - Step: Checkout code
  - Step: Set up Go 1.23
  - Step: Set up Docker Buildx
  - Step: Run `scripts/start-local-cluster.sh`
  - Step: Build keel Docker image (`docker build -t keel:e2e .`)
  - Step: Load image into kind (`kind load docker-image keel:e2e --name keel-dev`)
  - Step: Apply keel manifests (`kubectl apply -f tests/e2e/manifests/`)
  - Step: Wait for keel deployment ready
  - Step: Create test namespace and apply test deployment
  - Step: Run verification script
  - Step: Cleanup kind cluster (runs always, even on failure)

## Testing

- [ ] Test workflow locally using `act` or manual kind cluster setup
- [ ] Verify end-to-end flow: cluster up → keel deployed → test deployment updated → cleanup