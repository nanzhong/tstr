#!/usr/bin/env bash

set -eo pipefail

deploy_path="${BASH_SOURCE%/*}/../deploy"

options=$(getopt -l "environment:,image-tag:,pg-dsn:,ui-access-token:" -o "e:t:d:u:" -a -- "$@")
eval set -- "$options"

environment=""
image_tag=""
pg_dsn=""
ui_access_token=""

while true
do
  case $1 in
    -t|--image-tag)
      shift
      image_tag=$1
      ;;
    -e|--environment)
      shift
      environment=$1
      ;;
    -d|--pg-dsn)
      shift
      pg_dsn=$1
      ;;
    -u|--ui-access-token)
      shift
      ui_access_token=$1
      ;;
    --)
      shift
      break;;
  esac
  shift
done

cd "$deploy_path/$environment"

kustomize edit set image "nanzhong/tstr:$image_tag"
kustomize edit add secret tstr \
  --from-literal="pg_dsn=$pg_dsn" \
  --from-literal="ui_access_token=$ui_access_token"
kustomize build | kubectl apply -f -
kubectl -n "$environment" rollout status deployment/tstr --timeout 60s
