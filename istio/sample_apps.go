package istio

import (
	"context"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshkit/utils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
)

func (istio *Istio) installSampleApp(namespace string, del bool, templates []adapter.Template) (string, error) {
	st := status.Installing

	if del {
		st = status.Removing
	}

	for _, template := range templates {
		err := istio.applyManifest([]byte(template.String()), del, namespace)
		if err != nil {
			return st, ErrSampleApp(err)
		}
	}

	return status.Installed, nil
}

func (istio *Istio) patchWithEnvoyFilter(namespace string, del bool, app string, templates []adapter.Template, patchObject string) (string, error) {
	st := status.Deploying

	if del {
		st = status.Removing
	}

	jsonContents, err := utils.ReadFileSource(patchObject)
	if err != nil {
		return st, ErrEnvoyFilter(err)
	}

	_, err = istio.KubeClient.AppsV1().Deployments(namespace).Patch(context.TODO(), app, types.MergePatchType, []byte(jsonContents), metav1.PatchOptions{})
	if err != nil {
		return st, ErrEnvoyFilter(err)
	}

	for _, template := range templates {
		contents, err := utils.ReadFileSource(string(template))
		if err != nil {
			return st, ErrEnvoyFilter(err)
		}

		err = istio.applyManifest([]byte(contents), del, namespace)
		if err != nil {
			return st, ErrEnvoyFilter(err)
		}
	}

	return status.Deployed, nil
}
func (istio *Istio) applyPolicy(namespace string, del bool, templates []adapter.Template) (string, error) {
	st := status.Deploying

	if del {
		st = status.Removing
	}

	for _, template := range templates {
		contents, err := utils.ReadFileSource(string(template))
		if err != nil {
			return st, ErrApplyPolicy(err)
		}

		err = istio.applyManifest([]byte(contents), del, namespace)
		if err != nil {
			return st, ErrApplyPolicy(err)
		}
	}
	return status.Deployed, nil
}

// LoadToMesh is used to mark deployment for automatic sidecar injection (or not)
func (istio *Istio) LoadToMesh(namespace string, service string, remove bool) error {
	deploy, err := istio.KubeClient.AppsV1().Deployments(namespace).Get(context.TODO(), service, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if deploy.ObjectMeta.Labels == nil {
		deploy.ObjectMeta.Labels = map[string]string{}
	}
	deploy.ObjectMeta.Labels["istio-injection"] = "enabled"

	if remove {
		delete(deploy.ObjectMeta.Labels, "istio-injection")
	}

	_, err = istio.KubeClient.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// LoadNamespaceToMesh is used to mark namespaces for automatic sidecar injection (or not)
func (istio *Istio) LoadNamespaceToMesh(namespace string, remove bool) error {
	if istio.KubeClient == nil {
		return ErrNilClient
	}

	ns, err := istio.KubeClient.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		return ErrLoadNamespace(err, namespace)
	}

	if ns.ObjectMeta.Labels == nil {
		ns.ObjectMeta.Labels = map[string]string{}
	}
	ns.ObjectMeta.Labels["istio-injection"] = "enabled"

	if remove {
		delete(ns.ObjectMeta.Labels, "istio-injection")
	}

	_, err = istio.KubeClient.CoreV1().Namespaces().Update(context.TODO(), ns, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
