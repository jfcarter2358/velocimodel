#! /usr/bin/env bash

export AUTH_MANAGER_CONFIG_PATH="/home/auth-manager/data/config.json"
export AUTH_MANAGER_SLEEP="5"

RUN_DIR=$(dirname $0)

pushd "${RUN_DIR}"
./frontend
popd
