# Собираем в гошке
FROM amd64/golang:1.23 AS build

WORKDIR /go/src/

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /go/src/

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
RUN CGO_ENABLED=0 go build -o /opt/order/order-app ./public/...

# Собираем migrator
RUN git clone https://github.com/golang-migrate/migrate

WORKDIR /go/src/migrate/

RUN git checkout v4.18.2
RUN CGO_ENABLED=0 go build -o /go/src/db -ldflags='-X main.Version=4.18' -tags 'github postgres' ./cmd/migrate

# На выходе тонкий образ
FROM amd64/alpine:3.9

COPY --from=build /opt/order/order-app /opt/order/order-app
COPY --from=build /go/src/db /opt/db

ENV APP_NAME=order

