# Design: End-to-End Test for Keel in GitHub Actions

## Architecture Overview

The e2e test creates a local Kubernetes cluster using kind, deploys Keel, creates a test deployment with Keel annotations, and verifies that Keel automatically updates the deployment when a newer image version is detected via polling.

```
┌─────────────────────────────────────────────────────────────────┐
│                    GitHub Actions Runner                         │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                    Kind Cluster (keel-dev)                 │  │
│  │  ┌─────────────────┐    ┌─────────────────────────────┐   │  │
│  │  │  Keel Pod       │───▶│  Test Deployment            │   │  │
│  │  │  (polling mode) │    │  (keelhq/push-workflow-     │   │  │
│  │  │                 │    │   example:0.1.0 → 0.10.0)   │   │  │
│  │  └────────┬────────┘    └─────────────────────────────┘   │  │
│  │           │                                                │  │
│  │           │ polls DockerHub for newer versions             │  │
│  │           ▼                                                │  │
│  │  ┌─────────────────┐                                       │  │
│  │  │  DockerHub      │                                       │  │
│  │  │  (public)       │                                       │  │
│  │  └─────────────────┘                                       │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## Key Decisions

### Decision 1: Use Existing Polling-Based E2E Pattern
**Choice:** Follow the existing `tests/acceptance_polling_test.go` pattern  
**Rationale:** The project already has working e2e tests that use polling. This approach:
- Uses public DockerHub images (no credentials needed)
- Has predictable version tags (`keelhq/push-workflow-example`)
- Is already proven to work in the codebase

### Decision 2: Run Keel In-Cluster (Not --no-incluster)
**Choice:** Deploy Keel as a pod inside the kind cluster  
**Rationale:** 
- More realistic production scenario
- Tests RBAC permissions properly
- Simpler than running keel binary outside the cluster

### Decision 3: Build and Load Local Image
**Choice:** Build keel Docker image and load it into kind cluster  
**Rationale:**
- Tests the actual built artifact
- Avoids needing to push to a registry
- `kind load docker-image` makes this straightforward

### Decision 4: Use Simple Kubernetes Manifests (Not Helm)
**Choice:** Apply keel deployment via kubectl with simple YAML manifests  
**Rationale:**
- Simpler setup, fewer dependencies
- Helm chart is complex with many options
- E2E test only needs core functionality

## Component Design

### 1. GitHub Actions Workflow (`.github/workflows/e2e.yml`)

```yaml
# Triggers: push, PR, manual
# Steps:
# 1. Checkout
# 2. Setup Go
# 3. Setup Docker
# 4. Run start-local-cluster.sh
# 5. Build keel Docker image
# 6. Load image into kind
# 7. Apply keel manifests (RBAC + Deployment)
# 8. Wait for keel to be ready
# 9. Create test deployment
# 10. Wait for update and verify
# 11. Cleanup (kind delete cluster)
```

### 2. Test Manifests (`tests/e2e/manifests/`)

**keel.yaml** - Minimal keel deployment:
- ServiceAccount with RBAC
- Deployment with polling enabled
- Environment: `POLL=1`, `POLL_DEFAULTSCHEDULE=@every 5s`

**test-deployment.yaml** - Test workload:
- Labels: `keel.sh/policy: major`, `keel.sh/trigger: poll`
- Annotation: `keel.sh/pollSchedule: @every 5s`
- Image: `keelhq/push-workflow-example:0.1.0`

### 3. Verification Script (`tests/e2e/verify.sh`)

Simple bash script that:
1. Waits up to 120 seconds for deployment image to change
2. Checks if image was updated to expected version (`0.10.0`)
3. Exits with appropriate code

## Codebase Patterns Found

- **Existing e2e tests**: Located in `tests/` directory, use Go test framework
- **Local cluster script**: `scripts/start-local-cluster.sh` creates `keel-dev` cluster
- **Test image**: `keelhq/push-workflow-example` has versions 0.1.0 through 0.10.0
- **Keel labels/annotations**: Defined in `types/types.go` (e.g., `types.KeelPolicyLabel`)
- **CI workflow**: `.github/workflows/ci.yml` shows existing GitHub Actions patterns

## Files to Create/Modify

| File | Action | Purpose |
|------|--------|---------|
| `.github/workflows/e2e.yml` | Create | New e2e workflow |
| `tests/e2e/manifests/keel.yaml` | Create | Keel deployment manifest |
| `tests/e2e/manifests/test-deployment.yaml` | Create | Test workload manifest |
| `tests/e2e/verify.sh` | Create | Verification script |

## Constraints & Gotchas

1. **Kind cluster needs Docker**: GitHub Actions `ubuntu-latest` has Docker pre-installed
2. **Image loading**: Must use `kind load docker-image` before referencing local images
3. **Poll schedule**: Use short intervals (5s) for fast test feedback
4. **Timeout handling**: kind cluster creation can take 1-2 minutes
5. **RBAC**: Keel needs permissions to list/update deployments across namespaces