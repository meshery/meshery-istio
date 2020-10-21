#!/bin/sh

set -e

if ! /tmp/istio/istio-$ISTIO_VERSION/bin/istioctl x uninstall --purge -y; then
  exit 1
fi
