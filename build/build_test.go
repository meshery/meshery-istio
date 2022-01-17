package build

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/layer5io/meshery-adapter-library/adapter"
)

//TestCreateComponents can be used to directly invoke CreateComponents for writing new definitions and schemas onto the file system.
// This test is left here only for development purposes
func TestCreateComponents(t *testing.T) {
	wd, _ := os.Getwd()
	err := adapter.CreateComponents(adapter.StaticCompConfig{
		URL:     DefaultGenerationURL,
		Method:  DefaultGenerationMethod,
		Path:    filepath.Join(wd, "../templates", "oam", "workloads"),
		DirName: LatestVersion,
		Config:  NewConfig(LatestVersion),
		Force:   true,
	})
	if err != nil {
		t.Fatalf("Failed to generate components: %s", err.Error())
	}
}
