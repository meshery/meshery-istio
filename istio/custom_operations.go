// Package istio - Common operations for the adapter
package istio

import (
	"github.com/layer5io/meshery-adapter-library/status"
)

func (istio *Istio) applyCustomOperation(namespace string, manifest string, isDel bool) (string, error) {
	st := status.Starting
	if istio.KubeClient == nil {
		return st, ErrNilClient
	}

	err := istio.applyManifest([]byte(manifest), isDel, namespace)
	if err != nil {
		return st, ErrCustomOperation(err)
	}

	return status.Completed, nil
}
