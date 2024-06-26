#!/bin/bash

echo 'Starting raft node'

DEFAULT_NODE_BIN_PATH=.

DEFAULT_RPC_PORT=6666
DEFAULT_API_PORT=7777

# Define the directories
CLUSTER_DIR_ROOT=./cluster-data

DATA_PATH=${CLUSTER_DIR_ROOT}/log
METADATA_PATH=${CLUSTER_DIR_ROOT}/log/metadata.json
STATE_PATH=${CLUSTER_DIR_ROOT}/state/state.json
SNAPSHOT_PATH=${CLUSTER_DIR_ROOT}/snapshot

# Default join mode is set to cluster-file unless
# it is set to Kubernetes mode in which case it will
# look for the file containing labels
DEFAULT_JOIN_MODE='cluster-file'

ELECTION_TIMEOUT_MILLIS=2000
HEARTBEAT_INTERVAL_MILLIS=500

RPC_TIMEOUT_MILLIS=1250
API_TIMEOUT_MILLIS=1500
API_FWD_TIMEOUT_MILLIS=2000

MAX_CONN_RETRY_ATTEMPTS=3

# Common parameters for both modes of cluster formation
read -r -d '' NODE_ARGS <<EOM
    --api-port=${DEFAULT_API_PORT}
    --rpc-port=${DEFAULT_RPC_PORT}
    --log-entry-path=${DATA_PATH} 
    --log-metadata-path=${METADATA_PATH}
    --raft-state-path=${STATE_PATH} 
    --election-timeout=${ELECTION_TIMEOUT_MILLIS}
    --rpc-timeout=${RPC_TIMEOUT_MILLIS}
    --api-timeout=${API_TIMEOUT_MILLIS}
    --api-fwd-timeout=${API_FWD_TIMEOUT_MILLIS}
    --max-conn-retry-attempts=${MAX_CONN_RETRY_ATTEMPTS}
    --snapshot-path=${SNAPSHOT_PATH}
    --cluster-config-path=${CLUSTER_CONFIG_PATH}
EOM

NODE_ARGS="${NODE_ARGS} --join-mode=k8s --cluster-config-path=${CLUSTER_DIR_ROOT}/cluster"
touch arguments.txt
echo "Arguments for the node: ${NODE_ARGS}" > arguments.txt
