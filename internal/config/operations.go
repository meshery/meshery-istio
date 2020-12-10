package config

import (
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/common"
	"github.com/layer5io/meshery-adapter-library/meshes"
)

var (
	ServiceName = "service_name"
)

func getOperations(dev adapter.Operations) adapter.Operations {
	versions, _ := getLatestReleaseNames(3)

	// Add Istio networking resources to sample applications
	dev[common.BookInfoOperation].Templates = append(dev[common.BookInfoOperation].Templates, "file://templates/bookinfo-gateway.yaml")
	dev[common.HTTPBinOperation].Templates = append(dev[common.HTTPBinOperation].Templates, "file://templates/httpbin-gateway.yaml")
	dev[common.ImageHubOperation].Templates = append(dev[common.ImageHubOperation].Templates, "file://templates/imagehub-gateway.yaml")
	dev[common.EmojiVotoOperation].Templates = append(dev[common.EmojiVotoOperation].Templates, "file://templates/emojivoto-gateway.yaml")

	dev[IstioOperation] = &adapter.Operation{
		Type:                 int32(meshes.OpCategory_INSTALL),
		Description:          "Istio Service Mesh",
		Versions:             versions,
		AdditionalProperties: map[string]string{},
	}

	dev[LabelNamespace] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Automatic Sidecar Injection",
	}

	dev[PrometheusAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Add-on: Prometheus",
		Templates: []adapter.Template{
			"https://raw.githubusercontent.com/istio/istio/master/samples/addons/prometheus.yaml",
		},
		AdditionalProperties: map[string]string{
			ServiceName: "prometheus",
			PatchFile:   "file://templates/service-patch.json",
		},
	}

	dev[GrafanaAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Add-on: Grafana",
		Templates: []adapter.Template{
			"https://raw.githubusercontent.com/istio/istio/master/samples/addons/grafana.yaml",
		},
		AdditionalProperties: map[string]string{
			ServiceName: "grafana",
			PatchFile:   "file://templates/service-patch.json",
		},
	}

	dev[KialiAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Add-on: Kiali",
		Templates: []adapter.Template{
			"https://raw.githubusercontent.com/istio/istio/master/samples/addons/kiali.yaml",
		},
		AdditionalProperties: map[string]string{
			ServiceName: "kiali",
			PatchFile:   "file://templates/service-patch.json",
		},
	}

	dev[JaegerAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Add-on: Jaeger",
		Templates: []adapter.Template{
			"https://raw.githubusercontent.com/istio/istio/master/samples/addons/jaeger.yaml",
		},
		AdditionalProperties: map[string]string{
			ServiceName: "jaeger",
			PatchFile:   "file://templates/service-patch.json",
		},
	}

	dev[ZipkinAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Add-on: Zipkin",
		Templates: []adapter.Template{
			"https://raw.githubusercontent.com/istio/istio/master/samples/addons/extras/zipkin.yaml",
		},
		AdditionalProperties: map[string]string{
			ServiceName: "zipkin",
			PatchFile:   "file://templates/service-patch.json",
		},
	}

	dev[IstioVetOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_VALIDATE),
		Description: "Analyze Running Configuration",
	}

	dev[EnvoyFilterOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Envoy Filter for Image Hub",
		Versions:    adapter.NoneVersion,
		Templates: []adapter.Template{
			"file://templates/imagehub-filter.yaml",
		},
		AdditionalProperties: map[string]string{
			ServiceName: "api-v1",
			PatchFile:   "file://templates/imagehub-patch.json",
		},
	}

	dev[DenyAllPolicyOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Policy: Deny-All",
		Templates: []adapter.Template{
			"file://templates/policy-denyall.yaml",
		},
	}

	dev[StrictMTLSPolicyOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Policy: Strict MTLS",
		Templates: []adapter.Template{
			"file://templates/policy-strict.yaml",
		},
	}

	dev[MutualMTLSPolicyOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Policy: Mutual MTLS",
		Templates: []adapter.Template{
			"file://templates/policy-mutual.yaml",
		},
	}

	dev[DisableMTLSPolicyOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Policy: Disable MTLS",
		Templates: []adapter.Template{
			"file://templates/policy-disable.yaml",
		},
	}

	return dev
}
