module github.com/layer5io/meshery-istio

go 1.13

replace github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200723152044-916f10574334

require (
	github.com/Azure/go-autorest/autorest/adal v0.9.0 // indirect
	github.com/aspenmesh/istio-vet v0.0.0-20200806222806-9c8e9a962b9f
	github.com/ghodss/yaml v1.0.0
	github.com/golang/protobuf v1.4.2
	github.com/gophercloud/gophercloud v0.4.0 // indirect
	github.com/layer5io/gokit v0.1.12
	github.com/layer5io/learn-layer5/smi-conformance v0.0.0-20200916172547-79cf11334bd7
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	golang.org/x/net v0.0.0-20200927032502-5d4f70055728
	google.golang.org/grpc v1.32.0
	helm.sh/helm/v3 v3.4.0
	istio.io/client-go v0.0.0-20200708142230-d7730fd90478
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	rsc.io/letsencrypt v0.0.3 // indirect
)
