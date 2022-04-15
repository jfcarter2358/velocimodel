FROM golang:1.18.0-alpine

WORKDIR /model-manager-build
COPY src/model-manager /model-manager-build
RUN env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o model-manager

FROM alpine:latest

RUN adduser --disabled-password model-manager

WORKDIR /home/model-manager

COPY --from=0 /model-manager-build/model-manager ./

RUN apk update \
    && apk add bash

SHELL ["/bin/bash", "-c"]

ADD src/model-manager/start-model-manager.sh /home/model-manager/start-model-manager.sh
RUN chmod u+x /home/model-manager/start-model-manager.sh

RUN chown -R model-manager:model-manager /home/model-manager

USER model-manager

CMD ["./start-model-manager.sh"]
