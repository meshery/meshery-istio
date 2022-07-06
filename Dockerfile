FROM golang:1.17 as builder

ARG VERSION
ARG GIT_COMMITSHA
WORKDIR /build
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download
# Copy the go source
COPY main.go main.go
COPY internal/ internal/
COPY istio/ istio/
# Build
COPY build/ build/
RUN GOPROXY=direct,https://proxy.golang.org CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -ldflags="-w -s -X main.version=$VERSION -X main.gitsha=$GIT_COMMITSHA" -a -o meshery-istio main.go

FROM alpine:3.15 as jsonschema-util
RUN apk add --no-cache curl
WORKDIR /
RUN UTIL_VERSION=$(curl -L -s https://api.github.com/repos/layer5io/kubeopenapi-jsonschema/releases/latest | \
	grep tag_name | sed "s/ *\"tag_name\": *\"\\(.*\\)\",*/\\1/" | \
	grep -v "rc\.[0-9]$"| head -n 1 ) \
	&& curl -L https://github.com/layer5io/kubeopenapi-jsonschema/releases/download/${UTIL_VERSION}/kubeopenapi-jsonschema -o kubeopenapi-jsonschema \
	&& chmod +x /kubeopenapi-jsonschema

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/nodejs:16
ENV DISTRO="debian"
ENV GOARCH="amd64"
ENV SERVICE_ADDR="meshery-istio"
ENV MESHERY_SERVER="http://meshery:9081"
COPY templates/ ./templates
WORKDIR /
COPY --from=builder /build/meshery-istio .
COPY --from=jsonschema-util /kubeopenapi-jsonschema /root/.meshery/bin/kubeopenapi-jsonschema
ENTRYPOINT ["/meshery-istio"]
