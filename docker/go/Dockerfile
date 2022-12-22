# Базовый образ
FROM golang:1.19

RUN mkdir /app
WORKDIR /app

COPY ../.. go.mod ./
COPY ../.. go.sum ./
RUN go mod download && go mod verify

COPY ../.. ./

EXPOSE 8080

RUN go build -o main .

CMD ["/app/main"]