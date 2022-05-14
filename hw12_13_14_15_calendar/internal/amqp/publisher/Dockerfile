FROM golang:1.17.8-alpine as build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o scheduler cmd/scheduler/main.go

FROM alpine:3.14

ENV BIN_FILE "/app/scheduler"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /configs/scheduler_config.toml
COPY ./configs/scheduler_config.toml ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}
