package build

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/layer5io/meshery-adapter-library/adapter"

	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/manifests"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var DefaultGenerationMethod string
var DefaultGenerationURL string
var LatestVersion string
var WorkloadPath string

//in library
type StaticCompConfig struct {
	URL     string //URL
	Method  string //Use the constants exported by package. Manifests or Helm
	Path    string //Where to store the directory.(Each directory will have an array of definitions and schemas)
	DirName string //The directory's name. By convention, it should be the version name
	Config  manifests.Config
	Force   bool //When set to true, if the file with same name already exists, they will be overriden
}

//in library
func CreateComponents(scfg StaticCompConfig) error {
	dir := filepath.Join(scfg.Path, scfg.DirName)
	_, err := os.Stat(dir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(dir, 0777)
		if err != nil {
			return err
		}
	}
	var comp *manifests.Component
	switch scfg.Method {
	case adapter.Manifests:
		comp, err = manifests.GetFromManifest(scfg.URL, manifests.SERVICE_MESH, scfg.Config)
	case adapter.HelmCHARTS:
		comp, err = manifests.GetFromHelm(scfg.URL, manifests.SERVICE_MESH, scfg.Config)
	default:
		return err
	}
	if comp == nil {
		return errors.New("nil components")
	}

	for i, def := range comp.Definitions {
		schema := comp.Schemas[i]
		name := GetNameFromWorkloadDefinition([]byte(def))
		defFileName := name + "_definition.json"
		schemaFileName := name + ".meshery.layer5io.schema.json"
		err := writeToFile(filepath.Join(dir, defFileName), []byte(def), scfg.Force)
		if err != nil {
			return err
		}
		err = writeToFile(filepath.Join(dir, schemaFileName), []byte(schema), scfg.Force)
		if err != nil {
			return err
		}
	}
	return nil
}

//create a file with this filename and stuff the string
func writeToFile(path string, data []byte, force bool) error {
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) { //There some other error than non existence of file
		return err
	}

	if err == nil { //file already exists
		if !force { // Dont override existing file, skip it
			fmt.Println("File already exists,skipping...")
			return nil
		}
		err := os.Remove(path) //Remove the existing file, before overriding it
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(path, data, 0777)
}

func GetNameFromWorkloadDefinition(definition []byte) string {
	var wd v1alpha1.WorkloadDefinition
	err := json.Unmarshal(definition, &wd)
	if err != nil {
		return ""
	}
	return wd.Spec.DefinitionRef.Name
}

//Should stay here
func NewConfig(version string) manifests.Config {
	return manifests.Config{
		Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_ISTIO)],
		MeshVersion: version,
		Filter: manifests.CrdFilter{
			RootFilter:    []string{"$[?(@.kind==\"CustomResourceDefinition\")]"},
			NameFilter:    []string{"$..[\"spec\"][\"names\"][\"kind\"]"},
			VersionFilter: []string{"$[0]..spec.versions[0]"},
			GroupFilter:   []string{"$[0]..spec"},
			SpecFilter:    []string{"$[0]..openAPIV3Schema.properties.spec"},
			ItrFilter:     []string{"$[?(@.spec.names.kind"},
			ItrSpecFilter: []string{"$[?(@.spec.names.kind"},
			VField:        "name",
			GField:        "group",
		},
	}
}
func init() {
	wd, _ := os.Getwd()
	WorkloadPath = filepath.Join(wd, "templates", "oam", "workloads")
	versions, _ := utils.GetLatestReleaseTagsSorted("istio", "istio")
	if len(versions) == 0 {
		return //Since
	}
	LatestVersion = versions[len(versions)-1]
	DefaultGenerationMethod = adapter.Manifests
	DefaultGenerationURL = "https://raw.githubusercontent.com/istio/istio/" + LatestVersion + "/manifests/charts/base/crds/crd-all.gen.yaml"
}
