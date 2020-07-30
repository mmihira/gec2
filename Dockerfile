FROM golang:alpine as builder
RUN apk update \
    && apk add --virtual build-dependencies \
        build-base \
        gcc \
        wget \
        git \
    && apk add \
        bash

COPY . /
RUN rm -rf ./src/deploy_context
RUN rm -rf ./src/examples
RUN rm -rf ./src/scripts
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o main .

# Final image
FROM alpine:3.8
RUN apk add --no-cache --virtual .build-deps openssl
RUN apk add --no-cache ca-certificates
RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf
COPY --from=builder /src/main /main
RUN touch /config.yaml
RUN touch /credentials
RUN touch /sshKey
RUN mkdir context
RUN mkdir logs
RUN mkdir roles
ENTRYPOINT ["./main"]
