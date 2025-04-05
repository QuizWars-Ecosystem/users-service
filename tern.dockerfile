FROM golang:1.24.2-alpine AS builder

RUN go install github.com/jackc/tern/v2@latest

COPY ./migrations /migrations

ENTRYPOINT ["tern"]