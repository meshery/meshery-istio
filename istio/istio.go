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
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/layer5io/meshery-istio/meshes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
	iClient.k8sDynamicClient = ic.k8sDynamicClient
	return &meshes.CreateMeshInstanceResponse{}, nil
}

func (iClient *IstioClient) deleteResource(ctx context.Context, res schema.GroupVersionResource, data *unstructured.Unstructured) error {
	if iClient.k8sDynamicClient == nil {
		return errors.New("mesh client has not been created")
	}

	err := iClient.k8sDynamicClient.Resource(res).Namespace(data.GetNamespace()).Delete(data.GetName(), &metav1.DeleteOptions{})
	if err != nil {
		err = errors.Wrapf(err, "unable to delete the requested resource")
		logrus.Error(err)
		return err
	}
	logrus.Infof("Deleted Resource of type: %s and name: %s", data.GetKind(), data.GetName())
	return nil
}

// MeshName just returns the name of the mesh the client is representing
func (iClient *IstioClient) MeshName(context.Context, *meshes.MeshNameRequest) (*meshes.MeshNameResponse, error) {
	return &meshes.MeshNameResponse{Name: "Istio"}, nil
}

func (iClient *IstioClient) applyRulePayload(ctx context.Context, namespace string, newBytes []byte, delete bool) error {
	if iClient.k8sDynamicClient == nil {
		return errors.New("mesh client has not been created")
	}

	jsonBytes, err := yaml.YAMLToJSON(newBytes)
	if err != nil {
		err = errors.Wrapf(err, "unable to convert yaml to json")
		logrus.Error(err)
		return err
	}
	data := &unstructured.Unstructured{}
	err = data.UnmarshalJSON(jsonBytes)
	if err != nil {
		err = errors.Wrapf(err, "unable to unmarshal json created from yaml")
		logrus.Error(err)
		return err
	}
	groupVersion := strings.Split(data.GetAPIVersion(), "/")
	fmt.Printf("groupVersion: %v\n", groupVersion)

	group := groupVersion[0]
	version := groupVersion[1]

	kind := strings.ToLower(data.GetKind())
	if !data.IsList() {
		kind += "s"
	}

	res := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: kind,
	}

	if namespace != "" {
		data.SetNamespace(namespace)
	}

	if delete {
		return iClient.deleteResource(ctx, res, data)
	}

	_, err = iClient.k8sDynamicClient.Resource(res).Namespace(data.GetNamespace()).Create(data, metav1.CreateOptions{})
	if err != nil {
		err = errors.Wrapf(err, "unable to create the requested resource, attempting to update")
		logrus.Warn(err)

		data, err = iClient.k8sDynamicClient.Resource(res).Namespace(data.GetNamespace()).Get(data.GetName(), metav1.GetOptions{})
		if err != nil {
			err = errors.Wrap(err, "unable to retrieve the resource with a matching name, while attempting to apply the config")
			logrus.Error(err)
			return err
		}

		if _, err = iClient.k8sDynamicClient.Resource(res).Namespace(data.GetNamespace()).Update(data, metav1.UpdateOptions{}); err != nil {
			err = errors.Wrap(err, "unable to update resource with the given name, while attempting to apply the config")
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
	op, ok := supportedOps[arReq.OpName]
	if !ok {
		return nil, fmt.Errorf("error: %s is not a valid operation name", arReq.OpName)
	}

	if arReq.OpName == customOpName && arReq.CustomBody == "" {
		return nil, fmt.Errorf("error: yaml body is empty for %s operation", arReq.OpName)
	}

	yamlFile := ""
	if arReq.OpName == customOpName {
		yamlFile = arReq.CustomBody
	} else {
		tmpl, err := template.ParseFiles(path.Join("istio", "config_templates", op.templateName))
		if err != nil {
			err = errors.Wrapf(err, "unable to parse template")
			logrus.Error(err)
			return nil, err
		}
		buf := bytes.NewBufferString("")
		err = tmpl.Execute(buf, map[string]string{
			"user_name": arReq.Username,
			"namespace": arReq.Namespace,
		})
		if err != nil {
			err = errors.Wrapf(err, "unable to execute template")
			logrus.Error(err)
			return nil, err
		}
		yamlFile = buf.String()
	}
	if err := iClient.applyConfigChange(ctx, yamlFile, arReq.Namespace, arReq.DeleteOp); err != nil {
		return nil, err
	}
	return &meshes.ApplyRuleResponse{}, nil
}

func (iClient *IstioClient) applyConfigChange(ctx context.Context, yamlFile, namespace string, delete bool) error {
	yamls := strings.Split(yamlFile, "---")

	for _, yml := range yamls {
		if strings.TrimSpace(yml) != "" {
			if err := iClient.applyRulePayload(ctx, namespace, []byte(yml), delete); err != nil {
				return err
			}
		}
	}
	return nil
}

// SupportedOperations - returns a list of supported operations on the mesh
func (iClient *IstioClient) SupportedOperations(context.Context, *meshes.SupportedOperationsRequest) (*meshes.SupportedOperationsResponse, error) {
	result := map[string]string{}
	for key, op := range supportedOps {
		result[key] = op.name
	}
	return &meshes.SupportedOperationsResponse{
		Ops: result,
	}, nil
}
