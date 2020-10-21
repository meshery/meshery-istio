#!/bin/sh

set -e

: "${ISTIO_VERSION:=}"
: "${ARCH:=amd64}"
: "${DISTRO:=linux}"

OS=`uname -s`
if [ "$OS" = "Linux" ]; then
  DISTRO="linux"
  URL="https://github.com/istio/istio/releases/download/$ISTIO_VERSION/istio-$ISTIO_VERSION-$DISTRO-$ARCH.tar.gz"
elif [ "$OS" = "Darwin" ]; then
  DISTRO="osx"
  URL="https://github.com/istio/istio/releases/download/$ISTIO_VERSION/istio-$ISTIO_VERSION-$DISTRO.tar.gz"
else
  exit 1
fi

if [ -z "$DISTRO" ]; then
  exit 2
fi


if ! type "grep" > /dev/null 2>&1; then
  exit 3;
fi
if ! type "curl" > /dev/null 2>&1; then
  exit 4;
fi
if ! type "tar" > /dev/null 2>&1; then
  exit 5;
fi
if ! type "gzip" > /dev/null 2>&1; then
  exit 6;
fi

if ! curl -s --head $URL | head -n 1 | grep "HTTP/1.[01] [23].." > /dev/null; then
  exit 7;
fi

if ! curl -L "$URL" | tar xz; then
  exit 8;
fi

if [ "$ISTIO_MODE" = "operator" ]; then
  if ! ./istio-$ISTIO_VERSION/bin/istioctl operator init; then
    exit 9;
  fi
else 
  if ! ./istio-$ISTIO_VERSION/bin/istioctl install --set profile=$ISTIO_PROFILE --set meshConfig.accessLogFile=/dev/stdout; then
  	exit 10;
  fi
fi

if ! mv istio-$ISTIO_VERSION /tmp/istio/.; then
  echo "Already installed"
fi


