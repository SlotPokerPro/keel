# Design: Local Kubernetes Setup Script for Keel Development

## Architecture Overview

A simple shell script that creates a local Kubernetes cluster using **kind** (Kubernetes in Docker), which Keel can connect to via standard kubeconfig.

```
┌────────────────────────────────────────────────────────────┐
│                    Developer Machine                        │
│                                                             │
│  ┌──────────────────┐         ┌──────────────────────────┐ │
│  │   Keel Process   │         │   kind Container         │ │
│  │                  │  K8s    │   ┌──────────────────┐   │ │
│  │ keel --no-incluster ──────▶│   │  Control Plane   │   │ │
│  │                  │  API    │   │  (single node)   │   │ │
│  │  Port 9300 (UI)  │         │   └──────────────────┘   │ │
│  └──────────────────┘         │          │               │ │
│           │                   │   Port 6443 (mapped)     │ │
│           ▼                   └──────────────────────────┘ │
│    ~/.kube/config                                          │
│    (kind-keel-dev context)                                 │
└────────────────────────────────────────────────────────────┘
```

## Key Decisions

### Decision 1: Use kind over k3s/k3d
**Rationale**: 
- kind is simpler and more reliable in Docker environments
- No cgroup issues that affect k3s in containerized dev environments
- Widely used in CI/CD and local development
- Single binary, easy to install

### Decision 2: Single shell script, no wrapper
**Rationale**: Task is simple - just start a cluster. A shell script is sufficient. No need for Python/Go wrappers.

### Decision 3: Named cluster "keel-dev"
**Rationale**: Clear naming allows easy identification and avoids conflicts with other kind clusters the developer might have.

## Implementation Details

### Script: `scripts/start-local-cluster.sh`

Location in keel repo: `keel/scripts/start-local-cluster.sh`

**Components:**
1. Check for Docker
2. Install kind if not present
3. Create cluster if not exists
4. Export kubeconfig
5. Print usage instructions

### Cluster Configuration

| Setting | Value |
|---------|-------|
| Name | `keel-dev` |
| Nodes | 1 (control-plane only) |
| K8s version | kind default (latest stable) |
| API port | Auto-assigned by kind |

### Kubeconfig

kind automatically merges into `~/.kube/config` with context `kind-keel-dev`.

## Discovered Patterns

From `keel/cmd/keel/main.go`:
- `--no-incluster` flag uses `~/.kube/config` by default
- Can override with `KUBERNETES_CONFIG` env var
- Can specify context with `KUBERNETES_CONTEXT` env var

From `keel/.test/e2e-kind.sh`:
- Keel already has a kind-based e2e test script
- Uses `kind create cluster` with specific K8s version

## Usage

```bash
# Start cluster
./scripts/start-local-cluster.sh

# Run Keel (in another terminal)
cd cmd/keel && go build && ./keel --no-incluster

# Or with make
make run  # Already uses --no-incluster

# Clean up
kind delete cluster --name keel-dev
```

## Environment Variables for Keel

When running Keel locally:
```bash
export KUBECONFIG=~/.kube/config
export KUBERNETES_CONTEXT=kind-keel-dev  # Optional, if multiple contexts
keel --no-incluster
```

## Implementation Notes

- **Script tested successfully**: Cluster creation, idempotency, and kubectl access all verified working
- **kind version**: Script works with kind v0.20.0+ (tested with v0.27.0)
- **Kubeconfig**: kind automatically merges config into `~/.kube/config` with context `kind-keel-dev`
- **CGO requirement**: Keel requires `CGO_ENABLED=1` for sqlite support (needs gcc). This is a build requirement, unrelated to this script.
- **Deployments work**: Verified `kubectl create deployment` works against the cluster
- **Idempotency confirmed**: Running script multiple times correctly detects existing cluster
