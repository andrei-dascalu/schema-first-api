.PHONY: install-go-tools install-golangci-lint install-oapi-codegen \
	install-sqlc install-goose tool-versions \
	generate generate-backend generate-frontend \
	dev-db dev-db-down migrate-up migrate-down migration-create \
	build run-backend run-frontend \
	lint lint-backend lint-frontend \
	test test-unit test-integration test-frontend

REPO_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
LOCAL_BIN := $(REPO_ROOT)/.bin

# Container runtime: prefer Podman, fall back to Docker if it's not on PATH.
CONTAINER_ENGINE := $(shell command -v podman >/dev/null 2>&1 && echo podman || echo docker)

# Pinned repository-local Go CLI tools. Versioned stamps force a reinstall
# after a pin changes, even when a binary with the same name already exists.
GOLANGCI_LINT_VERSION := v2.13.2
OAPI_CODEGEN_VERSION := v2.8.0
SQLC_VERSION := v1.31.1
GOOSE_VERSION := v3.28.0

GOLANGCI_LINT := $(LOCAL_BIN)/golangci-lint
OAPI_CODEGEN := $(LOCAL_BIN)/oapi-codegen
SQLC := $(LOCAL_BIN)/sqlc
GOOSE := $(LOCAL_BIN)/goose

GOLANGCI_LINT_STAMP := $(LOCAL_BIN)/.golangci-lint-$(GOLANGCI_LINT_VERSION)
OAPI_CODEGEN_STAMP := $(LOCAL_BIN)/.oapi-codegen-$(OAPI_CODEGEN_VERSION)
SQLC_STAMP := $(LOCAL_BIN)/.sqlc-$(SQLC_VERSION)
GOOSE_STAMP := $(LOCAL_BIN)/.goose-$(GOOSE_VERSION)

# Load PORT/DATABASE_URL/JWT_SECRET/CORS_ORIGINS from backend/.env (see
# backend/.env.example) so migrate-up/down and run-backend work without
# manually exporting them first. Note: values here take precedence over a
# same-named shell export, since Makefile assignments (including `include`d
# ones) override the environment.
-include backend/.env
export PORT DATABASE_URL JWT_SECRET CORS_ORIGINS

# --- local Go tools ----------------------------------------------------
install-golangci-lint:
	@mkdir -p "$(LOCAL_BIN)"
	@if [ ! -x "$(GOLANGCI_LINT)" ] || [ ! -f "$(GOLANGCI_LINT_STAMP)" ]; then \
		echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION) in $(LOCAL_BIN)..."; \
		GOBIN="$(LOCAL_BIN)" go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) && \
			touch "$(GOLANGCI_LINT_STAMP)"; \
	fi

install-oapi-codegen:
	@mkdir -p "$(LOCAL_BIN)"
	@if [ ! -x "$(OAPI_CODEGEN)" ] || [ ! -f "$(OAPI_CODEGEN_STAMP)" ]; then \
		echo "Installing oapi-codegen $(OAPI_CODEGEN_VERSION) in $(LOCAL_BIN)..."; \
		GOBIN="$(LOCAL_BIN)" go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION) && \
			touch "$(OAPI_CODEGEN_STAMP)"; \
	fi

install-sqlc:
	@mkdir -p "$(LOCAL_BIN)"
	@if [ ! -x "$(SQLC)" ] || [ ! -f "$(SQLC_STAMP)" ]; then \
		echo "Installing sqlc $(SQLC_VERSION) in $(LOCAL_BIN)..."; \
		GOBIN="$(LOCAL_BIN)" go install github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION) && \
			touch "$(SQLC_STAMP)"; \
	fi

install-goose:
	@mkdir -p "$(LOCAL_BIN)"
	@if [ ! -x "$(GOOSE)" ] || [ ! -f "$(GOOSE_STAMP)" ]; then \
		echo "Installing goose $(GOOSE_VERSION) in $(LOCAL_BIN)..."; \
		GOBIN="$(LOCAL_BIN)" go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION) && \
			touch "$(GOOSE_STAMP)"; \
	fi

install-go-tools: install-golangci-lint install-oapi-codegen install-sqlc install-goose

tool-versions: install-go-tools
	@"$(GOLANGCI_LINT)" version
	@"$(OAPI_CODEGEN)" -version
	@"$(SQLC)" version
	@"$(GOOSE)" -version

# --- codegen -----------------------------------------------------------
# Regenerate everything derived from docs/api_schema.yaml: the Go server
# (oapi-codegen), sqlc's typed query code, and the frontend client
# (@hey-api/openapi-ts). Run after editing the schema or any *.sql query.
generate: generate-backend generate-frontend

generate-backend: install-oapi-codegen install-sqlc
	cd backend && go generate ./...
	cd backend && "$(SQLC)" generate

generate-frontend:
	cd frontend && npm run generate

# --- local infra ---------------------------------------------------------
dev-db:
	$(CONTAINER_ENGINE) compose -f compose.yml up -d --wait

dev-db-down:
	$(CONTAINER_ENGINE) compose -f compose.yml down

migrate-up:
	cd backend && go run ./cmd/migrate up

migrate-down:
	cd backend && go run ./cmd/migrate down

migration-create: install-goose
	@test -n "$(NAME)" || (echo "NAME is required, e.g. make migration-create NAME=add_event_owner" && exit 1)
	"$(GOOSE)" -dir backend/internal/adapters/postgres/migrations -s create "$(NAME)" sql

# --- build & run -----------------------------------------------------------
build:
	cd backend && go build ./...

run-backend:
	cd backend && go run ./cmd/api

run-frontend:
	cd frontend && npm run dev

# --- lint --------------------------------------------------------------
# Generated code (internal/api/generated.go, internal/store/postgres/db/*,
# frontend/src/api/generated/**) is excluded via .golangci.yml / eslint's
# ignores config rather than here.
lint: lint-backend lint-frontend

lint-backend: install-golangci-lint
	cd backend && "$(GOLANGCI_LINT)" run ./...

lint-frontend:
	cd frontend && npm run lint

# --- tests ---------------------------------------------------------------
# Unit tests carry no build tag; integration tests are tagged `integration`
# and spin up Postgres via testcontainers-go (needs a running container
# runtime - see the Podman/Docker + testcontainers note in CLAUDE.md).
test: test-unit test-integration test-frontend

test-unit:
	cd backend && go test ./...

# testcontainers-go talks to the container runtime's Docker-compatible API. A
# real Docker daemon exposes that at its usual socket, so no extra setup is
# needed there. Podman's rootless machine (macOS/Windows) needs DOCKER_HOST
# pointed explicitly at its API socket, and Ryuk (testcontainers' cleanup
# sidecar) disabled since it doesn't work against rootless Podman.
ifeq ($(CONTAINER_ENGINE),podman)
test-integration:
	DOCKER_HOST="unix://$$(podman machine inspect --format '{{.ConnectionInfo.PodmanSocket.Path}}')" \
	TESTCONTAINERS_RYUK_DISABLED=true \
	sh -c 'cd backend && go test -tags=integration ./...'
else
test-integration:
	cd backend && go test -tags=integration ./...
endif

test-frontend:
	cd frontend && npm run test
