#! /usr/bin/env bash

SLEEP_AMOUNT="${MODEL_MANAGER_SLEEP:-"0"}"
echo "Sleeping for ${SLEEP_AMOUNT} seconds"
sleep "${SLEEP_AMOUNT}"

RUN_DIR=$(dirname $0)

pushd "${RUN_DIR}"
./model-manager
popd
