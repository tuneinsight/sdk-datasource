APP_VERSION := $(scripts/version.sh)
BASE_PATH := /

# --- swagger
.PHONY:	go-swagger-preprocess go-swagger-validate go-swagger-clean go-swagger-gen
GO_SWAGGER ?= docker run --rm -u "$(USER_GROUP)" -v "$(PWD)":/src -w /src quay.io/goswagger/swagger:v0.25.0
YQ ?= docker run --rm -u "$(USER_GROUP)" -v "$(PWD)":/src -w /src mikefarah/yq:4.0.0
JQ ?= docker run --rm -u "$(USER_GROUP)" -v "$(PWD)":/src -w /src endeveit/docker-jq jq

go-swagger-preprocess:
	$(GO_SWAGGER) mixin --format=yaml --output=./geco-api/geco-flattened.yml \
		./geco-api/geco.yml ./geco-api/ml.yml ./geco-api/aggregation.yml ./geco-api/key-protocols.yml  #./geco-api/kaplan-meier.yml ./geco-api/set-intersection.yml
	$(GO_SWAGGER) flatten --with-flatten=minimal --format=yaml \
		--output=./geco-api/geco-flattened.yml ./geco-api/geco-flattened.yml
	$(YQ) -i eval '.info.version = "$(APP_VERSION)"' ./geco-api/geco-flattened.yml
	$(YQ) -i eval '.basePath = "$(BASE_PATH)"' ./geco-api/geco-flattened.yml

go-swagger-validate:
	$(GO_SWAGGER) validate ./geco-api/geco-flattened.yml

go-swagger-clean:
	rm -rf pkg/models

go-swagger-gen: go-swagger-preprocess go-swagger-validate
	$(GO_SWAGGER) generate model \
		--spec=./geco-api/geco-flattened.yml \
		--target=./ \
   		--model-package=pkg/models