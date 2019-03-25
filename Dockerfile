FROM golang:alpine as builder
RUN apk update \
    && apk add --virtual build-dependencies \
        build-base \
        gcc \
        wget \
        git \
    && apk add \
        bash

RUN mkdir /src
COPY go.mod /src
COPY go.sum /src
COPY *.go /src/
COPY config/* /src/config/
COPY ec2Query/* /src/ec2Query/
COPY ssh/* /src/ssh/

WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o main .

# Final image
FROM alpine:3.8
RUN apk add --no-cache --virtual .build-deps openssl
RUN apk add --no-cache ca-certificates
RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf
COPY --from=builder /src/main /main
RUN touch /config.yml
RUN touch /credentials
RUN touch /config.yml
RUN touch /sshKey
COPY entrypoint.sh ./entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]

