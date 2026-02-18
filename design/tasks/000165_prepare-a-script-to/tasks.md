# Implementation Tasks

## Script Setup
- [~] Create `keel/scripts/` directory if it doesn't exist
- [~] Create `keel/scripts/start-local-cluster.sh` script file
- [~] Make script executable (`chmod +x`)

## Script Implementation
- [ ] Add shebang and error handling (`set -euo pipefail`)
- [ ] Add function to check Docker is installed and running
- [ ] Add function to check/install kind binary
- [ ] Add function to check if cluster `keel-dev` already exists
- [ ] Add function to create cluster if it doesn't exist
- [ ] Add verification that `kubectl get nodes` works
- [ ] Print success message with instructions for running Keel

## Cleanup Helper
- [ ] Add usage message showing how to delete cluster (`kind delete cluster --name keel-dev`)

## Testing
- [ ] Test script on clean environment (no kind installed)
- [ ] Test script when cluster already exists (idempotency)
- [ ] Test that `keel --no-incluster` can connect to the cluster
- [ ] Verify deployments can be created and Keel can watch them

## Documentation
- [ ] Add brief comments in script explaining each step
- [ ] Update keel readme.md with local development instructions (optional)