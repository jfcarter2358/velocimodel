FROM golang:1.18.0-alpine

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
    git

SHELL ["/bin/bash", "-c"]

ADD src/asset-manager/start-asset-manager.sh /home/asset-manager/start-asset-manager.sh
RUN chmod u+x /home/asset-manager/start-asset-manager.sh

RUN chown -R asset-manager:asset-manager /home/asset-manager

USER asset-manager

CMD ["./start-asset-manager.sh"]
