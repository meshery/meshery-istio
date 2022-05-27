package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"

	"cuelang.org/go/cue"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/manifests"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var DefaultGenerationMethod string
var DefaultGenerationURL string
var LatestVersion string
var WorkloadPath string
var AllVersions []string

const Component = "Istio"

//NewConfig creates the configuration for creating components
func NewConfig(version string) manifests.Config {
	return manifests.Config{
		Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_ISTIO)],
		Type:        Component,
		MeshVersion: version,
		CrdFilter: manifests.CueCrdFilter{
			IdentifierExtractor: func(rootCRDCueVal cue.Value) (cue.Value, error) {
				res := rootCRDCueVal.LookupPath(cue.ParsePath("spec.names.kind"))
				if !res.Exists() {
					return res, fmt.Errorf("Could not find the value")
				}
				return res.Value(), nil
			},
			NameExtractor: func(rootCRDCueVal cue.Value) (cue.Value, error) {
				res := rootCRDCueVal.LookupPath(cue.ParsePath("spec.names.kind"))
				if !res.Exists() {
					return res, fmt.Errorf("Could not find the value")
				}
				return res.Value(), nil
			},
			VersionExtractor: func(rootCRDCueVal cue.Value) (cue.Value, error) {
				res := rootCRDCueVal.LookupPath(cue.ParsePath("spec.versions[0].name"))
				if !res.Exists() {
					return res, fmt.Errorf("Could not find the value")
				}
				return res.Value(), nil
			},
			GroupExtractor: func(rootCRDCueVal cue.Value) (cue.Value, error) {
				res := rootCRDCueVal.LookupPath(cue.ParsePath("spec.group"))
				if !res.Exists() {
					return res, fmt.Errorf("Could not find the value")
				}
				return res.Value(), nil
			},
			SpecExtractor: func(rootCRDCueVal cue.Value) (cue.Value, error) {
				res := rootCRDCueVal.LookupPath(cue.ParsePath("spec.versions[0].schema.openAPIV3Schema.properties.spec"))
				if !res.Exists() {
					return res, fmt.Errorf("Could not find the value")
				}
				return res.Value(), nil
			},
		},
		ExtractCrds: func(manifest string) []string {
			crds := strings.Split(manifest, "---")
			// trim the spaces
			for _, crd := range crds {
				crd = strings.TrimSpace(crd)
			}
			return crds
		},
	}
}
func init() {
	wd, _ := os.Getwd()
	WorkloadPath = filepath.Join(wd, "templates", "oam", "workloads")
	AllVersions, _ = utils.GetLatestReleaseTagsSorted("istio", "istio")
	if len(AllVersions) == 0 {
		return
	}
	LatestVersion = AllVersions[len(AllVersions)-1]
	DefaultGenerationMethod = adapter.Manifests
	DefaultGenerationURL = "https://raw.githubusercontent.com/istio/istio/" + LatestVersion + "/manifests/charts/base/crds/crd-all.gen.yaml"
}
