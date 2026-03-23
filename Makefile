SHELL := /bin/bash

# ⚙️ Configuration
APP             ?= alerter
DOCKER_COMPOSE  := COMPOSE_BAKE=true docker compose

DOCKERFILE ?= Dockerfile

IMAGE_REG  ?= ghcr.io/bruli
IMAGE_NAME := $(IMAGE_REG)/$(APP)
VERSION    ?= 0.4.0
CURRENT_IMAGE := $(IMAGE_NAME):$(VERSION)

# Default goal
.DEFAULT_GOAL := help

# 📚 Declare all phony targets
.PHONY: docker-logs docker-down docker-exec docker-ps docker-up \
        test test-functional clean fmt help \
        encryptVault decryptVault build deploy security docker-login docker-push-image


GOLANGCI_LINT_VERSION ?= v2.11.4
# ────────────────────────────────────────────────────────────────
# 🐳 Docker
# ────────────────────────────────────────────────────────────────
docker-up:
	@set -euo pipefail; \
	echo "🚀 Starting services with Docker Compose..."; \
	$(DOCKER_COMPOSE) up -d --build

docker-down:
	@set -euo pipefail; \
	echo "🛑 Stopping and removing Docker Compose services..."; \
	$(DOCKER_COMPOSE) down

docker-ps:
	@set -euo pipefail; \
	echo "📋 Active services:"; \
	$(DOCKER_COMPOSE) ps

docker-exec:
	@set -euo pipefail; \
	test -n "${SVC:-}" || { echo "❌ Please specify SVC=<service>"; exit 2; }; \
	echo "🔎 Opening shell inside $$SVC..."; \
	$(DOCKER_COMPOSE) exec $$SVC sh

docker-logs:
	@set -euo pipefail; \
	echo "👀 Showing logs for container $(APP) (CTRL+C to exit)..."; \
	docker logs -f $(APP)

# ────────────────────────────────────────────────────────────────
# 🧹 Code quality: format, lint, tests
# ────────────────────────────────────────────────────────────────
fmt:
	@set -euo pipefail; \
	echo "👉 Formatting code with gofumpt..."; \
	go tool gofumpt -w .

security:
	@set -euo pipefail; \
	echo "👉 Check security"; \
	go tool govulncheck ./...

.PHONY: install-lint
install-lint:
	@set -euo pipefail; \
    echo "🔧 Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."; \
    	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: lint
lint: install-lint
	@set -euo pipefail; \
	echo "🚀 Executing golangci-lint..."; \
    golangci-lint run ./...

test:
	@set -euo pipefail; \
	echo "🧪 Running unit tests (race, JSON → tparse)..."; \
	go test -race ./... -json -cover -coverprofile=coverage.out| go tool tparse -all

clean:
	@set -euo pipefail; \
	echo "🧹 Cleaning local artifacts..."; \
	rm -rf bin dist coverage .*cache || true; \
	go clean -testcache

check: fmt lint security test
	echo "✅ Format, linter and tests success."

docker-login:
	echo "🔐 Logging into Docker registry...";
	echo "$$CR_PAT" | docker login ghcr.io -u bruli --password-stdin


docker-push-image: docker-login
	echo "🐳 Building and pushing Docker image $(CURRENT_IMAGE) for (prod)...";
	docker buildx build \
		-t $(CURRENT_IMAGE) \
		-f $(DOCKERFILE) \
		--push \
		.
	 echo "✅ Image $(CURRENT_IMAGE) pushed successfully."

# ────────────────────────────────────────────────────────────────
# ℹ️ Help
# ────────────────────────────────────────────────────────────────
help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:' Makefile | awk -F':' '{print "  - " $$1}'
