FROM golang:1.18.0-alpine

WORKDIR /service-manager-build
COPY src/service-manager /service-manager-build
RUN env GOOS=linux CGO_ENABLED=0 go build -v -o service-manager

FROM alpine:latest

RUN adduser --disabled-password service-manager

WORKDIR /home/service-manager

COPY --from=0 /service-manager-build/service-manager ./

RUN apk update \
    && apk add bash

SHELL ["/bin/bash", "-c"]

RUN mkdir /home/service-manager/data
COPY data /home/service-manager/data

ADD src/service-manager/start-service-manager.sh /home/service-manager/start-service-manager.sh
RUN chmod u+x /home/service-manager/start-service-manager.sh

RUN chown -R service-manager:service-manager /home/service-manager

USER service-manager

CMD ["./start-service-manager.sh"]
