module github.com/layer5io/meshery-istio

go 1.16

replace (
	// github.com/googleapis/gnostic/openapiv2 => github.com/google/gnostic/openapiv2 v0.5.5
	//github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200723152044-916f10574334
	github.com/spf13/afero => github.com/spf13/afero v1.5.1 // Until viper bug is resolved #1161
	go.opentelemetry.io/otel v1.3.0 => go.opentelemetry.io/otel v0.11.0
	go.opentelemetry.io/otel/sdk v1.3.0 => go.opentelemetry.io/otel/sdk v0.11.0
	gopkg.in/ini.v1 => github.com/go-ini/ini v1.62.0
//golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6
)

require (
	github.com/aspenmesh/istio-vet v0.0.0-20200806222806-9c8e9a962b9f
	github.com/layer5io/meshery-adapter-library v0.5.6
	github.com/layer5io/meshkit v0.5.26
	github.com/layer5io/service-mesh-performance v0.3.4
	gopkg.in/yaml.v2 v2.4.0
	istio.io/client-go v1.14.0
	k8s.io/apimachinery v0.24.2
	k8s.io/client-go v0.24.2
)

require (
	helm.sh/helm/v3 v3.9.0 // indirect
	k8s.io/kubectl v0.24.1 // indirect
)

require github.com/containerd/containerd v1.6.6 // indirect
