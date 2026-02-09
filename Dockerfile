FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/product-api ./cmd/api


FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/product-api /app/product-api

EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/app/product-api"]
