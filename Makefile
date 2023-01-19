.PHONY:	go-imports go-lint


go-imports:
	@echo Checking correct formatting of files
	@{ \
  		GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports; \
		files=$$( goimports -w -l . ); \
		if [ -n "$$files" ]; then \
		echo "Files not properly formatted: $$files"; \
		exit 1; \
		fi; \
		if ! go vet ./...; then \
		exit 1; \
		fi \
	}

go-lint:
	@echo Checking linting of files
	@{ \
		GO111MODULE=off go get -u golang.org/x/lint/golint; \
		el="_test.go"; \
		lintfiles=$$( golint ./... | egrep -v "$$el" ); \
		if [ -n "$$lintfiles" ]; then \
		echo "Lint errors:"; \
		echo "$$lintfiles"; \
		exit 1; \
		fi \
	}

## Use docker mounting the current directory to /
GO_SWAGGER ?= docker run --platform linux/amd64 --rm -u "$(USER_GROUP)" -v "$(CURDIR)":/ -w / quay.io/goswagger/swagger:v0.25.0

go-swagger-gen: ## Generate go model code from swagger spec
	$(GO_SWAGGER) -q generate model \
		--accept-definitions-only \
		--target=./pkg/ \
		--spec=./swagger.yml