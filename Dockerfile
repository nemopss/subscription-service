FROM --platform=linux/arm64 golang:1.24.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM --platform=linux/arm64 alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
RUN ls -l /app
CMD ["./main"]
