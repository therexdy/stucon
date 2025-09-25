FROM docker.io/library/golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./cmd/stucon ./cmd/main.go

EXPOSE 8080

CMD ["./cmd/stucon"]
