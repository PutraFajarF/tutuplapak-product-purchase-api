FROM golang:1.19.5-alpine as builder

RUN \
    apk add --no-cache bash git openssh && \
    apk --no-cache add curl && \
    apk --no-cache add vim && \
    apk --no-cache add procps-dev && \
    apk --no-cache add busybox-extras

ADD ./ /app
RUN cd /app
WORKDIR /app/cmd

RUN go build -o main .

# Run stage
FROM alpine:3.13
WORKDIR /app/cmd
COPY --from=builder /app/cmd/main .

RUN apk add --no-cache tzdata
CMD ["/app/cmd/main"]