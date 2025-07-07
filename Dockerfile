FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates postgresql-client
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrate.sh .
RUN chmod +x migrate.sh

EXPOSE 8080

CMD ["./main"]