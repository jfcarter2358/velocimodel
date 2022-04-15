#! /usr/bin/env bash

export MODEL_MANAGER_DB_USERNAME="ceresdb"
export MODEL_MANAGER_DB_PASSWORD="ceresdb"
export MODEL_MANAGER_DB_HOST="localhost"
export MODEL_MANAGER_DB_NAME="velocimodel"
export MODEL_MANAGER_DB_PORT="7437"
export MODEL_MANAGER_HTTP_PORT="9003"
export MODEL_MANGER_API_SERVER_URL="http://localhost:9004"

RUN_DIR=$(dirname $0)

pushd "${RUN_DIR}"
./model-manager
popd
