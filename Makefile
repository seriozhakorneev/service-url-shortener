include .env.example
export

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

compose-up: ### Run docker-compose
	docker-compose up --build -d postgres
.PHONY: compose-up

compose-logs: ### Run docker-compose logs
	docker-compose logs -f
.PHONY: compose-up

compose-up-integration-test: ### Run docker-compose with integration test
	docker-compose up --build --abort-on-container-exit --exit-code-from integration
.PHONY: compose-up-integration-test

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

run: ### app run
	go mod tidy && go mod download && \
	GIN_MODE=debug GRPC_DEBUG=true CGO_ENABLED=0 go run -tags migrate ./cmd/app
.PHONY: run

docker-rm-volume: ### remove docker volume
	docker volume rm service-url-shortener_pg-data
.PHONY: docker-rm-volume

lint: ### check by all linters
	dotenv-linter \
	&& golangci-lint run \
	&& git ls-files -c --exclude='Dockerfile*' --ignored | xargs hadolint
.PHONY: lint

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

linter-hadolint: ### check by hadolint linter
	git ls-files --exclude='Dockerfile*' --ignored | xargs hadolint
.PHONY: linter-hadolint

linter-dotenv: ### check by dotenv linter
	dotenv-linter
.PHONY: linter-dotenv

test: ### run test
	go test -v -cover -race ./internal/...
.PHONY: test

integration-test: ### run integration-test
	go clean -testcache && go test -v ./integration-test/...
.PHONY: integration-test

migrate-create:  ### create new migration
	migrate create -ext sql -dir migrations 'init_schema'
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' -verbose up
.PHONY: migrate-up

migrate-down: ### migration down
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' -verbose down
.PHONY: migrate-down