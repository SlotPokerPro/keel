# Implementation Tasks

## Script Setup
- [x] Create `keel/scripts/` directory if it doesn't exist
- [x] Create `keel/scripts/start-local-cluster.sh` script file
- [x] Make script executable (`chmod +x`)

## Script Implementation
- [x] Add shebang and error handling (`set -euo pipefail`)
- [x] Add function to check Docker is installed and running
- [x] Add function to check/install kind binary
- [x] Add function to check if cluster `keel-dev` already exists
- [x] Add function to create cluster if it doesn't exist
- [x] Add verification that `kubectl get nodes` works
- [x] Print success message with instructions for running Keel

## Cleanup Helper
- [x] Add usage message showing how to delete cluster (`kind delete cluster --name keel-dev`)

## Testing
- [x] Test script on clean environment (no kind installed)
- [x] Test script when cluster already exists (idempotency)
- [x] Test that `keel --no-incluster` can connect to the cluster
- [x] Verify deployments can be created and Keel can watch them

## Documentation
- [x] Add brief comments in script explaining each step
- [x] Update keel readme.md with local development instructions (optional)

## Implementation Notes

- Script tested successfully: cluster creation, idempotency, kubectl access all working
- Deployments can be created in the cluster (`kubectl create deployment nginx-test --image=nginx:latest`)
- Kind automatically configures `~/.kube/config` with context `kind-keel-dev`
- Note: Keel requires `CGO_ENABLED=1` for sqlite support - unrelated to this script
- Keel connects via standard kubeconfig when using `--no-incluster` flag