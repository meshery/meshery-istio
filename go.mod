module github.com/layer5io/meshery-istio

go 1.13

replace github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200723152044-916f10574334

require (
	github.com/Azure/go-autorest/autorest/adal v0.9.0
	github.com/aspenmesh/istio-client-go v0.0.0-20200122202704-9695ccefca79
	github.com/aspenmesh/istio-vet v0.0.0-20190708155745-fac74e54e2e6
	github.com/cnf/structhash v0.0.0-20180104161610-62a607eb0224 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.4.2
	github.com/gophercloud/gophercloud v0.4.0 // indirect
	github.com/layer5io/learn-layer5/smi-conformance v0.0.0-20200806140644-44993b65fd37
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	golang.org/x/net v0.0.0-20200520004742-59133d7f0dd7
	google.golang.org/grpc v1.31.1
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v0.17.3
)
