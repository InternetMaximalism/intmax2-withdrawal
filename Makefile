ifneq (${GIT_USE},)
ifeq ($(shell git tag --contains HEAD),)
  VERSION := $(shell git rev-parse --short HEAD)
else
  VERSION := $(shell git tag --contains HEAD)
endif
endif

ifneq ($(goproxy),)
  re_build_arg := --build-arg goproxy="$(goproxy)"
endif

ifeq ($(shell uname -s),Darwin)
	SED_COMMAND := gsed
else
	SED_COMMAND := sed
endif

BUILDNAME := intmax2-withdrawal
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GOLDFLAGS += -X intmax2-withdrawal/configs/buildvars.Version=$(VERSION)
GOLDFLAGS += -X intmax2-withdrawal/configs/buildvars.BuildTime=$(BUILDTIME)
GOLDFLAGS += -X intmax2-withdrawal/configs/buildvars.BuildName=$(BUILDNAME)
GOFLAGS = -ldflags "$(GOLDFLAGS)"

.DEFAULT_GOAL := default

.PHONY: default
default: gen lint build

.PHONY: build
build:
	go build -v -o $(BUILDNAME) $(GOFLAGS) ./cmd/

.PHONY: gen
gen: format-proto
	buf generate -v --debug --timeout=2m --template api/proto/withdrawal_service/buf.gen.yaml api/proto/withdrawal_service
	buf generate -v --debug --timeout=2m --template api/proto/withdrawal_service/buf.gen.tagger.yaml api/proto/withdrawal_service
	go generate -v ./...
	cp -rf docs/swagger/withdrawal third_party/OpenAPI/withdrawal_service
	cp -rf third_party/OpenAPI/_default/* third_party/OpenAPI/withdrawal_service
# withdrawal
ifneq (${SWAGGER_USE},)
ifneq (${SWAGGER_BUILD_MODE},)
	$(SED_COMMAND) -i "s/SWAGGER_VERSION/$(VERSION)/g" third_party/OpenAPI/withdrawal_service/withdrawal/apidocs.swagger.json
else
	$(SED_COMMAND) -i "s/SWAGGER_VERSION/v0.0.0/g" third_party/OpenAPI/withdrawal_service/withdrawal/apidocs.swagger.json
endif
ifneq (${SWAGGER_HOST_URL},)
	$(SED_COMMAND) -i "s/SWAGGER_HOST_URL/${SWAGGER_HOST_URL}/g" third_party/OpenAPI/withdrawal_service/withdrawal/apidocs.swagger.json
else
	$(SED_COMMAND) -i "s/SWAGGER_HOST_URL//g" third_party/OpenAPI/withdrawal_service/withdrawal/apidocs.swagger.json
endif
ifneq (${SWAGGER_BASE_PATH},)
	$(SED_COMMAND) -i "s/SWAGGER_BASE_PATH/${SWAGGER_BASE_PATH}/g" third_party/OpenAPI/withdrawal_service/withdrawal/apidocs.swagger.json
else
	$(SED_COMMAND) -i "s/SWAGGER_BASE_PATH/\//g" third_party/OpenAPI/withdrawal_service/withdrawal/apidocs.swagger.json
endif
endif

.PHONY: format-proto
format-proto: ## format all protos
	clang-format -i api/proto/withdrawal_service/withdrawal/v1/withdrawal.proto

.PHONY: tools
tools:
	go install -v go.uber.org/mock/mockgen@v0.5.0
	go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
	go install -v github.com/bufbuild/buf/cmd/buf@v1.34.0
	go install -v github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.16.1
	go install -v github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.16.1
	go install -v google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0
	go install -v google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
	go install -v github.com/srikrsna/protoc-gen-gotag@v1.0.1

.PHONY: run
run: gen ## starting application and dependency services
# translate `SWAGGER_USE=true GIT_USE=true HTTP_PORT=8082 GRPC_PORT=10002 CMD_RUN="run" make run` => ./intmax2-withdrawal run
	go run $(GOFLAGS) ./cmd ${CMD_RUN}

.PHONY: cp
cp:
	cp -f build/env.docker.withdrawal-server.example build/env.docker.withdrawal-server

.PHONY: up
up: cp ## starting application and dependency services
	docker compose -f build/docker-compose.yml up

.PHONY: build-up
build-up: cp down cp ## rebuilding containers and starting application and dependency services
	docker compose -f build/docker-compose.yml build $(re_build_arg)
	docker compose -f build/docker-compose.yml up

.PHONY: start-build-up
start-build-up: down ## rebuilding containers and starting application and dependency services
	make cp
	docker compose -f build/docker-compose.yml up -d intmax2-withdrawal-postgres
	docker compose -f build/docker-compose.yml up -d intmax2-withdrawal-ot-collector

.PHONY: down
down: cp
	docker compose -f build/docker-compose.yml down
	rm -f build/env.docker.withdrawal-server

.PHONY: clean-all
clean-all: cp down
	rm -f build/env.docker.withdrawal-server
	rm -rf build/sql_dbs/intmax2-withdrawal-postgres

.PHONY: lint
lint:
	buf lint api/proto/withdrawal_service
	golangci-lint run --timeout=20m --fix ./...