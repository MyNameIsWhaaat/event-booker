FROM golang:1.25-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/api ./cmd/eventbooker
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/worker ./cmd/eventbooker-worker

FROM alpine:3.20

WORKDIR /app
RUN apk add --no-cache ca-certificates
RUN mkdir -p /data && chmod -R 777 /data

COPY --from=build /bin/api /app/api
COPY --from=build /bin/worker /app/worker
COPY internal/web /app/internal/web

EXPOSE 8080