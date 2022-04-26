BIN_CALENDAR := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_sender"
CALENDAR_IMG="calendar:develop"
SHEDULER_IMG="scheduler:develop"
SENDER_IMG="sender:develop"
INTEG_TEST_ING="calendar_tests:integration"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

.PHONY: build run version test lint postgres createdb dropdb migrateup migratedown sqlc-gen api-gen generate \
build-calend-img build-sched-img build-send-img build-integ-tests-img \
run-calend-img run-sched-img run-send-img \
integration-tests up down 

build:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar; \
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler; \
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(BIN_CALENDAR) -config ./configs/calendar_config.toml -log ./log/logs.log

build-calend-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(CALENDAR_IMG) \
		-f build/Dockerfile .

build-sched-img:
	docker build \
		-t ${SHEDULER_IMG} \
		-f internal/amqp/publisher/Dockerfile .

build-send-img:
	docker build \
		-t ${SENDER_IMG} \
		-f internal/amqp/consumer/Dockerfile .

build-integ-tests-img:
	docker build \
		-t ${INTEG_TEST_ING} \
		-f tests/integration/integration-tests.Dockerfile .

run-calend-img: build-calend-img
	docker run $(CALENDAR_IMG)

run-sched-img: build-sched-img
	docker run ${SHEDULER_IMG}

run-send-img: build-send-img
	docker run ${SENDER_IMG}

version: build
	$(BIN) version

test:
	go test -race ./internal/storage/memory ./api

run-integration-tests:
	cd deployments; docker-compose -f docker-compose.yml -f docker-compose.test.yml up \
		--exit-code-from integration_tests

integration-tests: run-integration-tests
	cd deployments; docker-compose down --remove-orphans

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

createdb:
	docker exec -it pgOtus createdb --username=alex --owner=alex calendar

dropdb:
	docker exec -it pgOtus dropdb calendar

migrateup:
	migrate -path ./migrations -database "postgresql://alex:secret@localhost:5432/calendar?sslmode=disable" -verbose up

migratedown:
	migrate -path ./migrations -database "postgresql://alex:secret@localhost:5432/calendar?sslmode=disable" -verbose down

sqlc-gen:
	cd internal/storage/sql/sqlc/config; sqlc generate

api-gen:
	cd api; buf generate

up:
	cd deployments; docker-compose up -d

down: 
	cd deployments; docker-compose down --remove-orphans

generate: api-gen