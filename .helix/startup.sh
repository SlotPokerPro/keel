#!/bin/bash
set -euo pipefail

# Project startup script for Keel development
# This runs when agents start working on this project
# Idempotent - safe to run multiple times
# Uses k3s (lightweight Kubernetes) for local development
#
# NOTE: Running Kubernetes inside Docker containers requires proper cgroup v2 support
# with memory controller enabled. If this fails, see workarounds at the end of the script.

