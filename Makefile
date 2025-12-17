main_package_path = ./cmd/server
binary_name = argus

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	@test -z "$(shell git status --porcelain)"

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## deps: install all necessary development tools
.PHONY: deps
deps:
	go install github.com/air-verse/air@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest

## audit: run quality control checks
.PHONY: audit
audit: test
	go mod tidy -diff
	go mod verify
	go fmt ./...
	go vet ./...
	staticcheck ./...
	govulncheck ./...

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## test/integration: run integration tests with coverage
.PHONY: test/integration
test/integration:
	go test -v -tags=integration -coverpkg=./internal/analyzer -coverprofile=/tmp/integration_coverage.out ./cmd/integration
	go tool cover -html=/tmp/integration_coverage.out

## test/ci-cover: run tests with coverage for CI badge
.PHONY: test/ci-cover
test/ci-cover:
	go test -v -buildvcs -covermode=atomic -coverprofile=coverage.raw.out ./...

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy
	go fmt ./...

## build: build the application
.PHONY: build
build:
	go build -o=/tmp/bin/${binary_name} ${main_package_path}

## run: run the application
.PHONY: run
run: build
	/tmp/bin/${binary_name}

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	air \
	--build.cmd "go build -o /tmp/bin/${binary_name} ${main_package_path}" --build.bin "/tmp/bin/${binary_name}" --build.delay "100" \
	--build.exclude_dir "" \
	--build.include_ext "go, tpl, tmpl, html" \
	--misc.clean_on_exit "true"

# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

## push: push changes to the remote Git repository
.PHONY: push
push: confirm audit no-dirty
	git push