FROM golang:1.18.1-alpine

RUN apk update && apk add git

WORKDIR /auth-manager-build
COPY src/auth-manager /auth-manager-build
RUN env GOOS=linux CGO_ENABLED=0 go build -v -o auth-manager

FROM alpine:latest

RUN adduser --disabled-password auth-manager

WORKDIR /home/auth-manager

COPY --from=0 /auth-manager-build/auth-manager ./

RUN apk update \
    && apk add \
    bash

SHELL ["/bin/bash", "-c"]

RUN mkdir /home/auth-manager/data
COPY data/auth-manager /home/auth-manager/data

ADD src/auth-manager/start-auth-manager.sh /home/auth-manager/start-auth-manager.sh
RUN chmod u+x /home/auth-manager/start-auth-manager.sh

COPY src/auth-manager/ui /home/auth-manager

RUN chown -R auth-manager:auth-manager /home/auth-manager

USER auth-manager

ENTRYPOINT ["./start-auth-manager.sh"]
