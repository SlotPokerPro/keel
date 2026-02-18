# Implementation Tasks

## Setup

- [x] Add `update_openapi` target to `Makefile` that installs swag and runs `swag init`
- [x] Add general API info annotations to `cmd/keel/main.go` (@title, @version, @description, @host, @BasePath)

## Core Endpoints

- [x] Add Swaggo annotations to `healthHandler` in `pkg/http/http.go`
- [x] Add Swaggo annotations to `versionHandler` in `pkg/http/http.go`
- [x] Add Swaggo annotations to `userInfoHandler` in `pkg/http/http.go`

## Auth Endpoints (`pkg/http/auth.go`)

- [x] Add annotations to `loginHandler`
- [x] Add annotations to `logoutHandler`
- [x] Add annotations to `refreshHandler`

## Approvals Endpoints (`pkg/http/approvals_endpoint.go`)

- [x] Add annotations to `approvalsHandler` (GET /v1/approvals)
- [x] Add annotations to `approvalSetHandler` (PUT /v1/approvals)
- [x] Add annotations to `approvalApproveHandler` (POST /v1/approvals)

## Other Admin Endpoints

- [x] Add annotations to `resourcesHandler` in `pkg/http/resources_endpoint.go`
- [x] Add annotations to `policyUpdateHandler` in `pkg/http/policy_endpoint.go`
- [x] Add annotations to `trackedHandler` in `pkg/http/tracked_endpoint.go`
- [x] Add annotations to `trackSetHandler` in `pkg/http/tracked_endpoint.go`
- [x] Add annotations to `adminAuditLogHandler` in `pkg/http/audit_endpoint.go`
- [x] Add annotations to `statsHandler` in `pkg/http/stats_endpoint.go`

## Webhook Endpoints

- [~] Add annotations to `nativeHandler` in `pkg/http/native_webhook_trigger.go`
- [ ] Add annotations to `dockerHubHandler` in `pkg/http/dockerhub_webhook_trigger.go`
- [ ] Add annotations to `jfrogHandler` in `pkg/http/jfrog_webhook_trigger.go`
- [ ] Add annotations to `quayHandler` in `pkg/http/quay_webhook_trigger.go`
- [ ] Add annotations to `azureHandler` in `pkg/http/azure_webhook_trigger.go`
- [ ] Add annotations to `githubHandler` in `pkg/http/github_webhook_trigger.go`
- [ ] Add annotations to `harborHandler` in `pkg/http/harbor_webhook_trigger.go`
- [ ] Add annotations to `registryNotificationHandler` in `pkg/http/registry_notifications.go`

## Verification

- [ ] Run `make update_openapi` and verify docs are generated
- [ ] Verify generated `docs/swagger.json` contains all endpoints
- [ ] Add `docs/` to `.gitignore` or commit generated files (team decision)