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

type supportedOperation struct {
	// a unique identifier
	key string
	// a friendly name
	name string
	// the template file name
	templateName string

	resetOp bool
}

var supportedOps = []supportedOperation{
	{
		key:          "istio_1",
		name:         "Shift All traffic to version V1 of all the services",
		templateName: "shift_all_traffic_to_v1_of_all_services.tmpl",
	},
	{
		key:          "istio_2",
		name:         "Shift logged in user traffic to version V2 of reviews service",
		templateName: "shift_user_traffic_to_v2_of_reviews.tmpl",
	},
	{
		key:          "istio_3",
		name:         "Inject a HTTP abort for ratings service for the logged in user",
		templateName: "inject_abort_for_ratings_service_for_user.tmpl",
	},
	{
		key:          "istio_4",
		name:         "Inject a HTTP delay for ratings service for the logged in user",
		templateName: "inject_delay_for_ratings_service_for_user.tmpl",
	},
	{
		key:          "istio_5",
		name:         "Shift 50 percent of traffic to version V3 of reviews service",
		templateName: "shift_50_percent_of_traffic_to_v3_of_reviews.tmpl",
	},
	{
		key:          "istio_6",
		name:         "Shift 100 percent of traffic to version V3 of reviews service",
		templateName: "shift_all_traffic_to_v3_of_reviews.tmpl",
	},
	{
		key:     "istio_7",
		name:    "Reset all applied rules",
		resetOp: true,
	},
}
