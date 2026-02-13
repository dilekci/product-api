FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/product ./services/product
COPY shared ./shared

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/product-service ./services/product/cmd/api


FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/product-service /app/product-service

EXPOSE 8081

USER nonroot:nonroot
ENTRYPOINT ["/app/product-service"]
