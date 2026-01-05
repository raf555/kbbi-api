FROM golang:1.25 AS builder

WORKDIR /app

COPY pkg/kbbi/ ./pkg/kbbi/
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o main ./cmd/kbbi

# --------------------------------

FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /app/main .
ENTRYPOINT ["./main"]
CMD []
