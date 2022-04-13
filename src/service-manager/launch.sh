#! /usr/bin/env bash

export SERVICE_MANAGER_DB_USERNAME="ceresdb"
export SERVICE_MANAGER_DB_PASSWORD="ceresdb"
export SERVICE_MANAGER_DB_HOST="localhost"
export SERVICE_MANAGER_DB_NAME="velocimodel"
export SERVICE_MANAGER_PARAMS_PATH="./data/config.json"
export SERVICE_MANAGER_SECRETS_PATH="./data/secrets.json"
export SERVICE_MANAGER_DB_PORT="7437"
export SERVICE_MANAGER_HTTP_PORT="9001"
export SERVICE_MANAGER_ENCRYPTION_KEY="abcdefghijklmnopqrstuvqxyz012345"

RUN_DIR=$(dirname $0)

pushd "${RUN_DIR}"
./service-manager
popd
