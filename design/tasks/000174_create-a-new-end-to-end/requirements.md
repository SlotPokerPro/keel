# Requirements: End-to-End Test for Keel in GitHub Actions

## Overview

Create an end-to-end test that validates Keel's core functionality (automatic deployment updates) using a local Kubernetes cluster in GitHub Actions CI.

## User Stories

### US-1: CI Verification of Keel Functionality
**As a** developer  
**I want** automated e2e tests running in GitHub Actions  
**So that** I can verify Keel correctly updates deployments before merging changes

### US-2: Local Cluster Testing
**As a** contributor  
**I want** tests to use `scripts/start-local-cluster.sh`  
**So that** the test setup is consistent with local development

## Acceptance Criteria

### AC-1: GitHub Actions Workflow
- [ ] New workflow file exists at `.github/workflows/e2e.yml`
- [ ] Workflow triggers on push to master, pull requests, and manual dispatch
- [ ] Workflow uses `ubuntu-latest` runner

### AC-2: Cluster Setup
- [ ] Workflow calls `scripts/start-local-cluster.sh` to create kind cluster
- [ ] Cluster is named `keel-dev` (as defined in the script)
- [ ] kubectl is configured to use the kind cluster context

### AC-3: Keel Deployment
- [ ] Keel is built from source and deployed to the cluster
- [ ] Keel runs with polling enabled for the e2e test
- [ ] Keel has necessary RBAC permissions to update deployments

### AC-4: Test Deployment
- [ ] A test deployment is created with Keel annotations (`keel.sh/policy`, `keel.sh/trigger`)
- [ ] Deployment uses a public image with multiple semver tags (e.g., `keelhq/push-workflow-example`)
- [ ] Initial image version is set to an older version

### AC-5: Verification
- [ ] Test waits for Keel to update the deployment to a newer version
- [ ] Test has appropriate timeout (e.g., 2-3 minutes)
- [ ] Test reports clear pass/fail status

### AC-6: Cleanup
- [ ] Kind cluster is deleted after test completion
- [ ] Cleanup runs even if test fails

## Non-Functional Requirements

- Test should complete within 10 minutes
- Test should not require external credentials or secrets
- Test should use publicly available container images