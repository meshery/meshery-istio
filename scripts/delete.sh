#!/bin/sh

set -e

if ! /tmp/istio/bin/istioctl x uninstall --purge -y; then
  exit 1
fi
