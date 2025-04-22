FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o main ./src/cmd/

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates libc6-compat

# Copiar específicamente los archivos necesarios
COPY --from=builder /app/src/database/ddl.sql .
COPY --from=builder /app/main .

RUN mkdir -p /app/data

ENV DB_PATH=/app/data/series.db
ENV DDL_PATH=/app/ddl.sql

EXPOSE 8080


CMD ["./main"]
