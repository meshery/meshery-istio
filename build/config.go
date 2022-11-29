package build

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"

	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/manifests"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var DefaultGenerationMethod string
var DefaultGenerationURL string
var LatestVersion string
var WorkloadPath string
var MeshModelPath string
var AllVersions []string

const Component = "Istio"

var meshmodelmetadata = map[string]interface{}{
	"Primary Color":   "#466BB0",
	"Secondary Color": "#93b0e6",
	"Shape":           "circle",
	"Logo URL":        "https://github.com/istio/istio/blob/master/logo/istio-bluelogo-whitebackground-unframed.svg",
	"SVG_Color":       "<svg xmlns=\"http://www.w3.org/2000/svg\" version=\"1.1\" viewBox=\"0 0 160 240\"><g id=\"logo\" fill=\"#466BB0\"><rect id=\"background\" width=\"160\" height=\"240\" fill=\"#fff\"/><polygon id=\"hull\" points=\"0 210 160 210 60 240\"/><polygon id=\"mainsail\" points=\"0 200 60 190 60 80\"/><polygon id=\"headsail\" points=\"70 190 160 200 70 0\"/></g></svg>\n",
	"SVG_White":       "<svg viewBox=\"6.386270046234131 3.7419400215148926 18.227430820465088 25.258059978485107\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M6.38627 24.7904L13.2215 23.738V12.1613L6.38627 24.7904Z\" fill=\"white\"/><path d=\"M6.38627 25.8427H24.6137L13.2215 29L6.38627 25.8427Z\" fill=\"white\"/><path d=\"M14.3608 23.7379L24.6137 24.7904L14.3608 3.74194V23.7379Z\" fill=\"white\"/></svg>",
}
var MeshModelConfig = adapter.MeshModelConfig{ //Move to build/config.go
	Category:    "Orchestration & Management",
	SubCategory: "Service Mesh",
	Metadata:    meshmodelmetadata,
}

// NewConfig creates the configuration for creating components
func NewConfig(version string) manifests.Config {
	return manifests.Config{
		Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_ISTIO)],
		Type:        Component,
		MeshVersion: version,
		CrdFilter: manifests.NewCueCrdFilter(manifests.ExtractorPaths{
			NamePath:    "spec.names.kind",
			IdPath:      "spec.names.kind",
			VersionPath: "spec.versions[0].name",
			GroupPath:   "spec.group",
			SpecPath:    "spec.versions[0].schema.openAPIV3Schema.properties.spec"}, false),
		ExtractCrds: func(manifest string) []string {
			crds := strings.Split(manifest, "---")
			return crds
		},
	}
}
func init() {
	wd, _ := os.Getwd()
	WorkloadPath = filepath.Join(wd, "templates", "oam", "workloads")
	MeshModelPath = filepath.Join(wd, "templates", "meshmodel", "components")
	AllVersions, _ = utils.GetLatestReleaseTagsSorted("istio", "istio")
	if len(AllVersions) == 0 {
		return
	}
	LatestVersion = AllVersions[len(AllVersions)-1]
	DefaultGenerationMethod = adapter.Manifests
	DefaultGenerationURL = "https://raw.githubusercontent.com/istio/istio/" + LatestVersion + "/manifests/charts/base/crds/crd-all.gen.yaml"
}
