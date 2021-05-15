package config

import (
	"gopkg.in/yaml.v2"
)

// Structs generated from template/imagehub-filter.yaml
// EnvoyFilter
type EnvoyFilter struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

// VmConfig
type VmConfig struct {
	Code             Code                  `yaml:"code"`
	Configuration    VmConfigConfiguration `yaml:"configuration"`
	Runtime          string                `yaml:"runtime"`
	VmId             string                `yaml:"vmId"`
	AllowPrecompiled bool                  `yaml:"allow_precompiled"`
}

// VmConfigConfiguration
type VmConfigConfiguration struct {
	Type  string `yaml:"@type"`
	Value string `yaml:"value"`
}

// Labels
type Labels struct {
	App     string `yaml:"app"`
	Version string `yaml:"version"`
}

// Metadata
type Metadata struct {
	Name string `yaml:"name"`
}

// ConfigPatches
type ConfigPatches struct {
	ApplyTo string `yaml:"applyTo"`
	Match   Match  `yaml:"match"`
	Patch   Patch  `yaml:"patch"`
}

// Match
type Match struct {
	Context  string   `yaml:"context"`
	Proxy    Proxy    `yaml:"proxy"`
	Listener Listener `yaml:"listener"`
}

// Proxy
type Proxy struct {
	ProxyVersion string `yaml:"proxyVersion"`
}

// Value
type Value struct {
	Name        string      `yaml:"name"`
	TypedConfig TypedConfig `yaml:"typed_config"`
}

// Config
type FilterConfig struct {
	Configuration Configuration `yaml:"configuration"`
	RootId        string        `yaml:"root_id"`
	VmConfig      VmConfig      `yaml:"vmConfig"`
}

// Configuration
type Configuration struct {
	Type  string `yaml:"@type"`
	Value string `yaml:"value"`
}

// Local
type Local struct {
	Filename string `yaml:"filename"`
}

// Listener
type Listener struct {
	PortNumber  int         `yaml:"portNumber"`
	FilterChain FilterChain `yaml:"filterChain"`
}

// FilterChain
type FilterChain struct {
	Filter Filter `yaml:"filter"`
}

// Filter
type Filter struct {
	Name      string    `yaml:"name"`
	SubFilter SubFilter `yaml:"subFilter"`
}

// TypedConfigValue
type TypedConfigValue struct {
	Config FilterConfig `yaml:"config"`
}

// Code
type Code struct {
	Local Local `yaml:"local"`
}

// WorkloadSelector
type WorkloadSelector struct {
	Labels Labels `yaml:"labels"`
}

// Spec
type Spec struct {
	ConfigPatches    []ConfigPatches  `yaml:"configPatches"`
	WorkloadSelector WorkloadSelector `yaml:"workloadSelector"`
}

// SubFilter
type SubFilter struct {
	Name string `yaml:"name"`
}

// Patch
type Patch struct {
	Operation string `yaml:"operation"`
	Value     Value  `yaml:"value"`
}

// TypedConfig
type TypedConfig struct {
	Type    string           `yaml:"@type"`
	TypeUrl string           `yaml:"type_url"`
	Value   TypedConfigValue `yaml:"value"`
}

// GenerateImageHubTemplate() generates an EnvoyFilter config
// containing the json object to be given to imagehub's
// rate limit filter
func GenerateImagehubEnvoyFilter(encodedValue string) (string, error) {

	// generate the defaults
	ef := generateDefaults()

	if len(encodedValue) == 0 {
		return defaultConfigYAML(), nil
	}

	// this field will contain the base64 encoded json object for rate limit filter
	ef.Spec.ConfigPatches[0].Patch.Value.TypedConfig.Value.Config.VmConfig.Configuration.Value = encodedValue

	// Marshalling to yaml after making changes
	newYAML, err := yaml.Marshal(ef)
	if err != nil {
		return defaultConfigYAML(), err
	}

	return string(newYAML), nil
}

// returns the YAML as a string from an EnvoyFilter struct
func defaultConfigYAML() string {
	ef := generateDefaults()

	configYAML, err := yaml.Marshal(ef)
	if err != nil {
		return ""
	}

	return string(configYAML)
}

// generates the default EnvoyFilter config for the rate limit filter
func generateDefaults() EnvoyFilter {

	var configPatches []ConfigPatches

	defaultPatch := Patch{
		Operation: "INSERT_BEFORE",
		Value: Value{
			Name: "envoy.filter.http.wasm",
			TypedConfig: TypedConfig{
				Type:    "type.googleapis.com/udpa.type.v1.TypedStruct",
				TypeUrl: "type.googleapis.com/envoy.extensions.filters.http.wasm.v3.Wasm",
				Value: TypedConfigValue{
					Config: FilterConfig{
						Configuration: Configuration{
							Type:  "type.googleapis.com/google.protobuf.StringValue",
							Value: "rate_limit_filter",
						},
						RootId: "rate_limit_filter",
						VmConfig: VmConfig{
							Code: Code{
								Local: Local{
									Filename: "/var/lib/imagehub/filter.wasm",
								},
							},
							Configuration: VmConfigConfiguration{
								Type: "type.googleapis.com/google.protobuf.StringValue",
								// a default value to use in case no configuration is supplied
								Value: "WwogIHsKICAgICJuYW1lIjogIi9wdWxsIiwKICAgICJydWxlIjp7CiAgICAgICJydWxlVHlwZSI6ICJyYXRlLWxpbWl0ZXIiLAogICAgICAicGFyYW1ldGVycyI6WwogICAgICAgIHsiaWRlbnRpZmllciI6ICJFbnRlcnByaXNlIiwgImxpbWl0IjogMTAwMH0sCiAgICAgICAgeyJpZGVudGlmaWVyIjogIlRlYW0iLCAibGltaXQiOiAxMDB9LAogICAgICAgIHsiaWRlbnRpZmllciI6ICJQZXJzb25hbCIsICJsaW1pdCI6IDEwfQogICAgICBdCiAgICB9CiAgfSwKICB7CiAgICAibmFtZSI6ICIvYXV0aCIsCiAgICAicnVsZSI6ewogICAgICAicnVsZVR5cGUiOiAibm9uZSIKICAgIH0KICB9LAogIHsKICAgICJuYW1lIjogIi9zaWdudXAiLAogICAgInJ1bGUiOnsKICAgICAgInJ1bGVUeXBlIjogIm5vbmUiCiAgICB9CiAgfSwKICB7CiAgICAibmFtZSI6ICIvdXBncmFkZSIsCiAgICAicnVsZSI6ewogICAgICAicnVsZVR5cGUiOiAibm9uZSIKICAgIH0KICB9Cl0=",
							},
							Runtime:          "envoy.wasm.runtime.v8",
							VmId:             "rate_limit_filter",
							AllowPrecompiled: true,
						},
					},
				},
			},
		},
	}
	cfgPatch0 := ConfigPatches{
		ApplyTo: "HTTP_FILTER",
		Match: Match{
			Context: "SIDECAR_INBOUND",
			Proxy: Proxy{
				ProxyVersion: "^1\\.9.*",
			},
			Listener: Listener{
				PortNumber: 9091,
				FilterChain: FilterChain{
					Filter: Filter{
						Name: "envoy.http_connection_manager",
						SubFilter: SubFilter{
							Name: "envoy.router",
						},
					},
				},
			},
		},
		Patch: defaultPatch,
	}

	ef := EnvoyFilter{
		ApiVersion: "networking.istio.io/v1alpha3",
		Kind:       "EnvoyFilter",
		Metadata: Metadata{
			Name: "imagehub-filter",
		},
		Spec: Spec{
			ConfigPatches: append(configPatches, cfgPatch0),
			WorkloadSelector: WorkloadSelector{
				Labels: Labels{
					App:     "api",
					Version: "v1",
				},
			},
		},
	}

	return ef
}
