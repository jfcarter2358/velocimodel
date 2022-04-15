#! /usr/bin/env bash

export ASSET_MANAGER_DB_USERNAME="ceresdb"
export ASSET_MANAGER_DB_PASSWORD="ceresdb"
export ASSET_MANAGER_DB_HOST="localhost"
export ASSET_MANAGER_DB_NAME="velocimodel"
export ASSET_MANAGER_DATA_PATH="/home/asset-manager/data/assets"
export ASSET_MANAGER_DB_PORT="7437"
export ASSET_MANAGER_HTTP_PORT="9002"
export ASSET_MANGER_API_SERVER_URL="http://localhost:9004"

RUN_DIR=$(dirname $0)

pushd "${RUN_DIR}"
./asset-manager
popd
