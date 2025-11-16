TARGET ?= peerly
MODULE_NAME := $(shell grep '^module ' go.mod | cut -d' ' -f2)

BIN_DIR := ./.bin

SHELL := /bin/bash
GREEN  = \033[1;32m
YELLOW = \033[1;33m
PURPLE = \033[1;35m
RED    = \033[1;31m
RESET  = \033[0m

FORMAT_TOOLS = mvdan.cc/gofumpt@latest \
        github.com/daixiang0/gci@latest \
        github.com/segmentio/golines@latest \
        github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.4.0

API_TOOLS = \
    github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest \
    github.com/daveshanley/vacuum@latest \
    github.com/oasdiff/oasdiff@latest

STORAGE_TOOLS = \
	github.com/sqlc-dev/sqlc/cmd/sqlc@latest

COVERAGE_FILE ?= coverage.out
COVERAGE_EXCLUDE_PATTERNS := \
	$(MODULE_NAME)/api/docs% \
	$(MODULE_NAME)/mocks% \
	$(MODULE_NAME)/cmd% \
	$(MODULE_NAME)/internal/infrastructure/flags% \
	$(MODULE_NAME)/internal/infrastructure/configs% \
	$(MODULE_NAME)/pkg/logger% \
	$(MODULE_NAME)/api/middlewares% \
	$(MODULE_NAME)/api/routers% \
	$(MODULE_NAME)/api/mapper% \
	$(MODULE_NAME)/api/routers/v1% \
	$(MODULE_NAME)/api/handlers/v1% \
	$(MODULE_NAME)/internal/infrastructure/storages/postgres% \
	$(MODULE_NAME)/internal/infrastructure/repositories% \
	$(MODULE_NAME)/api/servers%
COVERAGE_PACKAGES := $(filter-out $(COVERAGE_EXCLUDE_PATTERNS),$(shell go list ./...))

DOCKER_DEV_COMPOSE := dev.docker-compose.yaml
DOCKER_PROD_COMPOSE := prod.docker-compose.yaml

SQLC := sqlc
MIGRATE := migrate

MIGRATIONS_DIR := migrations
SQLC_CONFIG := sqlc.yaml

OPENAPI_FILE := api/docs/openapi.yml
OPENAPI_GEN_DIR := api/generated
OPENAPI_GEN_CONFIG_DIR := oapi-codegen
OPENAPI_GEN_MODELS := models.yaml
OPENAPI_GEN_SERVER := server.yaml
OPENAPI_GEN_SPEC := spec.yaml

OPENAPI_LINT_RULES := .vacuum-rules.yaml

OAPI_CODEGEN := oapi-codegen
VACUUM := vacuum
OASDIFF := oasdiff
