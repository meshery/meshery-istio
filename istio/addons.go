package istio

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshery-istio/internal/config"
)

// AddonTemplate is as a container for addon templates
var AddonTemplate = map[string]adapter.Template{
	config.PrometheusAddon: "https://raw.githubusercontent.com/istio/istio/master/samples/addons/prometheus.yaml",
	config.GrafanaAddon:    "https://raw.githubusercontent.com/istio/istio/master/samples/addons/grafana.yaml",
	config.KialiAddon:      "https://raw.githubusercontent.com/istio/istio/master/samples/addons/kiali.yaml",
	config.JaegerAddon:     "https://raw.githubusercontent.com/istio/istio/master/samples/addons/jaeger.yaml",
	config.ZipkinAddon:     "https://raw.githubusercontent.com/istio/istio/master/samples/addons/extras/zipkin.yaml",
}

// InstallAddon installs the specified addon in the given namespace
func (istio *Istio) InstallAddon(namespace string, del bool, addon string) (string, error) {
	// Some of addons will have a different install process
	// than these template based addons
	switch addon {
	case config.PrometheusAddon, config.GrafanaAddon, config.KialiAddon, config.JaegerAddon, config.ZipkinAddon:
		return istio.installAddonFromTemplate(namespace, del, AddonTemplate[addon])
	}

	return "", ErrAddonInvalidConfig(fmt.Errorf("%s is invalid addon", addon))
}

// installAddonFromTemplate installs/uninstalls an addon in the given namespace
//
// the template defines the manifest's link/location which needs to be used to
// install the addon
func (istio *Istio) installAddonFromTemplate(namespace string, del bool, template adapter.Template) (string, error) {
	istio.Log.Info(fmt.Sprintf("Requested action is delete: %v", del))
	st := status.Installing

	if del {
		st = status.Removing
	}

	contents, err := readFileSource(string(template))
	if err != nil {
		return st, ErrAddonFromTemplate(err)
	}

	err = istio.applyManifest([]byte(contents), del, namespace)
	// Specifically choosing to ignore kiali dashboard's error.
	// Referring to: https://github.com/kiali/kiali/issues/3112
	if err != nil && !strings.Contains(err.Error(), "no matches for kind \"MonitoringDashboard\" in version \"monitoring.kiali.io/v1alpha1\"") {
		return st, ErrAddonFromTemplate(err)
	}

	return status.Installed, nil
}
