#! /usr/bin/env bash

SLEEP_AMOUNT="${AUTH_MANAGER_SLEEP:-"0"}"
echo "Sleeping for ${SLEEP_AMOUNT} seconds"
sleep "${SLEEP_AMOUNT}"

RUN_DIR=$(dirname $0)

pushd "${RUN_DIR}"
./auth-manager
popd
