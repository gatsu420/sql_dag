FROM golang:1.22.4-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod .
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o main

FROM alpine:latest
RUN apk add --no-cache graphviz
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/sample_query.sql . 
EXPOSE 8080
CMD ["./main"]
