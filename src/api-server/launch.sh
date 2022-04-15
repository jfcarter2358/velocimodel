#! /usr/bin/env bash

export API_SERVER_HTTP_PORT="9004"
export API_SERVER_SERVICE_MANAGER_URL="http://localhost:9001"

RUN_DIR=$(dirname $0)

pushd "${RUN_DIR}"
./api-server
popd
