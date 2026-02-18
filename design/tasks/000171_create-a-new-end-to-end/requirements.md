# Requirements: End-to-End Test with Kind Cluster in GitHub Actions

## Overview

Add an end-to-end test workflow to GitHub Actions that uses the existing `scripts/start-local-cluster.sh` script to spin up a kind (Kubernetes in Docker) cluster and verify Keel can automatically update deployments.

## User Stories

### US1: Automated E2E Testing in CI
As a developer, I want e2e tests to run automatically on PRs and pushes to master, so that I can catch integration issues before merging.

### US2: Validate Auto-Update Functionality
As a developer, I want to verify that Keel correctly updates a deployment when a new image version is detected, so that I have confidence in the core functionality.

## Acceptance Criteria

1. **Workflow triggers** on push to master and pull requests
2. **Kind cluster** is created using `scripts/start-local-cluster.sh`
3. **Keel binary** is built and runs against the cluster
4. **Test deployment** is created with Keel policy labels
5. **Auto-update** is verified (deployment image changes from old to new version)
6. **Workflow completes** within a reasonable timeout (~10 minutes)
7. **Cleanup** happens automatically (kind cluster deleted)

## Out of Scope

- Helm chart deployment tests (existing tests cover this separately)
- Multi-architecture testing
- Performance benchmarking