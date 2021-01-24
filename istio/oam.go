package istio

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshery-adapter-library/common"
	"github.com/layer5io/meshery-istio/internal/config"
	"github.com/layer5io/meshery-istio/istio/oam/core/v1alpha1"
	"github.com/layer5io/meshkit/errors"
)

// HandleComponents handles the processing of OAM components
func (istio *Istio) HandleComponents(comps []v1alpha1.Component, isDel bool) error {
	var errs []error
	for _, comp := range comps {
		if comp.Spec.Type == "IstioMesh" {
			if err := handleComponentIstioMesh(istio, comp, isDel); err != nil {
				errs = append(errs, err)
			}

			continue
		}

		if err := handleComponentIstioAddon(istio, comp, isDel); err != nil {
			errs = append(errs, err)
		}
	}

	if err := mergeErrors(errs); err != nil {
		return errors.NewDefault("", err.Error())
	}

	return nil
}

// HandleApplicationConfiguration handles the processing of OAM application configuration
func (istio *Istio) HandleApplicationConfiguration(config v1alpha1.Configuration, isDel bool) error {
	var errs []error
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
		}
	}

	return mergeErrors(errs)
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

func handleComponentIstioMesh(istio *Istio, comp v1alpha1.Component, isDel bool) error {
	// Get the istio version from the settings
	// we are sure that the version of istio would be present
	// because the configuration is already validated against the schema
	version := comp.Spec.Settings["version"].(string)

	_, err := istio.installIstio(isDel, version, comp.Namespace)

	return err
}

func handleComponentIstioAddon(istio *Istio, comp v1alpha1.Component, isDel bool) error {
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
	}

	// Get the service
	svc := config.Operations[addonName].AdditionalProperties[common.ServiceName]

	// Get the patch
	patch := config.Operations[addonName].AdditionalProperties[config.PatchFile]

	// Get the templates
	templates := config.Operations[addonName].Templates

	_, err := istio.installAddon(comp.Namespace, isDel, svc, patch, templates)

	return err
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
