FROM golang:1.18.3-alpine3.16 AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY main.go ./
RUN go build -o main main.go

FROM alpine:3.16
WORKDIR /app
COPY ./templates ./
ADD templates ./templates
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["/app/main"]