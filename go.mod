module github.com/layer5io/meshery-istio

go 1.13

replace github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200723152044-916f10574334

require (
	cloud.google.com/go v0.46.3 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.0 // indirect
	github.com/aspenmesh/istio-vet v0.0.0-20200806222806-9c8e9a962b9f
	github.com/ghodss/yaml v1.0.0
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.4.2
	github.com/gophercloud/gophercloud v0.4.0 // indirect
	github.com/layer5io/gokit v0.1.12
	github.com/layer5io/learn-layer5/smi-conformance v0.0.0-20200916172547-79cf11334bd7
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	go.opencensus.io v0.22.3 // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a // indirect
	google.golang.org/genproto v0.0.0-20200731012542-8145dea6a485 // indirect
	google.golang.org/grpc v1.31.0
	google.golang.org/protobuf v1.25.0 // indirect
	helm.sh/helm/v3 v3.3.1
	istio.io/client-go v0.0.0-20200708142230-d7730fd90478
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
	rsc.io/letsencrypt v0.0.3 // indirect
)
