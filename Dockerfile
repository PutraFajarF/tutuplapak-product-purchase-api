FROM golang:1.23-alpine AS builder

RUN apk add --no-cache --update bash git openssh curl vim busybox-extras

WORKDIR /app
COPY . /app

RUN go mod tidy
RUN go build -o cmd/main ./cmd

# Run stage
FROM alpine:3.20
WORKDIR /app/cmd
COPY --from=builder /app/cmd/main .

RUN apk add --no-cache --update tzdata
USER nobody
CMD ["./main"]