FROM golang:1.17.8 as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux \
    go build -o sender cmd/sender/main.go

FROM alpine:3.14

ENV BIN_FILE "/app/sender"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /configs/sender_config.toml
COPY ./configs/sender_config.toml ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}