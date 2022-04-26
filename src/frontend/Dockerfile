FROM golang:1.18.0-alpine

WORKDIR /frontend-build
COPY src/frontend /frontend-build
RUN env GOOS=linux CGO_ENABLED=0 go build -v -o frontend

FROM alpine:latest

RUN adduser --disabled-password frontend

WORKDIR /home/frontend

COPY --from=0 /frontend-build/frontend ./

RUN apk update \
    && apk add bash

SHELL ["/bin/bash", "-c"]

ADD src/frontend/start-frontend.sh /home/frontend/start-frontend.sh
RUN chmod u+x /home/frontend/start-frontend.sh

COPY src/frontend/ui-dist /home/frontend

RUN chown -R frontend:frontend /home/frontend

USER frontend

CMD ["./start-frontend.sh"]
