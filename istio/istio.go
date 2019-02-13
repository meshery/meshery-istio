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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/layer5io/meshery-istio/meshes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime"
)

// 	SupportedOperations(context.Context, *SupportedOperationsRequest) (*SupportedOperationsResponse, error)

func (iClient *IstioClient) CreateMeshInstance(_ context.Context, k8sReq *meshes.CreateMeshInstanceRequest) (*meshes.CreateMeshInstanceResponse, error) {
	var k8sConfig []byte
	contextName := ""
	if k8sReq != nil {
		k8sConfig = k8sReq.K8SConfig
		contextName = k8sReq.ContextName
	}
	logrus.Debugf("received k8sConfig: %s", k8sConfig)
	logrus.Debugf("received contextName: %s", contextName)

	ic, err := newClient(k8sConfig, contextName)
	if err != nil {
		err = errors.Wrapf(err, "unable to create a new istio client")
		logrus.Error(err)
		return nil, err
	}
	iClient.k8s = ic.k8s
	iClient.istioConfigApi = ic.istioConfigApi
	iClient.istioNetworkingApi = ic.istioNetworkingApi
	return &meshes.CreateMeshInstanceResponse{}, nil
}

func (iClient *IstioClient) deleteAllCreatedResources(ctx context.Context, namespace string) {
	resourceNames := []string{"productpage", "ratings", "reviews", "details"}
	for _, rs := range resourceNames {
		iClient.deleteResource(ctx, virtualServices, namespace, rs)

	}
}

func (iClient *IstioClient) deleteResource(ctx context.Context, overallType, namespace, resName string) error {
	if iClient.istioNetworkingApi == nil {
		return errors.New("mesh client has not been created")
	}

	newRes := iClient.istioNetworkingApi.Delete().
		Namespace(namespace).
		Resource(overallType).SubResource(resName).Do()
	_, err := newRes.Get()
	if err != nil {
		err = errors.Wrapf(err, "unable to delete the requested resource")
		logrus.Error(err)
		return err
	}
	logrus.Infof("Deleted Resource of type: %s and name: %s", overallType, resName)
	return nil
}

// MeshName just returns the name of the mesh the client is representing
func (iClient *IstioClient) MeshName(context.Context, *meshes.MeshNameRequest) (*meshes.MeshNameResponse, error) {
	return &meshes.MeshNameResponse{Name: "Istio"}, nil
}

func (iClient *IstioClient) applyRulePayload(ctx context.Context, namespace string, newVSBytes []byte) error {
	if iClient.istioNetworkingApi == nil {
		return errors.New("mesh client has not been created")
	}
	vs := &VirtualService{}
	err := yaml.Unmarshal(newVSBytes, vs)
	if err != nil {
		err = errors.Wrapf(err, "unable to unmarshal yaml")
		logrus.Error(err)
		return err
	}
	vs.Kind = virtualservice
	vs.APIVersion = istioNetworkingGroupVersion.String()
	newVSBytesJ, err := json.Marshal(vs)
	if err != nil {
		err = errors.Wrapf(err, "unable to marshal virtual service map")
		logrus.Error(err)
		return err
	}

	newVSRes := iClient.istioNetworkingApi.Post().SetHeader("content-type", runtime.ContentTypeJSON).
		Namespace(namespace).
		Resource(virtualServices).Body(newVSBytesJ).Do()
	_, err = newVSRes.Get()
	if err != nil {
		newVSRes = iClient.istioNetworkingApi.Get().SetHeader("content-type", runtime.ContentTypeJSON).
			Namespace(namespace).Name(vs.ObjectMeta.Name).
			Resource(virtualServices).Do()
		newVSResInst, err := newVSRes.Get()
		if err != nil {
			err = errors.Wrapf(err, "unable to get the virtual service instance")
			logrus.Error(err)
			return err
		}
		vs1, _ := newVSResInst.(*VirtualService)
		vs.ObjectMeta.ResourceVersion = vs1.ObjectMeta.ResourceVersion
		newVSBytesJ, err := json.Marshal(vs)
		if err != nil {
			err = errors.Wrapf(err, "unable to marshal virtual service map")
			logrus.Error(err)
			return err
		}

		if _, err = iClient.istioNetworkingApi.Put().SetHeader("content-type", runtime.ContentTypeJSON).
			Namespace(namespace).Name(vs.ObjectMeta.Name).
			Resource(virtualServices).Body(newVSBytesJ).Do().Get(); err != nil {
			err = errors.Wrapf(err, "unable to get the virtual service instance from result")
			logrus.Error(err)
			return err
		}
	}
	return nil
}

// ApplyRule is a method invoked to apply a particular operation on the mesh in a namespace
func (iClient *IstioClient) ApplyOperation(ctx context.Context, arReq *meshes.ApplyRuleRequest) (*meshes.ApplyRuleResponse, error) {
	if arReq == nil {
		return nil, errors.New("mesh client has not been created")
	}

	// ApplyRule(ctx context.Context, opName, username, namespace string) error {
	yamlFile := ""
	reset := false
	for _, op := range supportedOps {
		if op.key == arReq.OpName {
			yamlFile = op.templateName
			if op.resetOp == true {
				reset = true
			}
		}
	}
	if yamlFile == "" && !reset {
		return nil, fmt.Errorf("error: %s is not a valid operation name", arReq.OpName)
	}
	if reset {
		iClient.deleteAllCreatedResources(ctx, arReq.Namespace)
		return &meshes.ApplyRuleResponse{}, nil
	}
	if err := iClient.applyConfigChange(ctx, path.Join("istio", "config_templates", yamlFile), arReq.Username, arReq.Namespace); err != nil {
		return nil, err
	}
	return &meshes.ApplyRuleResponse{}, nil
}

func (iClient *IstioClient) applyConfigChange(ctx context.Context, yamlFile, username, namespace string) error {
	iClient.deleteAllCreatedResources(ctx, namespace)

	tmpl := template.Must(template.ParseFiles(yamlFile))

	buf := bytes.NewBufferString("")
	err := tmpl.Execute(buf, map[string]string{
		"user_name": username,
		"namespace": namespace,
	})
	if err != nil {
		err = errors.Wrapf(err, "unable to parse template")
		logrus.Error(err)
		return err
	}
	completeYaml := buf.String()
	yamls := strings.Split(completeYaml, "---")

	for _, yml := range yamls {
		if strings.TrimSpace(yml) != "" {
			if err := iClient.applyRulePayload(ctx, namespace, []byte(yml)); err != nil {
				return err
			}
		}
	}
	return nil
}

// SupportedOperations - returns a list of supported operations on the mesh
func (iClient *IstioClient) SupportedOperations(context.Context, *meshes.SupportedOperationsRequest) (*meshes.SupportedOperationsResponse, error) {
	// Operations(ctx context.Context) (map[string]string, error) {
	result := map[string]string{}
	for _, op := range supportedOps {
		result[op.key] = op.name
	}
	return &meshes.SupportedOperationsResponse{
		Ops: result,
	}, nil
}
