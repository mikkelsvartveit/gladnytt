ARG GO_VERSION=1

# Build stage
FROM golang:${GO_VERSION}-alpine as builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=1 go build -v -o ./main ./src


# Production stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main /main
COPY --from=builder /app/templates /templates
COPY --from=builder /app/static /static

CMD ["./main"]
