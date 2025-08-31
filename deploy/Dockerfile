FROM golang:1.22-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o reviewservice ./cmd/reviewService


# -------- Runtime Stage --------
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/reviewservice .

EXPOSE 8080
CMD ["./reviewservice"]
