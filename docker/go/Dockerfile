# Базовый образ
FROM golang:1.19

WORKDIR /usr/src/app

COPY ../.. go.mod ./
COPY ../.. go.sum ./
RUN go mod download && go mod verify

COPY ../.. ./

EXPOSE 8080

RUN go build -v -o /usr/local/bin/app ./...

CMD ["app"]