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

import "github.com/layer5io/meshery-istio/meshes"

type supportedOperation struct {
	// a friendly name
	name string
	// the template file name
	templateName string
	opType       meshes.OpCategory
}

const (
	customOpCommand                          = "custom"
	runVet                                   = "istio_vet"
	installIstioCommand                      = "istio_install"
	installmTLSIstioCommand                  = "istio_mtls_install"
	installBookInfoCommand                   = "install_book_info"
	cbCommand                                = "cb1"
	installSMI                               = "install_smi"
	installHTTPBin                           = "install_http_bin"
	googleMSSampleApplication                = "google_microservices_demo_application"
	bookInfoDefaultDestinationRules          = "bookInfoDefaultDestinationRules"
	bookInfoRouteToV1AllServices             = "bookInfoRouteToV1AllServices"
	bookInfoRouteToReviewsV2ForJason         = "bookInfoRouteToReviewsV2ForJason"
	bookInfoCanary50pcReviewsV3              = "bookInfoCanary50pcReviewsV3"
	bookInfoCanary100pcReviewsV3             = "bookInfoCanary100pcReviewsV3"
	bookInfoInjectDelayForRatingsForJason    = "bookInfoInjectDelayForRatingsForJason"
	bookInfoInjectHTTPAbortToRatingsForJason = "bookInfoInjectHTTPAbortToRatingsForJason"
	bookInfoProductPageCircuitBreaking       = "bookInfoProductPageCircuitBreaking"
)

var supportedOps = map[string]supportedOperation{
	installIstioCommand: {
		name: "Latest Istio without mTLS",
		// templateName: "install_istio.tmpl",
		opType: meshes.OpCategory_INSTALL,
	},
	installmTLSIstioCommand: {
		name:   "Latest Istio with mTLS",
		opType: meshes.OpCategory_INSTALL,
	},
	installBookInfoCommand: {
		name: "Book Info Application",
		// templateName: "install_istio.tmpl",
		opType: meshes.OpCategory_SAMPLE_APPLICATION,
	},
	runVet: {
		name:   "Run istio-vet",
		opType: meshes.OpCategory_VALIDATE,
		// templateName: "istio_vet.tmpl",
		// appLabel:     "istio-vet",
		// returnLogs:   true,
	},
	cbCommand: {
		name:         "Configure circuit breaker with only one connection",
		opType:       meshes.OpCategory_CONFIGURE,
		templateName: "circuit_breaking.tmpl",
	},
	bookInfoDefaultDestinationRules: {
		name:   "Default Book info destination rules (defines subsets)",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoRouteToV1AllServices: {
		name:   "Route traffic to V1 of all Book info services",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoRouteToReviewsV2ForJason: {
		name:   "Route traffic to V2 of Book info reviews service for user Jason",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoCanary50pcReviewsV3: {
		name:   "Route 50% of the traffic to Book info reviews V3",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoCanary100pcReviewsV3: {
		name:   "Route 100% of the traffic to Book info reviews V3",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoInjectDelayForRatingsForJason: {
		name:   "Inject a 7s delay in the traffic to Book info ratings service for user Jason",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoInjectHTTPAbortToRatingsForJason: {
		name:   "Inject HTTP abort to Book info ratings service for user Jason",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoProductPageCircuitBreaking: {
		name:         "Configure circuit breaking with max 1 request per connection and max 1 pending request to Book info productpage service",
		opType:       meshes.OpCategory_CONFIGURE,
		templateName: "book_info_product_page_circuit_breaking.tmpl",
	},
	installSMI: {
		name:   "Service Mesh Interface (SMI) Istio Adapter",
		opType: meshes.OpCategory_INSTALL,
	},
	installHTTPBin: {
		name:         "HTTPbin Application",
		templateName: "httpbin.yaml",
		opType:       meshes.OpCategory_SAMPLE_APPLICATION,
	},
	customOpCommand: {
		name:   "Custom YAML",
		opType: meshes.OpCategory_CUSTOM,
	},
	googleMSSampleApplication: {
		name:   "Hipster Shop Application",
		opType: meshes.OpCategory_SAMPLE_APPLICATION,
	},
}
