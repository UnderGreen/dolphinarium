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

function build_and_push_galera() {
    announce-step "Build and push mysql-galera image"

    local registry=$(k8s-service-endpoint kube-registry 5000 kube-system)
    REGISTRY="$registry" make -C "$PROJECT_DIR/galera/" build
    REGISTRY="$registry" make -C "$PROJECT_DIR/galera/" push
}

function main() {
    deploy_registry
    build_and_push_galera
}

main "$@"
