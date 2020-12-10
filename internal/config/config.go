package config

import (
	"path"

	"github.com/layer5io/meshery-adapter-library/common"
	"github.com/layer5io/meshery-adapter-library/config"
	configprovider "github.com/layer5io/meshery-adapter-library/config/provider"
	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshkit/utils"
)

const (
	// Constants to use in log statements
	IstioOperation = "istio"
	LabelNamespace = "label-namespace"
	PatchFile      = "patch-file"

	// Istio vet operation
	IstioVetOperation = "istio-vet"

	// Configure Envoy filter operation
	EnvoyFilterOperation = "envoy-filter-operation"

	// Addons that the adapter supports
	PrometheusAddon = "prometheus-addon"
	GrafanaAddon    = "grafana-addon"
	KialiAddon      = "kiali-addon"
	JaegerAddon     = "jaeger-addon"
	ZipkinAddon     = "zipkin-addon"

	// Policies
	DenyAllPolicyOperation     = "deny-all-policy-operation"
	StrictMTLSPolicyOperation  = "strict-mtls-policy-operation"
	MutualMTLSPolicyOperation  = "mutual-mtls-policy-operation"
	DisableMTLSPolicyOperation = "disable-mtls-policy-operation"
)

var (
	configRootPath = path.Join(utils.GetHome(), ".meshery")

	Config = configprovider.Options{
		ServerConfig:   ServerConfig,
		MeshSpec:       MeshSpec,
		ProviderConfig: ProviderConfig,
		Operations:     Operations,
	}

	ServerConfig = map[string]string{
		"name":    "istio-adapter",
		"port":    "10000",
		"version": "v1.0.0",
	}

	MeshSpec = map[string]string{
		"name":     "istio",
		"status":   status.None,
		"traceurl": status.None,
		"version":  status.None,
	}

	ProviderConfig = map[string]string{
		configprovider.FilePath: configRootPath,
		configprovider.FileType: "yaml",
		configprovider.FileName: "istio",
	}

	// KubeConfig - Controlling the kubeconfig lifecycle with viper
	KubeConfig = map[string]string{
		configprovider.FilePath: configRootPath,
		configprovider.FileType: "yaml",
		configprovider.FileName: "kubeconfig",
	}

	Operations = getOperations(common.Operations)
)

// New creates a new config instance
func New(provider string) (config.Handler, error) {
	// Config provider
	switch provider {
	case configprovider.ViperKey:
		return configprovider.NewViper(Config)
	case configprovider.InMemKey:
		return configprovider.NewInMem(Config)
	}

	return nil, ErrEmptyConfig
}

func NewKubeconfigBuilder(provider string) (config.Handler, error) {
	opts := configprovider.Options{}
	opts.ProviderConfig = KubeConfig

	// Config provider
	switch provider {
	case configprovider.ViperKey:
		return configprovider.NewViper(opts)
	case configprovider.InMemKey:
		return configprovider.NewInMem(opts)
	}
	return nil, ErrEmptyConfig
}

// RootPath returns the config root path for the adapter
func RootPath() string {
	return configRootPath
}
