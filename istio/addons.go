package istio

import (
	"context"
	"fmt"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// AddonTemplate is as a container for addon templates
// var AddonTemplate = map[string]adapter.Template{
// 	config.PrometheusAddon: "https://raw.githubusercontent.com/istio/istio/master/samples/addons/prometheus.yaml",
// 	config.GrafanaAddon:    "https://raw.githubusercontent.com/istio/istio/master/samples/addons/grafana.yaml",
// 	config.KialiAddon:      "https://raw.githubusercontent.com/istio/istio/master/samples/addons/kiali.yaml",
// 	config.JaegerAddon:     "https://raw.githubusercontent.com/istio/istio/master/samples/addons/jaeger.yaml",
// 	config.ZipkinAddon:     "https://raw.githubusercontent.com/istio/istio/master/samples/addons/extras/zipkin.yaml",
// }

// // InstallAddon installs the specified addon in the given namespace
// func (istio *Istio) InstallAddon(namespace string, del bool, addon string) (string, error) {
// 	// Some of addons will have a different install process
// 	// than these template based addons
// 	switch addon {
// 	case config.PrometheusAddon, config.GrafanaAddon, config.KialiAddon, config.JaegerAddon, config.ZipkinAddon:
// 		return istio.installAddonFromTemplate(namespace, del, AddonTemplate[addon])
// 	}

// 	return "", ErrAddonInvalidConfig(fmt.Errorf("%s is invalid addon", addon))
// }

// installAddonFromTemplate installs/uninstalls an addon in the given namespace
//
// the template defines the manifest's link/location which needs to be used to
// install the addon
func (istio *Istio) installAddon(namespace string, del bool, service string, patch string, templates []adapter.Template) (string, error) {
	st := status.Installing

	if del {
		st = status.Removing
	}

	istio.Log.Debug(fmt.Sprintf("Overidden namespace: %s", namespace))
	namespace = ""

	for _, template := range templates {
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
	}

	jsonContents, err := readFileSource(patch)
	if err != nil {
		return st, ErrAddonFromTemplate(err)
	}

	_, err = istio.KubeClient.CoreV1().Services("istio-system").Patch(context.TODO(), service, types.MergePatchType, []byte(jsonContents), metav1.PatchOptions{})
	if err != nil {
		return st, ErrAddonFromTemplate(err)
	}

	return status.Installed, nil
}
