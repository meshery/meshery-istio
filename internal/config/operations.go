package config

import (
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/meshes"
)

var (
	ServiceName = "service_name"
)

func getOperations(dev adapter.Operations) adapter.Operations {
	versions, _ := getLatestReleaseNames(3)

	dev[IstioOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_INSTALL),
		Description: "Istio Service Mesh",
		Versions:    versions,
		Templates: []adapter.Template{
			"templates/istio.yaml",
		},
		AdditionalProperties: map[string]string{},
	}

	dev[LabelNamespace] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Label Namespace for Automatic Sidecar Injection",
	}

	dev[PrometheusAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Prometheus Monitoring",
	}

	dev[GrafanaAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Grafana Dashboard",
	}

	dev[KialiAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Kiali Dashboard",
	}

	dev[JaegerAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Jaeger Dashboard",
	}

	dev[ZipkinAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Zipkin Dashboard",
	}

	dev[IstioVetOpertation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_VALIDATE),
		Description: "Analyze Running Configuration",
	}

	return dev
}
