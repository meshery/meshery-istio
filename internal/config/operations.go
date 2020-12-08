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

	dev[IstioOperation] = &adapter.Operation{
		Type:                 int32(meshes.OpCategory_INSTALL),
		Description:          "Istio Service Mesh",
		Versions:             versions,
		AdditionalProperties: map[string]string{},
	}

	// Add Istio networking resources to sample applications
	dev[common.BookInfoOperation].Templates = append(dev[common.BookInfoOperation].Templates, "file://templates/bookinfo-gateway.yaml")
	dev[common.HTTPBinOperation].Templates = append(dev[common.HTTPBinOperation].Templates, "file://templates/httpbin-gateway.yaml")
	dev[common.ImageHubOperation].Templates = append(dev[common.ImageHubOperation].Templates, "file://templates/imagehub-gateway.yaml")
	dev[common.EmojiVotoOperation].Templates = append(dev[common.EmojiVotoOperation].Templates, "file://templates/emojivoto-gateway.yaml")

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

	dev[IstioVetOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_VALIDATE),
		Description: "Analyze Running Configuration",
	}

	dev[EnvoyFilterOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Envoy Filter for Imagehub",
		Versions:    adapter.NoneVersion,
		Templates: []adapter.Template{
			"file://templates/imagehub-filter.yaml",
		},
		AdditionalProperties: map[string]string{
			ServiceName:    "api-v1",
			EnvoyPatchFile: "file://templates/imagehub-patch.json",
		},
	}

	dev[DenyAllPolicyOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Deny-All Policy",
		Templates: []adapter.Template{
			"file://templates/policy-denyall.yaml",
		},
	}

	dev[StrictMTLSPolicyOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Strict MTLS Policy",
		Templates: []adapter.Template{
			"file://templates/policy-strict.yaml",
		},
	}

	dev[MutualMTLSPolicyOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Mutual MTLS Policy",
		Templates: []adapter.Template{
			"file://templates/policy-mutual.yaml",
		},
	}

	dev[DisableMTLSPolicyOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Disable MTLS Policy",
		Templates: []adapter.Template{
			"file://templates/policy-disable.yaml",
		},
	}

	return dev
}
