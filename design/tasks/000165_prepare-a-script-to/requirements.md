# Requirements: Local Kubernetes Setup Script for Keel Development

## Overview

Create a simple shell script to start a local Kubernetes cluster (kind or k3s) that Keel can connect to when running locally with `keel --no-incluster`.

## User Stories

### US1: As a developer, I want to start a local Kubernetes cluster
- **Given** I have Docker installed
- **When** I run the setup script
- **Then** a local Kubernetes cluster should be created and running
- **And** the kubeconfig should be configured for kubectl access

### US2: As a developer, I want to run Keel against the local cluster
- **Given** the local cluster is running
- **When** I run `keel --no-incluster`
- **Then** Keel should connect to the local cluster using my kubeconfig
- **And** Keel should be able to watch and update deployments

### US3: As a developer, I want the script to be idempotent
- **Given** the cluster is already running
- **When** I run the script again
- **Then** it should not fail or create duplicate clusters

### US4: As a developer, I want to easily tear down the cluster
- **Given** the cluster is running
- **When** I want to clean up
- **Then** there should be clear instructions or a command to delete the cluster

## Acceptance Criteria

1. **Cluster creation**: Script creates a working single-node Kubernetes cluster
2. **Kubeconfig setup**: Kubeconfig is available at `~/.kube/config` or exported via `KUBECONFIG`
3. **Keel compatibility**: Keel can connect using `--no-incluster` flag
4. **kubectl access**: `kubectl get nodes` works after script completes
5. **Idempotent**: Running the script multiple times doesn't cause errors
6. **Lightweight**: Uses kind (preferred) or k3d - both work well in Docker environments
7. **Simple**: Script is < 100 lines and easy to understand

## Out of Scope

- Installing Go or building Keel (user handles this separately)
- Starting Keel automatically (user runs `keel --no-incluster` manually)
- Production-grade cluster setup
- Multi-node cluster configuration