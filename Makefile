-include .env

# コンテナ操作
up: build
	docker-compose up -d

build:
	docker-compose build --no-cache

down:
	docker-compose down

go:
	docker-compose exec -it w3st-cms /bin/sh

db:
	docker exec -it ${DB_HOST} psql -U ${DB_USER} -d ${DB_NAME}

# ---------- Mock 生成 ----------
MOCKGEN = mockgen
SRC_DIR = domain
SRC_INTERFACE = interfaces
MOCK_DIR = mock
REPO_PKG = repositories
SERVICE_PKG = services

mockgen:
	$(MOCKGEN) -source=src/$(SRC) -destination=src/$(MOCK_DIR)/$(DST) -package=$(PKG)

mock-user:
	$(MOCKGEN) -source=src/$(SRC_DIR)/$(REPO_PKG)/users.go -destination=src/$(MOCK_DIR)/$(REPO_PKG)/mock_user_repository.go -package=mock_repositories

mock-auth:
	$(MOCKGEN) -source=src/$(SRC_INTERFACE)/$(SERVICE_PKG)/auth.go -destination=src/$(MOCK_DIR)/$(SERVICE_PKG)/mock_auth_service.go -package=mock_services

mock-tx:
	$(MOCKGEN) -source=src/$(SRC_DIR)/$(REPO_PKG)/transaction.go -destination=src/$(MOCK_DIR)/$(REPO_PKG)/mock_transaction_repository.go -package=mock_repositories

mock-all: mock-user mock-auth mock-tx

# ---------- Format / Lint ----------
GOFMT = gofmt
GOFUMPT = gofumpt
GOIMPORTS = goimports
GOLANGCI_LINT = golangci-lint
GOFILES := $(shell find src -name '*.go' -not -path "./vendor/*")

.PHONY: fmt lint fmt-lint

fmt:
	@echo "Running gofumpt..."
	$(GOFUMPT) -w $(GOFILES)
	@echo "Running goimports..."
	$(GOIMPORTS) -w $(GOFILES)
	@echo "Running gofmt..."
	$(GOFMT) -s -w $(GOFILES)

lint:
	cd src && $(GOLANGCI_LINT) run ./... --timeout 5m

fmt-lint: fmt lint
