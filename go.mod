module github.com/layer5io/meshery-istio

go 1.13

replace github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200723152044-916f10574334

require (
	github.com/aspenmesh/istio-vet v0.0.0-20200806222806-9c8e9a962b9f
	github.com/layer5io/meshery-adapter-library v0.1.7
	github.com/layer5io/meshkit v0.1.28
	github.com/onsi/ginkgo v1.13.0 // indirect
	golang.org/x/net v0.0.0-20200927032502-5d4f70055728 // indirect
	google.golang.org/grpc v1.32.0 // indirect
	istio.io/client-go v1.8.0
	k8s.io/apimachinery v0.18.12
	k8s.io/client-go v0.18.12
)
