.PHONY: all build run migrate tidy swag image help

CURRENT_DIR := $(shell pwd)
APP := $(shell basename ${CURRENT_DIR})
APP_CMD_DIR := ${CURRENT_DIR}/cmd/app
TAG := latest
ENV_TAG := latest

-include ./conf/dev.env

POSTGRESQL_URL := 'postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=${POSTGRES_SSLMODE}'

build:
	@rm -f ${CURRENT_DIR}/bin/${APP}
	@CGO_ENABLED=1 GOOS=linux go build -tags musl -ldflags="-X main.Version=${TAG}" -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

migrate-up:
	@migrate -database ${POSTGRESQL_URL} -path migrations/postgres up

migrate-down:
	@migrate -database ${POSTGRESQL_URL} -path migrations/postgres down

migrate-new: # make migrate-new name=file_name
	@migrate create -ext sql -dir migrations/postgres -seq $(name)

build-image:
	@docker build --rm -t ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG} .
	@docker tag ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG} ${REGISTRY}/${PROJECT_NAME}/${APP}:${ENV_TAG}

push-image:
	@docker push ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG}
	@docker push ${REGISTRY}/${PROJECT_NAME}/${APP}:${ENV_TAG}

swag-init:
	@swag init -g internal/transport/http/handler.go -o docs --parseVendor --parseDependency 1 --parseInternal 1 --parseDepth 1

# export APP_ENVIRONMENT=prod
run: 
	@go run ${APP_CMD_DIR}/main.go

.DEFAULT_GOAL := run