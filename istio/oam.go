package istio

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshery-adapter-library/common"
	"github.com/layer5io/meshery-istio/internal/config"
	"github.com/layer5io/meshkit/errors"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	"gopkg.in/yaml.v2"
)

// CompHandler is the type for functions which can handle OAM components
type CompHandler func(*Istio, v1alpha1.Component, bool) (string, error)

// HandleComponents handles the processing of OAM components
func (istio *Istio) HandleComponents(comps []v1alpha1.Component, isDel bool) (string, error) {
	var errs []error
	var msgs []string

	compFuncMap := map[string]CompHandler{
		"IstioMesh":            handleComponentIstioMesh,
		"VirtualService":       handleComponentVirtualService,
		"EnvoyFilterIstio":     handleComponentEnvoyFilter,
		"GrafanaIstioAddon":    handleComponentIstioAddon,
		"PrometheusIstioAddon": handleComponentIstioAddon,
		"ZipkinIstioAddon":     handleComponentIstioAddon,
		"JaegerIstioAddon":     handleComponentIstioAddon,
	}

	for _, comp := range comps {
		fnc, ok := compFuncMap[comp.Spec.Type]
		if !ok {
			return "", ErrInvalidOAMComponentType(comp.Spec.Type)
		}

		msg, err := fnc(istio, comp, isDel)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		msgs = append(msgs, msg)
	}

	if err := mergeErrors(errs); err != nil {
		return mergeMsgs(msgs), errors.NewDefault("", err.Error())
	}

	return mergeMsgs(msgs), nil
}

// HandleApplicationConfiguration handles the processing of OAM application configuration
func (istio *Istio) HandleApplicationConfiguration(config v1alpha1.Configuration, isDel bool) (string, error) {
	var errs []error
	var msgs []string
	for _, comp := range config.Spec.Components {
		for _, trait := range comp.Traits {
			if trait.Name == "mTLS" {
				namespaces := castSliceInterfaceToSliceString(trait.Properties["namespaces"].([]interface{}))
				policy := trait.Properties["policy"].(string)

				if err := handleMTLS(istio, namespaces, policy, isDel); err != nil {
					errs = append(errs, err)
				}
			}

			if trait.Name == "automaticsidecarinjection" {
				namespaces := castSliceInterfaceToSliceString(trait.Properties["namespaces"].([]interface{}))
				if err := handleNamespaceLabel(istio, namespaces, isDel); err != nil {
					errs = append(errs, err)
				}
			}

			msgs = append(msgs, fmt.Sprintf("applied trait \"%s\" on service \"%s\"", trait.Name, comp.ComponentName))
		}
	}

	if err := mergeErrors(errs); err != nil {
		return mergeMsgs(msgs), errors.NewDefault("", err.Error())
	}

	return mergeMsgs(msgs), nil
}

func handleMTLS(istio *Istio, namespaces []string, policy string, isDel bool) error {
	var errs []error
	for _, ns := range namespaces {
		policyName := fmt.Sprintf("%s-mtls-policy-operation", policy)

		if _, err := istio.applyPolicy(ns, isDel, config.Operations[policyName].Templates); err != nil {
			errs = append(errs, err)
		}
	}

	return mergeErrors(errs)
}

func handleNamespaceLabel(istio *Istio, namespaces []string, isDel bool) error {
	var errs []error
	for _, ns := range namespaces {
		if err := istio.LoadNamespaceToMesh(ns, isDel); err != nil {
			errs = append(errs, err)
		}
	}

	return mergeErrors(errs)
}

func handleComponentIstioMesh(istio *Istio, comp v1alpha1.Component, isDel bool) (string, error) {
	// Get the istio version from the settings
	// we are sure that the version of istio would be present
	// because the configuration is already validated against the schema
	version := comp.Spec.Settings["version"].(string)

	return istio.installIstio(isDel, version, comp.Namespace)
}

func handleComponentVirtualService(istio *Istio, comp v1alpha1.Component, isDel bool) (string, error) {
	return handleIstioCoreComponent(istio, comp, isDel, "networking.istio.io/v1beta1", "VirtualService")
}

func handleComponentEnvoyFilter(istio *Istio, comp v1alpha1.Component, isDel bool) (string, error) {
	return handleIstioCoreComponent(istio, comp, isDel, "networking.istio.io/v1alpha3", "EnvoyFilter")
}

func handleIstioCoreComponent(
	istio *Istio,
	comp v1alpha1.Component,
	isDel bool,
	apiVersion,
	kind string) (string, error) {
	component := map[string]interface{}{
		"apiVersion": apiVersion,
		"kind":       kind,
		"metadata": map[string]interface{}{
			"name":        comp.Name,
			"annotations": comp.Annotations,
			"labels":      comp.Labels,
		},
		"spec": comp.Spec.Settings,
	}

	// Convert to yaml
	yamlByt, err := yaml.Marshal(component)
	if err != nil {
		err = ErrParseIstioCoreComponent(err)
		istio.Log.Error(err)
		return "", err
	}

	msg := fmt.Sprintf("created %s \"%s\" in namespace \"%s\"", kind, comp.Name, comp.Namespace)
	if isDel {
		msg = fmt.Sprintf("deleted %s config \"%s\" in namespace \"%s\"", kind, comp.Name, comp.Namespace)
	}

	return msg, istio.applyManifest(yamlByt, isDel, comp.Namespace)
}

func handleComponentIstioAddon(istio *Istio, comp v1alpha1.Component, isDel bool) (string, error) {
	var addonName string

	switch comp.Spec.Type {
	case "GrafanaIstioAddon":
		addonName = config.GrafanaAddon
	case "PrometheusIstioAddon":
		addonName = config.PrometheusAddon
	case "ZipkinIstioAddon":
		addonName = config.ZipkinAddon
	case "JaegerIstioAddon":
		addonName = config.JaegerAddon
	default:
		return "", nil
	}

	// Get the service
	svc := config.Operations[addonName].AdditionalProperties[common.ServiceName]

	// Get the patches
	patches := make([]string, 0)
	patches = append(patches, config.Operations[addonName].AdditionalProperties[config.ServicePatchFile])
	patches = append(patches, config.Operations[addonName].AdditionalProperties[config.CPPatchFile])
	patches = append(patches, config.Operations[addonName].AdditionalProperties[config.ControlPatchFile])

	// Get the templates
	templates := config.Operations[addonName].Templates

	_, err := istio.installAddon(comp.Namespace, isDel, svc, patches, templates)

	msg := fmt.Sprintf("created service of type \"%s\"", comp.Spec.Type)
	if isDel {
		msg = fmt.Sprintf("deleted service of type \"%s\"", comp.Spec.Type)
	}

	return msg, err
}

func castSliceInterfaceToSliceString(in []interface{}) []string {
	var out []string

	for _, v := range in {
		cast, ok := v.(string)
		if ok {
			out = append(out, cast)
		}
	}

	return out
}

func mergeErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	var errMsgs []string

	for _, err := range errs {
		errMsgs = append(errMsgs, err.Error())
	}

	return fmt.Errorf(strings.Join(errMsgs, "\n"))
}

func mergeMsgs(strs []string) string {
	return strings.Join(strs, "\n")
}
