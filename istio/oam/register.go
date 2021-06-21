package oam

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-istio/internal/config"
)

var (
	basePath, _  = os.Getwd()
	workloadPath = filepath.Join(basePath, "templates", "oam", "workloads")
	traitPath    = filepath.Join(basePath, "templates", "oam", "traits")
)

// GenericStructure struct defines the body of the POST request that is sent to the OAM
// registry (Meshery)
//
// The body contains the
// 1. OAM definition, which is in accordance with the OAM spec
// 2. OAMRefSchema, which is json schema draft-4, draft-7 or draft-8 for the corresponding OAM object
// 3. Host is this service's grpc address in the form of `hostname:port`
type GenericStructure struct {
	OAMDefinition interface{} `json:"oam_definition,omitempty"`
	OAMRefSchema  string      `json:"oam_ref_schema,omitempty"`
	Host          string      `json:"host,omitempty"`
}

// RegisterWorkloads will register all of the workload definitions
// present in the path oam/workloads
//
// Registeration process will send POST request to $runtime/api/experimental/oam/workload
func RegisterWorkloads(runtime, host string) error {
	workloads := []string{
		"istiomesh",
		"grafanaistioaddon",
		"prometheusistioaddon",
		"zipkinistioaddon",
		"jaegeristioaddon",
		"virtualservice",
		"envoyfilter",
	}

	oamRDP := []adapter.OAMRegistrantDefinitionPath{}

	for _, workload := range workloads {
		defintionPath, schemaPath := generatePaths(workloadPath, workload)

		metadata := map[string]string{
			config.OAMAdapterNameMetadataKey: config.IstioOperation,
		}

		if strings.HasSuffix(workload, "addon") {
			metadata[config.OAMComponentCategoryMetadataKey] = "addon"
		}

		oamRDP = append(oamRDP, adapter.OAMRegistrantDefinitionPath{
			OAMDefintionPath: defintionPath,
			OAMRefSchemaPath: schemaPath,
			Host:             host,
			Metadata:         metadata,
		})
	}

	return adapter.
		NewOAMRegistrant(oamRDP, fmt.Sprintf("%s/api/experimental/oam/workload", runtime)).
		Register()
}

// RegisterTraits will register all of the trait definitions
// present in the path oam/traits
//
// Registeration process will send POST request to $runtime/api/experimental/oam/trait
func RegisterTraits(runtime, host string) error {
	traits := []string{
		"automaticsidecarinjection",
		"mtls",
	}

	oamRDP := []adapter.OAMRegistrantDefinitionPath{}

	for _, trait := range traits {
		defintionPath, schemaPath := generatePaths(traitPath, trait)

		oamRDP = append(oamRDP, adapter.OAMRegistrantDefinitionPath{
			OAMDefintionPath: defintionPath,
			OAMRefSchemaPath: schemaPath,
			Host:             host,
			Metadata: map[string]string{
				config.OAMAdapterNameMetadataKey: config.IstioOperation,
			},
		})
	}

	return adapter.
		NewOAMRegistrant(oamRDP, fmt.Sprintf("%s/api/experimental/oam/trait", runtime)).
		Register()
}

func generatePaths(path, name string) (defintionPath, schemaPath string) {
	definitionName := fmt.Sprintf("%s_definition.json", name)
	schemaName := fmt.Sprintf("%s.meshery.layer5.io.schema.json", name)

	return filepath.Join(path, definitionName), filepath.Join(path, schemaName)
}
