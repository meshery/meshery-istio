package oam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var (
	basePath, _  = os.Getwd()
	workloadPath = filepath.Join(basePath, "oam", "workloads")
	traitPath    = filepath.Join(basePath, "oam", "traits")
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
	}

	for _, workload := range workloads {
		oamcontent, err := readDefintionAndSchema(workloadPath, workload)
		if err != nil {
			return err
		}
		oamcontent.Host = host

		// Convert struct to json
		byt, err := json.Marshal(oamcontent)
		if err != nil {
			return err
		}

		reader := bytes.NewReader(byt)

		if err := register(fmt.Sprintf("%s/api/experimental/oam/%s", runtime, "workload"), reader); err != nil {
			return err
		}
	}

	return nil
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

	for _, trait := range traits {
		oamcontent, err := readDefintionAndSchema(traitPath, trait)
		if err != nil {
			return err
		}
		oamcontent.Host = host

		// Convert struct to json
		byt, err := json.Marshal(oamcontent)
		if err != nil {
			return err
		}

		reader := bytes.NewReader(byt)

		if err := register(fmt.Sprintf("%s/api/experimental/oam/%s", runtime, "trait"), reader); err != nil {
			return err
		}
	}

	return nil
}

func readDefintionAndSchema(path, name string) (*GenericStructure, error) {
	definitionName := fmt.Sprintf("%s_definition.json", name)
	schemaName := fmt.Sprintf("%s.meshery.layer5.io.schema.json", name)

	// Paths are constructed on the fly but are trusted hence,
	// #nosec
	definition, err := ioutil.ReadFile(filepath.Join(path, definitionName))
	if err != nil {
		return nil, err
	}

	var definitionMap map[string]interface{}
	if err := json.Unmarshal(definition, &definitionMap); err != nil {
		return nil, err
	}

	// Paths are constructed on the fly but are trusted hence,
	// #nosec
	schema, err := ioutil.ReadFile(filepath.Join(path, schemaName))
	if err != nil {
		return nil, err
	}

	return &GenericStructure{OAMDefinition: definitionMap, OAMRefSchema: string(schema)}, nil
}

func register(host string, content io.Reader) error {
	// host here is given by the application itself and is trustworthy hence,
	// #nosec
	resp, err := http.Post(host, "application/json", content)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf(
			"register process failed, host returned status: %s with status code %d",
			resp.Status,
			resp.StatusCode,
		)
	}

	return nil
}
