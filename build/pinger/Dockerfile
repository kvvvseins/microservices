# Собираем в гошке
FROM amd64/golang:1.23 AS build

ENV BIN_FILE=/opt/pinger/pinger-app
ENV CODE_DIR=/go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY .. ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} ./public/pinger/...

# На выходе тонкий образ
FROM amd64/alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="pinger"
LABEL MAINTAINERS="vladimir.kostin@vitech.team"

ENV BIN_FILE=/opt/pinger/pinger-app
COPY --from=build ${BIN_FILE} ${BIN_FILE}

CMD APP_NAME=pinger ${BIN_FILE}
