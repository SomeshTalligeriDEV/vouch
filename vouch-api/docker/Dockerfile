# ---- build stage ----
FROM golang:1.22-alpine AS build
WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build both binaries statically.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/api ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/worker ./cmd/worker

# ---- runtime stage ----
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates && \
    adduser -D -u 10001 vouch
COPY --from=build /out/api /app/api
COPY --from=build /out/worker /app/worker
USER vouch

EXPOSE 8080
# Default to the API; the worker process overrides the command in Fly/compose.
ENTRYPOINT ["/app/api"]
