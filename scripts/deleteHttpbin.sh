#!/bin/sh

set -e

if ! kubectl delete -f /tmp/istio/samples/httpbin/httpbin.yaml; then
  exit 1
fi
