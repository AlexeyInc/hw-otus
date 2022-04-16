BIN_CALENDAR := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"
SHED_IMG="sched:dev"
SEND_IMG="send:dev"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar; \
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./internal/amqp/publisher; \
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./internal/amqp/consumer

run: build
	$(BIN_CALENDAR) -config ./configs/calendar_config.toml -log ./log/logs.log

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

build-sched-img:
	docker build \
		-t ${SHED_IMG} \
		-f internal/amqp/publisher/docker/Dockerfile .

build-send-img:
	docker build \
		-t ${SEND_IMG} \
		-f internal/amqp/consumer/docker/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

run-sched-img: build-sched-img
	docker run ${SHED_IMG}

run-send-img: build-send-img
	docker run ${SEND_IMG}

version: build
	$(BIN) version

test:
	go test -race ./internal/storage/memory ./internal/storage/sql/sqlc

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

.PHONY: build run build-img run-img version test lint postgres createdb dropdb migrateup migratedown sqlc-gen api-gen build-sched-img build-send-img run-sched-img run-send-img
