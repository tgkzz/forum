FROM golang:1.20.1-alpine3.16 AS builder
RUN apk add build-base gcc
LABEL project-name = "forum"
WORKDIR /app
COPY . .
RUN go build -o forum main.go
FROM alpine:3.16 
WORKDIR /app
COPY --from=builder /app /app
CMD ["./forum"]
EXPOSE 4000
