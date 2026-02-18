# Requirements: Swaggo API Documentation for Keel HTTP Handlers

## Overview

Add Swagger/OpenAPI documentation to all HTTP API handlers in the Keel project using [swaggo/swag](https://github.com/swaggo/swag), and provide a Makefile command to generate the OpenAPI spec file.

## User Stories

### US1: API Documentation Generation
As a developer, I want to run `make update_openapi` to generate an up-to-date OpenAPI specification file so that I can share accurate API documentation with consumers.

### US2: Handler Documentation
As a developer, I want all HTTP handlers to have Swaggo annotations so that the generated OpenAPI spec accurately reflects the API's endpoints, request/response schemas, and parameters.

## Acceptance Criteria

### AC1: Makefile Command
- [ ] `make update_openapi` command exists in the Makefile
- [ ] Running the command generates/updates the OpenAPI spec file (e.g., `docs/swagger.json` or `docs/swagger.yaml`)
- [ ] Command installs swag CLI if not present

### AC2: Handler Annotations
All handlers in `pkg/http/` must have Swaggo annotations:
- [ ] `healthHandler` - GET /healthz
- [ ] `versionHandler` - GET /version
- [ ] Auth endpoints: `loginHandler`, `logoutHandler`, `refreshHandler`, `userInfoHandler`
- [ ] Approvals endpoints: `approvalsHandler`, `approvalSetHandler`, `approvalApproveHandler`
- [ ] Resources endpoint: `resourcesHandler`
- [ ] Policies endpoint: `policyUpdateHandler`
- [ ] Tracked images endpoints: `trackedHandler`, `trackSetHandler`
- [ ] Audit endpoint: `adminAuditLogHandler`
- [ ] Stats endpoint: `statsHandler`
- [ ] Webhook triggers: `nativeHandler`, `dockerHubHandler`, `jfrogHandler`, `quayHandler`, `azureHandler`, `githubHandler`, `harborHandler`, `registryNotificationHandler`

### AC3: Documentation Quality
- [ ] Each endpoint documents HTTP method, path, and description
- [ ] Request body schemas are documented where applicable
- [ ] Response schemas are documented with status codes
- [ ] Query parameters are documented where applicable
- [ ] Authentication requirements are noted

## Out of Scope

- Serving Swagger UI from the application
- Auto-generating client SDKs
- API versioning changes