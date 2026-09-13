FROM node:22-alpine AS html-minifier

WORKDIR /app

COPY assets/view/index.html .
RUN npx --yes html-minifier-terser@7.2.0 index.html \
	--collapse-whitespace \
	--minify-css \
	--minify-js \
	--remove-comments \
	--remove-redundant-attributes \
	--remove-script-type-attributes \
	--use-short-doctype \
	--output index.min.html

# --------------------------------

FROM golang:1.27 AS builder

ARG VERSION

WORKDIR /app

COPY pkg/kbbi/ ./pkg/kbbi/
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w -X github.com/raf555/kbbi-api/internal/version.Version=${VERSION}" -o main ./cmd/kbbi

# --------------------------------

FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=html-minifier /app/index.min.html ./assets/view/index.html
ENTRYPOINT ["./main"]
CMD []
