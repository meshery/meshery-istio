module github.com/layer5io/meshery-istio

go 1.13

replace github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200723152044-916f10574334

require (
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/layer5io/meshery-adapter-library v0.1.7
	github.com/layer5io/meshkit v0.1.28
	github.com/onsi/ginkgo v1.13.0 // indirect
	golang.org/x/net v0.0.0-20200927032502-5d4f70055728 // indirect
	google.golang.org/grpc v1.32.0 // indirect
	k8s.io/apimachinery v0.19.4
)
