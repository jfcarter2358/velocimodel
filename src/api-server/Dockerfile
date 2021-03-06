FROM golang:1.18.1-alpine

RUN apk update && apk add git

WORKDIR /api-server-build
COPY src/api-server /api-server-build
RUN env GOOS=linux CGO_ENABLED=0 go build -v -o api-server

FROM alpine:latest

RUN adduser --disabled-password api-server

WORKDIR /home/api-server

COPY --from=0 /api-server-build/api-server ./

RUN apk update \
    && apk add bash

SHELL ["/bin/bash", "-c"]

RUN mkdir /home/api-server/data
COPY data/api-server /home/api-server/data

ADD src/api-server/start-api-server.sh /home/api-server/start-api-server.sh
RUN chmod u+x /home/api-server/start-api-server.sh

RUN chown -R api-server:api-server /home/api-server

USER api-server

CMD ["./start-api-server.sh"]
