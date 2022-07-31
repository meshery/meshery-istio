# Copyright Meshery Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

include build/Makefile.core.mk
include build/Makefile.show-help.mk

#-----------------------------------------------------------------------------
# Environment Setup
#-----------------------------------------------------------------------------
BUILDER=buildx-multi-arch
ADAPTER=istio

v ?= 1.17.8 # Default go version to be used


#-----------------------------------------------------------------------------
# Docker-based Builds
#-----------------------------------------------------------------------------
.PHONY: docker docker-run lint proto-setup proto error test run run-force-dynamic-reg


## Lint check Golang
lint:
	golangci-lint run

## Retrieve protos
proto-setup:
	cd meshes
	wget https://raw.githubusercontent.com/layer5io/meshery/master/meshes/meshops.proto

## Generate protos
proto:	
	protoc -I meshes/ meshes/meshops.proto --go_out=plugins=grpc:./meshes/

## Build Adapter container image with "edge-latest" tag
docker:
	DOCKER_BUILDKIT=1 docker build -t layer5/meshery-$(ADAPTER):$(RELEASE_CHANNEL)-latest .

## Run Adapter container with "edge-latest" tag
docker-run:
	(docker rm -f meshery-$(ADAPTER)) || true
	docker run --name meshery-$(ADAPTER) -d \
	-p 10000:10000 \
	-e DEBUG=true \
	layer5/meshery-$(ADAPTER):$(RELEASE_CHANNEL)-latest

## Build and run Adapter locally
run:
	go$(v) mod tidy -compat=1.17; \
	DEBUG=true GOPROXY=direct GOSUMDB=off go run main.go

## Build and run Adapter locally; force component registration
run-force-dynamic-reg:
	FORCE_DYNAMIC_REG=true DEBUG=true GOPROXY=direct GOSUMDB=off go run main.go

## Run Meshery Error utility
error:
	go run github.com/layer5io/meshkit/cmd/errorutil -d . analyze -i ./helpers -o ./helpers

## Run Golang tests
test:
	export CURRENTCONTEXT="$(kubectl config current-context)" 
	echo "current-context:" ${CURRENTCONTEXT} 
	export KUBECONFIG="${HOME}/.kube/config"
	echo "environment-kubeconfig:" ${KUBECONFIG}
	GOPROXY=direct GOSUMDB=off GO111MODULE=on go test -v ./...
