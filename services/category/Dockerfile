FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/category ./services/category
COPY shared ./shared

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/category-service ./services/category/cmd/api


FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/category-service /app/category-service

EXPOSE 8082

USER nonroot:nonroot
ENTRYPOINT ["/app/category-service"]
