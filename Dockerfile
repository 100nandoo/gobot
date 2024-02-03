FROM golang:1.21.6-alpine AS builder
WORKDIR /build
ADD go.mod .
COPY . .
RUN go build -ldflags="-s -w" -o gobot main.go

FROM alpine
WORKDIR /app
COPY --from=builder /build/gobot /app/gobot
CMD ["./gobot"]
