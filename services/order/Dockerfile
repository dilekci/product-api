FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/order ./services/order
COPY shared ./shared

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/order-service ./services/order/cmd/api


FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/order-service /app/order-service

EXPOSE 8082

USER nonroot:nonroot
ENTRYPOINT ["/app/order-service"]
