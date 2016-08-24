#!/usr/bin/env bash

set -e

SCRIPT_DIR=$(dirname $0)
PROJECT_DIR="$SCRIPT_DIR"

source "$PROJECT_DIR/tools/shell-functions.sh"


function deploy_registry() {
	announce-step "Deploying docker registry"

	make -C "$PROJECT_DIR/registry/" start

	local registry="http://$(k8s-service-endpoint kube-registry 5000 kube-system)"
	wait-for-http "$registry" 5 30
}

function main() {
	deploy_registry
}

main "$@"
