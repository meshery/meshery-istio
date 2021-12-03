package istio

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshkit/utils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// installAddon installs/uninstalls an addon in the given namespace
//
// the template defines the manifest's link/location which needs to be used to
// install the addon
func (istio *Istio) installAddon(namespace string, del bool, service string, patches []string, templates []adapter.Template) (string, error) {
	st := status.Installing

	if del {
		st = status.Removing
	}

	istio.Log.Debug(fmt.Sprintf("Overidden namespace: %s", namespace))
	namespace = "istio-system"

	for _, template := range templates {
		if istio.KubeClient == nil {
			return st, ErrNilClient
		}
		err := istio.applyManifest([]byte(template.String()), del, namespace)
		// Specifically choosing to ignore kiali dashboard's error.
		// Referring to: https://github.com/kiali/kiali/issues/3112
		if err != nil && !strings.Contains(err.Error(), "no matches for kind \"MonitoringDashboard\" in version \"monitoring.kiali.io/v1alpha1\"") {
			if !strings.Contains(err.Error(), "clusterIP") {
				return st, ErrAddonFromTemplate(err)
			}
		}
	}

	for _, patch := range patches {
		if patch == "" {
			continue //avoid throwing error when a given patch key didn't exist for a specific addon type in operations
		}
		if !del {
			_, err := url.ParseRequestURI(patch)
			if err != nil {
				return st, ErrAddonFromTemplate(err)
			}

			content, err := utils.ReadFileSource(patch)
			if err != nil {
				return st, ErrAddonFromTemplate(err)
			}

			_, err = istio.KubeClient.CoreV1().Services(namespace).Patch(context.TODO(), service, types.MergePatchType, []byte(content), metav1.PatchOptions{})
			if err != nil {
				return st, ErrAddonFromTemplate(err)
			}
		}
	}

	return status.Installed, nil
}
