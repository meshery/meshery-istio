#!/bin/sh

set -e

if ! kubectl apply -f /tmp/istio/samples/httpbin/httpbin.yaml; then
  exit 1
fi
