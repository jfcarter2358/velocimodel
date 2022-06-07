FROM golang:1.18.1-alpine

RUN apk update && apk add git

WORKDIR /asset-manager-build
COPY src/asset-manager /asset-manager-build
RUN env GOOS=linux CGO_ENABLED=0 go build -v -o asset-manager

FROM alpine:latest

RUN adduser --disabled-password asset-manager

WORKDIR /home/asset-manager

COPY --from=0 /asset-manager-build/asset-manager ./

RUN apk update \
    && apk add \
    bash \
    git \
    openssh

SHELL ["/bin/bash", "-c"]

RUN mkdir /home/asset-manager/.ssh
RUN mkdir /home/asset-manager/data
COPY data/asset-manager /home/asset-manager/data

ADD src/asset-manager/start-asset-manager.sh /home/asset-manager/start-asset-manager.sh
RUN chmod u+x /home/asset-manager/start-asset-manager.sh

RUN chown -R asset-manager:asset-manager /home/asset-manager

USER asset-manager

CMD ["./start-asset-manager.sh"]
