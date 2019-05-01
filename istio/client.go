// Copyright 2019 Layer5.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package istio

import (
	"time"

	"github.com/layer5io/meshery-istio/meshes"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/ghodss/yaml"
)

// IstioClient represents an Istio client in Meshery
type IstioClient struct {
	config           *rest.Config
	k8sClientset     *kubernetes.Clientset
	k8sDynamicClient dynamic.Interface
	eventChan        chan *meshes.EventsResponse

	istioReleaseVersion     string
	istioReleaseDownloadURL string
	istioReleaseUpdatedAt   time.Time
}

func configClient(kubeconfig []byte, contextName string) (*rest.Config, error) {
	if len(kubeconfig) > 0 {
		ccfg, err := clientcmd.Load(kubeconfig)
		if err != nil {
			return nil, err
		}
		if contextName != "" {
			ccfg.CurrentContext = contextName
		}

		return clientcmd.NewDefaultClientConfig(*ccfg, &clientcmd.ConfigOverrides{}).ClientConfig()
	}
	return rest.InClusterConfig()
}

func newClient(kubeconfig []byte, contextName string) (*IstioClient, error) {
	kubeconfig = monkeyPatchingToSupportInsecureConn(kubeconfig)
	client := IstioClient{}
	config, err := configClient(kubeconfig, contextName)
	if err != nil {
		return nil, err
	}
	config.QPS = 100
	config.Burst = 200

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	client.k8sDynamicClient = dynamicClient

	k8sClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	client.k8sClientset = k8sClientset
	client.config = config

	return &client, nil
}

func monkeyPatchingToSupportInsecureConn(data []byte) []byte {
	config := map[string]interface{}{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		logrus.Warn(err)
		return data // we will skip this process
	}
	// logrus.Infof("unmarshalled config: %+#v", config)
	clusters, ok := config["clusters"].([]interface{})
	if !ok {
		logrus.Warn("unable to type cast clusters to a map array")
		return data
	}
	for _, clusterI := range clusters {
		cluster, ok := clusterI.(map[string]interface{})
		if !ok {
			logrus.Warn("unable to type case individual cluster to a map")
			continue
		}
		indCluster, ok := cluster["cluster"].(map[string]interface{})
		if !ok {
			logrus.Warn("unable to type case clusters.cluster to a map")
			continue
		}
		indCluster["insecure-skip-tls-verify"] = true // TODO: should we persist this back?
		delete(indCluster, "certificate-authority-data")
		delete(indCluster, "certificate-authority")
	}
	// logrus.Debugf("New config: %+#v", config)
	data1, err := yaml.Marshal(config)
	if err != nil {
		logrus.Warn(err)
		return data
	}
	return data1
}
