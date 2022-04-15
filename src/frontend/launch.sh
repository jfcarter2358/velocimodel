#! /usr/bin/env bash

export FRONTEND_HTTP_PORT="9000"
export FRONTEND_API_SERVER_URL="http://localhost:9004"

RUN_DIR=$(dirname $0)

pushd "${RUN_DIR}"
./frontend
popd
