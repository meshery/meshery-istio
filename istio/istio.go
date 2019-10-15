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
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/ghodss/yaml"
	"github.com/layer5io/meshery-istio/meshes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	hipsterShopIstioManifestsURL      = "https://raw.githubusercontent.com/GoogleCloudPlatform/microservices-demo/master/release/istio-manifests.yaml"
	hipsterShopKubernetesManifestsURL = "https://raw.githubusercontent.com/GoogleCloudPlatform/microservices-demo/master/release/kubernetes-manifests.yaml"
)

//CreateMeshInstance is called from UI
func (iClient *Client) CreateMeshInstance(_ context.Context, k8sReq *meshes.CreateMeshInstanceRequest) (*meshes.CreateMeshInstanceResponse, error) {
	var k8sConfig []byte
	contextName := ""
	if k8sReq != nil {
		k8sConfig = k8sReq.K8SConfig
		contextName = k8sReq.ContextName
	}
	// logrus.Debugf("received k8sConfig: %s", k8sConfig)
	logrus.Debugf("received contextName: %s", contextName)

	ic, err := newClient(k8sConfig, contextName)
	if err != nil {
		err = errors.Wrapf(err, "unable to create a new istio client")
		logrus.Error(err)
		return nil, err
	}
	iClient.k8sClientset = ic.k8sClientset
	iClient.k8sDynamicClient = ic.k8sDynamicClient
	iClient.eventChan = make(chan *meshes.EventsResponse, 100)
	iClient.config = ic.config
	return &meshes.CreateMeshInstanceResponse{}, nil
}

func (iClient *Client) createResource(ctx context.Context, res schema.GroupVersionResource, data *unstructured.Unstructured) error {
	_, err := iClient.k8sDynamicClient.Resource(res).Namespace(data.GetNamespace()).Create(data, metav1.CreateOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			// 	if err1 := iClient.deleteResource(ctx, res, data); err1 != nil {

			// 	}
			return errors.Wrap(err, "resource already exists")
		}
		err = errors.Wrapf(err, "unable to create the requested resource, attempting operation without namespace")
		logrus.Warn(err)
		if _, err = iClient.k8sDynamicClient.Resource(res).Create(data, metav1.CreateOptions{}); err != nil {
			if strings.Contains(err.Error(), "already exists") {
				return errors.Wrap(err, "resource already exists")
			}
			err = errors.Wrapf(err, "unable to create the requested resource, attempting to update")
			logrus.Error(err)
			return err
		}
	}
	logrus.Infof("Created Resource of type: %s and name: %s", data.GetKind(), data.GetName())
	return nil
}

func (iClient *Client) deleteResource(ctx context.Context, res schema.GroupVersionResource, data *unstructured.Unstructured) error {
	if iClient.k8sDynamicClient == nil {
		return errors.New("mesh client has not been created")
	}

	if res.Resource == "namespaces" && data.GetName() == "default" { // skipping deletion of default namespace
		return nil
	}

	// in the case with deployments, have to scale it down to 0 first and then delete. . . or else RS and pods will be left behind
	if res.Resource == "deployments" {
		data1, err := iClient.getResource(ctx, res, data)
		if err != nil {
			return err
		}
		depl := data1.UnstructuredContent()
		spec1 := depl["spec"].(map[string]interface{})
		spec1["replicas"] = 0
		data1.SetUnstructuredContent(depl)
		if err = iClient.updateResource(ctx, res, data1); err != nil {
			return err
		}
	}

	err := iClient.k8sDynamicClient.Resource(res).Namespace(data.GetNamespace()).Delete(data.GetName(), &metav1.DeleteOptions{})
	if err != nil {
		err = errors.Wrapf(err, "unable to delete the requested resource, attempting operation without namespace")
		logrus.Warn(err)

		err := iClient.k8sDynamicClient.Resource(res).Delete(data.GetName(), &metav1.DeleteOptions{})
		if err != nil {
			err = errors.Wrapf(err, "unable to delete the requested resource")
			logrus.Error(err)
			return err
		}
	}
	logrus.Infof("Deleted Resource of type: %s and name: %s", data.GetKind(), data.GetName())
	return nil
}

func (iClient *Client) getResource(ctx context.Context, res schema.GroupVersionResource, data *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	data1, err := iClient.k8sDynamicClient.Resource(res).Namespace(data.GetNamespace()).Get(data.GetName(), metav1.GetOptions{})
	if err != nil {
		err = errors.Wrap(err, "unable to retrieve the resource with a matching name, attempting operation without namespace")
		logrus.Warn(err)

		data1, err = iClient.k8sDynamicClient.Resource(res).Get(data.GetName(), metav1.GetOptions{})
		if err != nil {
			err = errors.Wrap(err, "unable to retrieve the resource with a matching name, while attempting to apply the config")
			logrus.Error(err)
			return nil, err
		}
	}
	logrus.Infof("Retrieved Resource of type: %s and name: %s", data.GetKind(), data.GetName())
	return data1, nil
}

func (iClient *Client) updateResource(ctx context.Context, res schema.GroupVersionResource, data *unstructured.Unstructured) error {
	if _, err := iClient.k8sDynamicClient.Resource(res).Namespace(data.GetNamespace()).Update(data, metav1.UpdateOptions{}); err != nil {
		if strings.Contains(err.Error(), "the server does not allow this method on the requested resource") {
			logrus.Error(err)
			return err
		}
		err = errors.Wrap(err, "unable to update resource with the given name, attempting operation without namespace")
		logrus.Warn(err)

		if _, err = iClient.k8sDynamicClient.Resource(res).Update(data, metav1.UpdateOptions{}); err != nil {
			if strings.Contains(err.Error(), "the server does not allow this method on the requested resource") {
				logrus.Error(err)
				return err
			}
			err = errors.Wrap(err, "unable to update resource with the given name, while attempting to apply the config")
			logrus.Error(err)
			return err
		}
	}
	logrus.Infof("Updated Resource of type: %s and name: %s", data.GetKind(), data.GetName())
	return nil
}

// MeshName just returns the name of the mesh the client is representing
func (iClient *Client) MeshName(context.Context, *meshes.MeshNameRequest) (*meshes.MeshNameResponse, error) {
	return &meshes.MeshNameResponse{Name: "Istio"}, nil
}

func (iClient *Client) applyRulePayload(ctx context.Context, namespace string, newBytes []byte, delete, isCustomOp bool) error {
	if iClient.k8sDynamicClient == nil {
		return errors.New("mesh client has not been created")
	}
	// logrus.Debugf("received yaml bytes: %s", newBytes)
	jsonBytes, err := yaml.YAMLToJSON(newBytes)
	if err != nil {
		err = errors.Wrapf(err, "unable to convert yaml to json")
		logrus.Error(err)
		return err
	}
	// logrus.Debugf("created json: %s, length: %d", jsonBytes, len(jsonBytes))
	if len(jsonBytes) > 5 { // attempting to skip 'null' json
		data := &unstructured.Unstructured{}
		err = data.UnmarshalJSON(jsonBytes)
		if err != nil {
			err = errors.Wrapf(err, "unable to unmarshal json created from yaml")
			logrus.Error(err)
			return err
		}
		if data.IsList() {
			err = data.EachListItem(func(r runtime.Object) error {
				dataL, _ := r.(*unstructured.Unstructured)
				return iClient.executeRule(ctx, dataL, namespace, delete, isCustomOp)
			})
			return err
		}
		return iClient.executeRule(ctx, data, namespace, delete, isCustomOp)
	}
	return nil
}

func (iClient *Client) executeRule(ctx context.Context, data *unstructured.Unstructured, namespace string, delete, isCustomOp bool) error {
	// logrus.Debug("========================================================")
	// logrus.Debugf("Received data: %+#v", data)
	if namespace != "" {
		data.SetNamespace(namespace)
	}
	groupVersion := strings.Split(data.GetAPIVersion(), "/")
	logrus.Debugf("groupVersion: %v", groupVersion)
	var group, version string
	if len(groupVersion) == 2 {
		group = groupVersion[0]
		version = groupVersion[1]
	} else if len(groupVersion) == 1 {
		version = groupVersion[0]
	}

	kind := strings.ToLower(data.GetKind())
	switch kind {
	case "logentry":
		kind = "logentries"
	case "kubernetes":
		kind = "kuberneteses"
	case "podsecuritypolicy":
		kind = "podsecuritypolicies"
	case "serviceentry":
		kind = "serviceentries"
	default:
		kind += "s"
	}

	res := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: kind,
	}
	logrus.Debugf("Computed Resource: %+#v", res)

	if delete {
		return iClient.deleteResource(ctx, res, data)
	}
	trackRetry := 0
RETRY:
	if err := iClient.createResource(ctx, res, data); err != nil {
		if isCustomOp {
			if err := iClient.deleteResource(ctx, res, data); err != nil {
				return err
			}
			time.Sleep(time.Second)
			if err := iClient.createResource(ctx, res, data); err != nil {
				return err
			}
			// data1, err := iClient.getResource(ctx, res, data)
			// if err != nil {
			// 	return err
			// }
			// if err = iClient.updateResource(ctx, res, data1); err != nil {
			// 	return err
			// }
		} else {
			data1, err := iClient.getResource(ctx, res, data)
			if err != nil {
				return err
			}
			data.SetCreationTimestamp(data1.GetCreationTimestamp())
			data.SetGenerateName(data1.GetGenerateName())
			data.SetGeneration(data1.GetGeneration())
			data.SetSelfLink(data1.GetSelfLink())
			data.SetResourceVersion(data1.GetResourceVersion())
			// data.DeepCopyInto(data1)
			if err = iClient.updateResource(ctx, res, data); err != nil {
				if strings.Contains(err.Error(), "the server does not allow this method on the requested resource") {
					logrus.Info("attempting to delete resource. . . ")
					iClient.deleteResource(ctx, res, data)
					trackRetry++
					if trackRetry <= 3 {
						goto RETRY
					} // else return error
				}
				return err
			}
			// return err
		}
	}
	return nil
}

func (iClient *Client) applyIstioCRDs(ctx context.Context, delete bool) error {
	crdYAMLs, err := iClient.getCRDsYAML()
	if err != nil {
		return err
	}
	logrus.Debug("processing crds. . .")
	for _, crdYAML := range crdYAMLs {
		if err := iClient.applyConfigChange(ctx, crdYAML, "", delete, false); err != nil {
			return err
		}
	}
	return nil
}

func (iClient *Client) labelNamespaceForAutoInjection(ctx context.Context, namespace string) error {
	ns := &unstructured.Unstructured{}
	res := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "namespaces",
	}
	ns.SetName(namespace)
	ns, err := iClient.getResource(ctx, res, ns)
	if err != nil {
		if strings.HasSuffix(err.Error(), "not found") {
			if err = iClient.createNamespace(ctx, namespace); err != nil {
				return err
			}

			ns := &unstructured.Unstructured{}
			ns.SetName(namespace)
			ns, err = iClient.getResource(ctx, res, ns)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	logrus.Debugf("retrieved namespace: %+#v", ns)
	if ns == nil {
		ns = &unstructured.Unstructured{}
		ns.SetName(namespace)
	}
	ns.SetLabels(map[string]string{
		"istio-injection": "enabled",
	})
	err = iClient.updateResource(ctx, res, ns)
	if err != nil {
		return err
	}
	return nil
}

func (iClient *Client) createNamespace(ctx context.Context, namespace string) error {
	logrus.Debugf("creating namespace: %s", namespace)
	yamlFileContents, err := iClient.executeTemplate(ctx, "", namespace, "namespace.yml")
	if err != nil {
		return err
	}
	if err := iClient.applyConfigChange(ctx, yamlFileContents, namespace, false, false); err != nil {
		return err
	}
	return nil
}

func (iClient *Client) executeTemplate(ctx context.Context, username, namespace, templateName string) (string, error) {
	tmpl, err := template.ParseFiles(path.Join("istio", "config_templates", templateName))
	if err != nil {
		err = errors.Wrapf(err, "unable to parse template")
		logrus.Error(err)
		return "", err
	}
	buf := bytes.NewBufferString("")
	err = tmpl.Execute(buf, map[string]string{
		"user_name": username,
		"namespace": namespace,
	})
	if err != nil {
		err = errors.Wrapf(err, "unable to execute template")
		logrus.Error(err)
		return "", err
	}
	return buf.String(), nil
}

func (iClient *Client) executeInstall(ctx context.Context, installmTLS bool, arReq *meshes.ApplyRuleRequest) error {
	arReq.Namespace = ""
	if arReq.DeleteOp {
		defer iClient.applyIstioCRDs(ctx, arReq.DeleteOp)
	} else {
		if err := iClient.applyIstioCRDs(ctx, arReq.DeleteOp); err != nil {
			return err
		}
	}
	yamlFileContents, err := iClient.getLatestIstioYAML(installmTLS)
	if err != nil {
		return err
	}
	if err := iClient.applyConfigChange(ctx, yamlFileContents, arReq.Namespace, arReq.DeleteOp, false); err != nil {
		return err
	}
	return nil
}

func (iClient *Client) executeBookInfoInstall(ctx context.Context, arReq *meshes.ApplyRuleRequest) error {
	if !arReq.DeleteOp {
		if err := iClient.labelNamespaceForAutoInjection(ctx, arReq.Namespace); err != nil {
			return err
		}
	}
	yamlFileContents, err := iClient.getBookInfoAppYAML()
	if err != nil {
		return err
	}
	if err := iClient.applyConfigChange(ctx, yamlFileContents, arReq.Namespace, arReq.DeleteOp, false); err != nil {
		return err
	}
	yamlFileContents, err = iClient.getBookInfoGatewayYAML()
	if err != nil {
		return err
	}
	if err := iClient.applyConfigChange(ctx, yamlFileContents, arReq.Namespace, arReq.DeleteOp, false); err != nil {
		return err
	}
	return nil
}

func (iClient *Client) executeHipsterShopInstall(ctx context.Context, arReq *meshes.ApplyRuleRequest) error {
	if !arReq.DeleteOp {
		if err := iClient.labelNamespaceForAutoInjection(ctx, arReq.Namespace); err != nil {
			return err
		}
	}
	hipsterShopFilecontents := func(fileURL string) (string, error) {
		resp, err := http.Get(fileURL)
		if err != nil {
			err = errors.Wrapf(err, "error getting data from %s", fileURL)
			logrus.Error(err)
			return "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				err = errors.Wrapf(err, "error parsing response from %s", fileURL)
				logrus.Error(err)
				return "", err
			}
			return string(body), nil
		}
		err = errors.Wrapf(err, "Call failed with response status: %s", resp.Status)
		logrus.Error(err)
		return "", err
	}

	kubernetesManifestsContent, err := hipsterShopFilecontents(hipsterShopKubernetesManifestsURL)
	if err != nil {
		return err
	}
	istioManifestsContent, err := hipsterShopFilecontents(hipsterShopIstioManifestsURL)
	if err != nil {
		return err
	}

	var yamlFileContents = fmt.Sprintf("%s\n---\n%s", kubernetesManifestsContent, istioManifestsContent)

	if err := iClient.applyConfigChange(ctx, yamlFileContents, arReq.Namespace, arReq.DeleteOp, false); err != nil {
		return err
	}

	return nil
}

// ApplyOperation is a method invoked to apply a particular operation on the mesh in a namespace
func (iClient *Client) ApplyOperation(ctx context.Context, arReq *meshes.ApplyRuleRequest) (*meshes.ApplyRuleResponse, error) {
	if arReq == nil {
		return nil, errors.New("mesh client has not been created")
	}

	op, ok := supportedOps[arReq.OpName]
	if !ok {
		return nil, fmt.Errorf("operation id: %s, error: %s is not a valid operation name", arReq.OperationId, arReq.OpName)
	}

	if arReq.OpName == customOpCommand && arReq.CustomBody == "" {
		return nil, fmt.Errorf("operation id: %s, error: yaml body is empty for %s operation", arReq.OperationId, arReq.OpName)
	}

	var yamlFileContents string
	var err error
	installWithmTLS := false
	isCustomOp := false

	switch arReq.OpName {
	case installmTLSIstioCommand:
		installWithmTLS = true
		fallthrough
	case installIstioCommand:
		go func() {
			opName1 := "deploying"
			if arReq.DeleteOp {
				opName1 = "removing"
			}
			if err := iClient.executeInstall(ctx, installWithmTLS, arReq); err != nil {
				iClient.eventChan <- &meshes.EventsResponse{
					OperationId: arReq.OperationId,
					EventType:   meshes.EventType_ERROR,
					Summary:     fmt.Sprintf("Error while %s Istio", opName1),
					Details:     err.Error(),
				}
				return
			}
			opName := "deployed"
			if arReq.DeleteOp {
				opName = "removed"
			}
			iClient.eventChan <- &meshes.EventsResponse{
				OperationId: arReq.OperationId,
				EventType:   meshes.EventType_INFO,
				Summary:     fmt.Sprintf("Istio %s successfully", opName),
				Details:     fmt.Sprintf("The latest version of Istio is now %s.", opName),
			}
			return
		}()
		return &meshes.ApplyRuleResponse{
			OperationId: arReq.OperationId,
		}, nil
	case googleMSSampleApplication:
		go func() {
			opName1 := "deploying"
			if arReq.DeleteOp {
				opName1 = "removing"
			}
			if err := iClient.executeHipsterShopInstall(ctx, arReq); err != nil {
				iClient.eventChan <- &meshes.EventsResponse{
					OperationId: arReq.OperationId,
					EventType:   meshes.EventType_ERROR,
					Summary:     fmt.Sprintf("Error while %s the Hipster Shop application", opName1),
					Details:     err.Error(),
				}
				return
			}
			opName := "deployed"
			if arReq.DeleteOp {
				opName = "removed"
			}
			iClient.eventChan <- &meshes.EventsResponse{
				OperationId: arReq.OperationId,
				EventType:   meshes.EventType_INFO,
				Summary:     fmt.Sprintf("The Hipster Shop application %s successfully", opName),
				Details:     fmt.Sprintf("The Hipster Shop is now %s.", opName),
			}
			return
		}()
		return &meshes.ApplyRuleResponse{
			OperationId: arReq.OperationId,
		}, nil

	case installBookInfoCommand:
		go func() {
			opName1 := "deploying"
			if arReq.DeleteOp {
				opName1 = "removing"
			}
			if err := iClient.executeBookInfoInstall(ctx, arReq); err != nil {
				iClient.eventChan <- &meshes.EventsResponse{
					OperationId: arReq.OperationId,
					EventType:   meshes.EventType_ERROR,
					Summary:     fmt.Sprintf("Error while %s the canonical Book Info App", opName1),
					Details:     err.Error(),
				}
				return
			}
			opName := "deployed"
			if arReq.DeleteOp {
				opName = "removed"
			}
			iClient.eventChan <- &meshes.EventsResponse{
				OperationId: arReq.OperationId,
				EventType:   meshes.EventType_INFO,
				Summary:     fmt.Sprintf("Book Info app %s successfully", opName),
				Details:     fmt.Sprintf("The Istio canonical Book Info app is now %s.", opName),
			}
			return
		}()
		return &meshes.ApplyRuleResponse{
			OperationId: arReq.OperationId,
		}, nil
	case bookInfoDefaultDestinationRules:
		yamlFileContents, err = iClient.getBookInfoDefaultDesinationRulesYAML()
		if err != nil {
			return nil, err
		}
	case bookInfoRouteToV1AllServices:
		yamlFileContents, err = iClient.getBookInfoRouteToV1AllServicesYAML()
		if err != nil {
			return nil, err
		}
	case bookInfoRouteToReviewsV2ForJason:
		yamlFileContents, err = iClient.getBookInfoRouteToReviewsV2ForJasonFile()
		if err != nil {
			return nil, err
		}
	case bookInfoCanary50pcReviewsV3:
		yamlFileContents, err = iClient.getBookInfoCanary50pcReviewsV3File()
		if err != nil {
			return nil, err
		}
	case bookInfoCanary100pcReviewsV3:
		yamlFileContents, err = iClient.getBookInfoCanary100pcReviewsV3File()
		if err != nil {
			return nil, err
		}
	case bookInfoInjectDelayForRatingsForJason:
		yamlFileContents, err = iClient.getBookInfoInjectDelayForRatingsForJasonFile()
		if err != nil {
			return nil, err
		}
	case bookInfoInjectHTTPAbortToRatingsForJason:
		yamlFileContents, err = iClient.getBookInfoInjectHTTPAbortToRatingsForJasonFile()
		if err != nil {
			return nil, err
		}
	case installSMI:
		if !arReq.DeleteOp && arReq.Namespace != "default" {
			iClient.createNamespace(ctx, arReq.Namespace)
		}
		yamlFileContents, err = getSMIYamls()
		if err != nil {
			return nil, err
		}
	case runVet:
		go iClient.runVet()
		return &meshes.ApplyRuleResponse{
			OperationId: arReq.OperationId,
		}, nil
	case customOpCommand:
		yamlFileContents = arReq.CustomBody
		isCustomOp = true
	default:
		if !arReq.DeleteOp {
			if err := iClient.labelNamespaceForAutoInjection(ctx, arReq.Namespace); err != nil {
				return nil, err
			}
		}
		yamlFileContents, err = iClient.executeTemplate(ctx, arReq.Username, arReq.Namespace, op.templateName)
		if err != nil {
			return nil, err
		}
	}

	go func() {
		logrus.Debug("in the routine. . . .")
		opName1 := "deploying"
		if arReq.DeleteOp {
			opName1 = "removing"
		}
		if err := iClient.applyConfigChange(ctx, yamlFileContents, arReq.Namespace, arReq.DeleteOp, isCustomOp); err != nil {
			iClient.eventChan <- &meshes.EventsResponse{
				OperationId: arReq.OperationId,
				EventType:   meshes.EventType_ERROR,
				Summary:     fmt.Sprintf("Error while %s \"%s\"", opName1, op.name),
				Details:     err.Error(),
			}
			return
		}
		opName := "deployed"
		if arReq.DeleteOp {
			opName = "removed"
		}
		iClient.eventChan <- &meshes.EventsResponse{
			OperationId: arReq.OperationId,
			EventType:   meshes.EventType_INFO,
			Summary:     fmt.Sprintf("\"%s\" %s successfully", op.name, opName),
			Details:     fmt.Sprintf("\"%s\" %s successfully", op.name, opName),
		}
	}()

	return &meshes.ApplyRuleResponse{
		OperationId: arReq.OperationId,
	}, nil
}

func (iClient *Client) applyConfigChange(ctx context.Context, yamlFileContents, namespace string, delete, isCustomOp bool) error {
	// yamls := strings.Split(yamlFileContents, "---")
	yamls, err := iClient.splitYAML(yamlFileContents)
	if err != nil {
		err = errors.Wrap(err, "error while splitting yaml")
		logrus.Error(err)
		return err
	}
	for _, yml := range yamls {
		if strings.TrimSpace(yml) != "" {
			if err := iClient.applyRulePayload(ctx, namespace, []byte(yml), delete, isCustomOp); err != nil {
				errStr := strings.TrimSpace(err.Error())
				if delete {
					if strings.HasSuffix(errStr, "not found") ||
						strings.HasSuffix(errStr, "the server could not find the requested resource") {
						// logrus.Debugf("skipping error. . .")
						continue
					}
				} else {
					if strings.HasSuffix(errStr, "already exists") {
						continue
					}
				}
				// logrus.Debugf("returning error: %v", err)
				return err
			}
		}
	}
	return nil
}

// SupportedOperations - returns a list of supported operations on the mesh
func (iClient *Client) SupportedOperations(context.Context, *meshes.SupportedOperationsRequest) (*meshes.SupportedOperationsResponse, error) {
	supportedOpsCount := len(supportedOps)
	result := make([]*meshes.SupportedOperation, supportedOpsCount)
	i := 0
	for k, sp := range supportedOps {
		result[i] = &meshes.SupportedOperation{
			Key:      k,
			Value:    sp.name,
			Category: sp.opType,
		}
		i++
	}
	return &meshes.SupportedOperationsResponse{
		Ops: result,
	}, nil
}

// StreamEvents - streams generated/collected events to the client
func (iClient *Client) StreamEvents(in *meshes.EventsRequest, stream meshes.MeshService_StreamEventsServer) error {
	logrus.Debugf("waiting on event stream. . .")
	for {
		select {
		case event := <-iClient.eventChan:
			logrus.Debugf("sending event: %+#v", event)
			if err := stream.Send(event); err != nil {
				err = errors.Wrapf(err, "unable to send event")

				// to prevent loosing the event, will re-add to the channel
				go func() {
					iClient.eventChan <- event
				}()
				logrus.Error(err)
				return err
			}
		default:
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func (iClient *Client) splitYAML(yamlContents string) ([]string, error) {
	yamlDecoder, ok := NewDocumentDecoder(ioutil.NopCloser(bytes.NewReader([]byte(yamlContents)))).(*YAMLDecoder)
	if !ok {
		err := fmt.Errorf("unable to create a yaml decoder")
		logrus.Error(err)
		return nil, err
	}
	defer yamlDecoder.Close()
	var err error
	n := 0
	data := [][]byte{}
	ind := 0
	for err == io.ErrShortBuffer || err == nil {
		// for {
		d := make([]byte, 1000)
		n, err = yamlDecoder.Read(d)
		// logrus.Debugf("Read this: %s, count: %d, err: %v", d, n, err)
		if len(data) == 0 || len(data) <= ind {
			data = append(data, []byte{})
		}
		if n > 0 {
			data[ind] = append(data[ind], d...)
		}
		if err == nil {
			logrus.Debugf("..............BOUNDARY................")
			ind++
		}
	}
	result := make([]string, len(data))
	for i, row := range data {
		r := string(row)
		r = strings.Trim(r, "\x00")
		logrus.Debugf("ind: %d, data: %s", i, r)
		result[i] = r
	}
	return result, nil
}
