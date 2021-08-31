package oam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-istio/internal/config"
	"github.com/layer5io/meshkit/utils/manifests"
)

var (
	basePath, _  = os.Getwd()
	workloadPath = filepath.Join(basePath, "templates", "oam", "workloads")
	traitPath    = filepath.Join(basePath, "templates", "oam", "traits")
)

type schemaDefinitionPathSet struct {
	oamDefinitionPath string
	jsonSchemaPath    string
	name              string
}

// RegisterWorkloads will register all of the workload definitions
// present in the path oam/workloads
//
// Registration process will send POST request to $runtime/api/oam/workload
func RegisterWorkloads(runtime, host string) error {
	oamRDP := []adapter.OAMRegistrantDefinitionPath{}

	pathSets, err := load(workloadPath)
	if err != nil {
		return err
	}

	for _, pathSet := range pathSets {
		metadata := map[string]string{
			config.OAMAdapterNameMetadataKey: config.IstioOperation,
		}

		if strings.HasSuffix(pathSet.name, "addon") {
			metadata[config.OAMComponentCategoryMetadataKey] = "addon"
		}

		oamRDP = append(oamRDP, adapter.OAMRegistrantDefinitionPath{
			OAMDefintionPath: pathSet.oamDefinitionPath,
			OAMRefSchemaPath: pathSet.jsonSchemaPath,
			Host:             host,
			Metadata:         metadata,
		})
	}

	return adapter.
		NewOAMRegistrant(oamRDP, fmt.Sprintf("%s/api/oam/workload", runtime)).
		Register()
}

// RegisterTraits will register all of the trait definitions
// present in the path oam/traits
//
// Registeration process will send POST request to $runtime/api/oam/trait
func RegisterTraits(runtime, host string) error {
	oamRDP := []adapter.OAMRegistrantDefinitionPath{}

	pathSets, err := load(traitPath)
	if err != nil {
		return err
	}

	for _, pathSet := range pathSets {
		metadata := map[string]string{
			config.OAMAdapterNameMetadataKey: config.IstioOperation,
		}

		oamRDP = append(oamRDP, adapter.OAMRegistrantDefinitionPath{
			OAMDefintionPath: pathSet.oamDefinitionPath,
			OAMRefSchemaPath: pathSet.jsonSchemaPath,
			Host:             host,
			Metadata:         metadata,
		})
	}

	return adapter.
		NewOAMRegistrant(oamRDP, fmt.Sprintf("%s/api/oam/trait", runtime)).
		Register()
}

func load(basePath string) ([]schemaDefinitionPathSet, error) {
	res := []schemaDefinitionPathSet{}

	if err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if matched, err := filepath.Match("*_definition.json", filepath.Base(path)); err != nil {
			return err
		} else if matched {
			nameWithPath := strings.TrimSuffix(path, "_definition.json")

			res = append(res, schemaDefinitionPathSet{
				oamDefinitionPath: path,
				jsonSchemaPath:    fmt.Sprintf("%s.meshery.layer5io.schema.json", nameWithPath),
				name:              filepath.Base(nameWithPath),
			})
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return res, nil
}

//RegisterWorkLoadsDynamically ...
func RegisterWorkLoadsDynamically(runtime, host string) error {
	release, err := config.GetLatestReleases(1)
	if err != nil {
		fmt.Println("Could not get latest stable release")
		return err
	}
	v := release[0].TagName
	m := manifests.Config{
		Name:        "Istio",
		MeshVersion: v,
		Filter: manifests.CrdFilter{
			RootFilter:    []string{"$[?(@.kind==\"CustomResourceDefinition\")]"},
			NameFilter:    []string{"$..[\"spec\"][\"names\"][\"kind\"]"},
			VersionFilter: []string{"$..spec.versions[0]", " --o-filter", "$[0]"},
			GroupFilter:   []string{"$..spec", " --o-filter", "$[]"},
			SpecFilter:    []string{"$..openAPIV3Schema.properties.spec", " --o-filter", "$[]"},
		},
	}
	fmt.Println("version ", v)
	url := "https://raw.githubusercontent.com/istio/istio/" + v + "/manifests/charts/base/crds/crd-all.gen.yaml"
	comp, err := manifests.GetFromManifest(url, manifests.SERVICE_MESH, m)
	if err != nil {
		return err
	}
	for i, def := range comp.Definitions {
		var ord adapter.OAMRegistrantData
		ord.OAMRefSchema = comp.Schemas[i]

		//Marshalling the stringified json
		ord.Host = host
		definitionMap := map[string]interface{}{}
		if err := json.Unmarshal([]byte(def), &definitionMap); err != nil {
			return err
		}
		// To be shifted in meshkit
		definitionMap["apiVersion"] = "core.oam.dev/v1alpha1"
		definitionMap["kind"] = "WorkloadDefinition"
		ord.OAMDefinition = definitionMap
		ord.Metadata = map[string]string{
			config.OAMAdapterNameMetadataKey: config.IstioOperation,
		}
		// send request to the register
		backoffOpt := backoff.NewExponentialBackOff()
		backoffOpt.MaxElapsedTime = 10 * time.Minute
		if err := backoff.Retry(func() error {
			contentByt, err := json.Marshal(ord)
			if err != nil {
				return backoff.Permanent(err)
			}
			content := bytes.NewReader(contentByt)
			// host here is given by the application itself and is trustworthy hence,
			// #nosec
			resp, err := http.Post(fmt.Sprintf("%s/api/oam/workload", runtime), "application/json", content)
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
		}, backoffOpt); err != nil {
			return err
		}
	}
	return nil
}
