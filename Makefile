# Kurso monorepo — developer commands.
# Run `make help` for the list.

# Host ports for local infra. Override POSTGRES_PORT/REDIS_PORT (env or `make`
# arg) if 5432/6379 are already taken locally, e.g. `POSTGRES_PORT=5433 make up`.
POSTGRES_PORT ?= 5432
REDIS_PORT    ?= 6379
DATABASE_URL  ?= postgres://kurso:kurso@localhost:$(POSTGRES_PORT)/kurso?sslmode=disable
REDIS_URL     ?= redis://localhost:$(REDIS_PORT)/0
FRONTENDS     := kurso-web kurso-admin kurso-partner
COMPOSE       := docker compose

# Export so `docker compose` (port mapping) and `go run` (app config) agree.
export POSTGRES_PORT
export REDIS_PORT
export DATABASE_URL
export REDIS_URL

.DEFAULT_GOAL := help

# ─── Help ────────────────────────────────────────────────────────────────────
.PHONY: help
help: ## Show this help
	@grep -hE '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

# ─── Full stack in Docker (production-parity; NOT needed for local dev) ───────
.PHONY: stack-up
stack-up: ## Build & start the whole stack in Docker (pg+redis+api+3 fronts)
	$(COMPOSE) up -d --build

.PHONY: stack-down
stack-down: ## Stop the Docker stack
	$(COMPOSE) down

.PHONY: stack-down-v
stack-down-v: ## Stop the Docker stack and delete its volumes
	$(COMPOSE) down -v

.PHONY: stack-logs
stack-logs: ## Tail Docker stack logs
	$(COMPOSE) logs -f

# ─── Optional: containerised Postgres+Redis (only if you DON'T run them locally)
.PHONY: db-up
db-up: ## Optional: start ONLY Postgres+Redis in Docker
	$(COMPOSE) up -d postgres redis

.PHONY: db-down
db-down: ## Optional: stop the containerised Postgres+Redis
	$(COMPOSE) stop postgres redis

# ─── Backend ─────────────────────────────────────────────────────────────────
.PHONY: dev
dev: ## Run API + all 3 frontends together (Ctrl+C stops everything)
	pnpm exec concurrently --kill-others --names api,web,admin,partner --prefix-colors blue,green,magenta,cyan \
		"cd kurso-api && go run ./cmd/kurso" \
		"cd kurso-web && pnpm dev" \
		"cd kurso-admin && pnpm dev" \
		"cd kurso-partner && pnpm dev"

.PHONY: run
run: ## Run only the API locally (go run) against your HOST Postgres/Redis
	cd kurso-api && go run ./cmd/kurso

.PHONY: db-create
db-create: ## One-time: create the kurso role + database in your HOST Postgres
	psql postgres -c "CREATE ROLE kurso WITH LOGIN PASSWORD 'kurso' CREATEDB;" 2>/dev/null || true
	createdb -O kurso kurso 2>/dev/null || true

.PHONY: build
build: build-api build-frontends ## Build everything

.PHONY: build-api
build-api: ## Build the API binary
	cd kurso-api && go build -o bin/kurso ./cmd/kurso

.PHONY: build-frontends
build-frontends: ## Build all three frontends
	@for d in $(FRONTENDS); do echo "==> build $$d"; (cd $$d && pnpm build) || exit 1; done

# ─── Migrations (goose) ──────────────────────────────────────────────────────
.PHONY: migrate
migrate: ## Apply all migrations
	cd kurso-api && go tool goose -dir migrations postgres "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down: ## Roll back the last migration
	cd kurso-api && go tool goose -dir migrations postgres "$(DATABASE_URL)" down

.PHONY: migrate-status
migrate-status: ## Show migration status
	cd kurso-api && go tool goose -dir migrations postgres "$(DATABASE_URL)" status

.PHONY: migrate-reset
migrate-reset: ## Roll back ALL migrations
	cd kurso-api && go tool goose -dir migrations postgres "$(DATABASE_URL)" reset

.PHONY: migrate-create
migrate-create: ## Create a new migration: make migrate-create name=add_something
	@test -n "$(name)" || (echo "usage: make migrate-create name=<name>"; exit 1)
	cd kurso-api && go tool goose -dir migrations create $(name) sql

# ─── Codegen ─────────────────────────────────────────────────────────────────
.PHONY: generate
generate: generate-api generate-frontends ## Regenerate all code from api/openapi.yaml

.PHONY: generate-api
generate-api: ## Regenerate Go types from the OpenAPI spec
	cd kurso-api && go tool oapi-codegen -config internal/adapter/http/openapi/config.yaml ../api/openapi.yaml

.PHONY: generate-frontends
generate-frontends: ## Regenerate TypeScript types in every frontend
	@for d in $(FRONTENDS); do echo "==> generate $$d"; (cd $$d && pnpm run openapi) || exit 1; done

# ─── Quality ─────────────────────────────────────────────────────────────────
.PHONY: test
test: ## Run backend tests
	cd kurso-api && go test ./...

.PHONY: lint
lint: lint-api lint-frontends ## Lint everything

.PHONY: lint-api
lint-api: ## Vet + gofmt check the backend
	cd kurso-api && go vet ./...
	@cd kurso-api && test -z "$$(gofmt -l .)" || (echo "gofmt: files need formatting:"; gofmt -l .; exit 1)

.PHONY: lint-frontends
lint-frontends: ## Lint all three frontends
	@for d in $(FRONTENDS); do echo "==> lint $$d"; (cd $$d && pnpm lint) || exit 1; done

.PHONY: fmt
fmt: ## Format the backend (gofmt)
	cd kurso-api && gofmt -w .

# ─── Setup ───────────────────────────────────────────────────────────────────
.PHONY: install
install: ## Install all dependencies (root hooks, Go modules, frontends)
	pnpm install
	cd kurso-api && go mod download
	@for d in $(FRONTENDS); do echo "==> install $$d"; (cd $$d && pnpm install) || exit 1; done

.PHONY: tidy
tidy: ## Tidy Go modules
	cd kurso-api && go mod tidy
