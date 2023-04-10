help: ## print available targets
	@cat $(MAKEFILE_LIST) | \
	grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

test: ## run unit tests for the StatBud server
	go test ./... -coverprofile=coverage.out

coverage: test ## run unit tests and code coverage for the StatBud server
	go tool cover -html=coverage.out
