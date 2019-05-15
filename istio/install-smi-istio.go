package istio

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const smiBasePath = "istio/config_templates/smi/"

var smiFiles = []string{
	"split_v1alpha1_trafficsplit_crd.yaml",
	"rbac.yaml",
	"operator.yaml",
}

func getSMIYamls() (string, error) {
	var result strings.Builder
	for _, smiFile := range smiFiles {
		fileContents, err := ioutil.ReadFile(path.Join(smiBasePath, smiFile))
		if err != nil {
			err = errors.Wrap(err, "unable to read file")
			logrus.Error(err)
			return "", err
		}
		result.Write(fileContents)
		result.WriteString("\n---\n")
	}
	logrus.Debugf("generated yaml: %s", result.String())
	return result.String(), nil
}
