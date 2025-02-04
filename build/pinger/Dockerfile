# Собираем в гошке
FROM amd64/golang:1.23 AS build

WORKDIR /go/src/

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY .. /go/src/

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
RUN CGO_ENABLED=0 go build -o /opt/pinger/pinger-app ./public/pinger/...

# На выходе тонкий образ
FROM amd64/alpine:3.9

COPY --from=build /opt/pinger/pinger-app /opt/pinger/pinger-app

ENV APP_NAME=pinger

CMD ["/opt/pinger/pinger-app"]
