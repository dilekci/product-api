FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/user ./services/user
COPY shared ./shared

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/user-service ./services/user/cmd/api


FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/user-service /app/user-service

EXPOSE 8083

USER nonroot:nonroot
ENTRYPOINT ["/app/user-service"]
