include .env


up: build
	docker-compose up -d

build:
	docker-compose build

down:
	docker-compose down

go:
	docker-compose exec -it w3st-cms /bin/sh

db:
	docker exec -it ${DB_HOST} psql -U ${DB_USER} -d ${DB_NAME}


#moke作成

# 変数定義
MOCKGEN=mockgen
SRC_DIR=domain
SRC_INTERFACE=interfaces
MOCK_DIR=mock
REPO_PKG=repositories
SERVICE_PKG=services

# 汎用的なモック生成ルール
# Usage: make mockgen SRC=xxx.go DST=yyy.go PKG=zzz
mockgen:
	$(MOCKGEN) -source=src/$(SRC) -destination=src/$(MOCK_DIR)/$(DST) -package=$(PKG)

# 便利なターゲット（よく使うものをここで一括指定）
mock-user:
	$(MOCKGEN) -source=src/$(SRC_DIR)/$(REPO_PKG)/users.go -destination=$src/(MOCK_DIR)/$(REPO_PKG)/mock_user_repository.go -package=mock_repositories

mock-auth:
	$(MOCKGEN) -source=src/$(SRC_INTERFACE)/$(SERVICE_PKG)/auth.go -destination=src/$(MOCK_DIR)/$(SERVICE_PKG)/mock_auth_service.go -package=mock_services

mock-tx:
	$(MOCKGEN) -source=src/$(SRC_DIR)/$(REPO_PKG)/transaction.go -destination=src/$(MOCK_DIR)/$(REPO_PKG)/mock_transaction_repository.go -package=mock_repositories

# まとめて実行
mock-all: mock-user mock-auth mock-tx


.PHONY: lint
lint:
	cd src && golangci-lint run ./... \
		--skip-dirs=mocks \
		--skip-dirs=tmp \
		--skip-files=cover.out \
		--skip-files=cover.html \
		--skip-files=.air.toml