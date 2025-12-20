GO_FILES = $(shell find . -type f -name '*.go')

.PHONY: build
build: $(_BINARY_NAME) ## Build the binary

$(_BINARY_NAME): $(GO_FILES)
	go build -o $(_BINARY_NAME) $(_MAIN_PACKAGE_PATH)

.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: clean
clean: ## Clean build artifacts
	rm -f $(_BINARY_NAME)

.PHONY: tidy
tidy: ## Tidy go modules
	go mod tidy
