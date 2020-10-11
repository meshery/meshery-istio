// Copyright 2020, Layer5 Authors
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
	customOpCommand = "custom"
	runVet          = "istio_vet"

	// Install istio
	installv173IstioCommand     = "istio_install_v173"
	installv173IstioCommandTls  = "istio_install_v173_tls"
	installOperatorIstioCommand = "istio_install_operator"
	installmTLSIstioCommand     = "istio_mtls_install"
	verifyInstallation          = "verify_installation" // requires
	installAddons               = "install_addons"
	injectLabels                = "inject_labels"

	// Bookinfo
	installBookInfoCommand                   = "install_book_info"
	cbCommand                                = "cb1"
	googleMSSampleApplication                = "google_microservices_demo_application"
	bookInfoDefaultDestinationRules          = "bookInfoDefaultDestinationRules"
	bookInfoRouteToV1AllServices             = "bookInfoRouteToV1AllServices"
	bookInfoRouteToReviewsV2ForJason         = "bookInfoRouteToReviewsV2ForJason"
	bookInfoCanary50pcReviewsV3              = "bookInfoCanary50pcReviewsV3"
	bookInfoCanary100pcReviewsV3             = "bookInfoCanary100pcReviewsV3"
	bookInfoInjectDelayForRatingsForJason    = "bookInfoInjectDelayForRatingsForJason"
	bookInfoInjectHTTPAbortToRatingsForJason = "bookInfoInjectHTTPAbortToRatingsForJason"
	bookInfoProductPageCircuitBreaking       = "bookInfoProductPageCircuitBreaking"

	// HTTPbin
	installHttpbinCommandV1 = "install_http_binv1"
	installHttpbinCommandV2 = "install_http_binv2"

	// Labs
	enablePrometheus = "enable_prometheus"
	enableGrafana    = "enable_grafana"
	enableKiali      = "enable_kiali"
	denyAllPolicy    = "deny_all_policy"
	strictMtls       = "strict_mtls"
	mutualMtls       = "mutual_mtls"
	disableMtls      = "disable_mtls"

	bookInfoRouteV1ForUser                 = "bookinfo_route_v1_for_user"
	bookInfoMirrorTrafficToV2              = "bookinfo_mirror_traffic_to_v2"
	bookInfoRetrySpecForReviews            = "bookinfo_retry_spec_for_reviews"
	bookInfoCanaryDeploy20V3               = "bookinfo_canary_deploy_20_v3"
	bookInfoCanaryDeploy80V3               = "bookinfo_canary_deploy_80_v3"
	bookInfoCanaryDeploy100V3              = "bookinfo_canary_deploy_100_v3"
	bookInfoInjectDelayFaultRatings        = "bookinfo_inject_delay_fault_ratings"
	bookInfoInjectDelayFaultReviews        = "bookinfo_inject_delay_fault_reviews"
	bookInfoConfigureConnectionPoolOutlier = "bookinfo_configure_connection_pool_outlier"
	bookInfoAllowGet                       = "bookinfo_allow_get"
	bookInfoAllowReviewsForUser            = "bookinfo_allow_reviews_for_user"

	// SMI conformance test
	smiConformanceCommand = "smiConformanceTest"
	installSMI            = "install_smi"
)

var supportedOps = map[string]supportedOperation{
	installv173IstioCommand: {
		name:   "Istio 1.7.3",
		opType: meshes.OpCategory_INSTALL,
	},
	installv173IstioCommandTls: {
		name:   "Istio 1.7.3 with mTLS",
		opType: meshes.OpCategory_INSTALL,
	},
	installOperatorIstioCommand: {
		name:   "Istio with Operator",
		opType: meshes.OpCategory_INSTALL,
	},
	installmTLSIstioCommand: {
		name:   "Istio 1.5.1 with mTLS",
		opType: meshes.OpCategory_INSTALL,
	},
	installBookInfoCommand: {
		name: "BookInfo Application",
		// templateName: "install_istio.tmpl",
		opType: meshes.OpCategory_SAMPLE_APPLICATION,
	},
	runVet: {
		name:   "Check configuration",
		opType: meshes.OpCategory_VALIDATE,
		// templateName: "istio_vet.tmpl",
		// appLabel:     "istio-vet",
		// returnLogs:   true,
	},
	cbCommand: {
		name:         "httpbin: Configure circuit breaker with only one connection",
		opType:       meshes.OpCategory_CONFIGURE,
		templateName: "circuit_breaking.tmpl",
	},
	bookInfoDefaultDestinationRules: {
		name:   "BookInfo: Default BookInfo destination rules (defines subsets)",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoRouteToV1AllServices: {
		name:   "BookInfo: Route traffic to V1 of all BookInfo services",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoRouteToReviewsV2ForJason: {
		name:   "BookInfo: Route traffic to V2 of BookInfo reviews service for user Jason",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoCanary50pcReviewsV3: {
		name:   "BookInfo: Route 50% of the traffic to BookInfo reviews V3",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoCanary100pcReviewsV3: {
		name:   "BookInfo: Route 100% of the traffic to BookInfo reviews V3",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoInjectDelayForRatingsForJason: {
		name:   "BookInfo: Inject a 7s delay in the traffic to BookInfo ratings service for user Jason",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoInjectHTTPAbortToRatingsForJason: {
		name:   "BookInfo: Inject HTTP abort to BookInfo ratings service for user Jason",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoProductPageCircuitBreaking: {
		name:         "BookInfo: Configure circuit breaking with max 1 request per connection and max 1 pending request to BookInfo productpage service",
		opType:       meshes.OpCategory_CONFIGURE,
		templateName: "book_info_product_page_circuit_breaking.tmpl",
	},
	installSMI: {
		name:   "Service Mesh Interface (SMI) Istio Adapter",
		opType: meshes.OpCategory_INSTALL,
	},
	installHttpbinCommandV1: {
		name:         "httpbin Application V1",
		templateName: "v1",
		opType:       meshes.OpCategory_SAMPLE_APPLICATION,
	},
	installHttpbinCommandV2: {
		name:         "httpbin Application V2 (Needs V1 installed)",
		templateName: "v2",
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
	smiConformanceCommand: {
		name:   "Run SMI conformance test",
		opType: meshes.OpCategory_VALIDATE,
	},
	enablePrometheus: {
		name:   "Enable Prometheus monitoring",
		opType: meshes.OpCategory_INSTALL,
	},
	enableGrafana: {
		name:   "Enable Grafana dashboard",
		opType: meshes.OpCategory_INSTALL,
	},
	enableKiali: {
		name:   "Enable Prometheus dashboard",
		opType: meshes.OpCategory_INSTALL,
	},
	denyAllPolicy: {
		name:   "Deny-All policy on the namespace",
		opType: meshes.OpCategory_CONFIGURE,
	},
	strictMtls: {
		name:   "Strict Mtls policy",
		opType: meshes.OpCategory_CONFIGURE,
	},
	mutualMtls: {
		name:   "Mutual Mtls policy",
		opType: meshes.OpCategory_CONFIGURE,
	},
	disableMtls: {
		name:   "Disable Mtls policy",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoRouteV1ForUser: {
		name:   "Configure bookinfo page to version v1",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoMirrorTrafficToV2: {
		name:   "Configure bookinfo page mirror traffic from v1 to v2",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoRetrySpecForReviews: {
		name:   "Configure bookinfo page to retry for reviews application",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoCanaryDeploy20V3: {
		name:   "Configure bookinfo to forward 20 percent traffic to v2",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoCanaryDeploy80V3: {
		name:   "Configure bookinfo to forward 80 percent traffic to v2",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoCanaryDeploy100V3: {
		name:   "Configure bookinfo to forward 100 percent traffic to v2",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoInjectDelayFaultRatings: {
		name:   "Configure bookinfo to add delay to ratings application",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoInjectDelayFaultReviews: {
		name:   "Configure bookinfo to add delay to ratings application",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoConfigureConnectionPoolOutlier: {
		name:   "Configure bookinfo for connection pool limits and outlier detection",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoAllowGet: {
		name:   "Configure bookinfo to allow only GET requests",
		opType: meshes.OpCategory_CONFIGURE,
	},
	bookInfoAllowReviewsForUser: {
		name:   "Configure bookinfo to allow reviews only for user",
		opType: meshes.OpCategory_CONFIGURE,
	},
}
