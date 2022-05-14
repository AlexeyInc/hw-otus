# Собираем в гошке
FROM golang:1.17.8-alpine as build

WORKDIR /app

# Кэшируем слои с модулями
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 GOOS=linux go build \
        -ldflags "$LDFLAGS" \
        -o calendar cmd/calendar/*

# На выходе тонкий образ
FROM alpine:3.14

ENV BIN_FILE "/app/calendar"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /config/config.toml
COPY ./configs/calendar_config.toml ${CONFIG_FILE}

ENV LOG_FILE /log/logs.log
COPY ./log/logs.log ${LOG_FILE}

EXPOSE 8080
EXPOSE 8081

CMD ${BIN_FILE} -config ${CONFIG_FILE} -log ${LOG_FILE}
