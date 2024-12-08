FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build ./cmd/notify

EXPOSE 8080

CMD ["./notify"]