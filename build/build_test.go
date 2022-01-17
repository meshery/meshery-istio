package build

import "testing"

//TestCreateComponents can be used to directly invoke CreateComponents for writing new definitions and schemas onto the file system.
// This test is left here only for development purposes
func TestCreateComponents(t *testing.T) {
	err := CreateComponents(StaticCompConfig{
		URL:     DefaultGenerationURL,
		Method:  DefaultGenerationMethod,
		Path:    WorkloadPath,
		DirName: LatestVersion,
		Config:  NewConfig(LatestVersion),
		Force:   true,
	})
	if err != nil {
		t.Fatalf("Failed to generate components: %s", err.Error())
	}
}
