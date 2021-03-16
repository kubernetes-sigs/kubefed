/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package status

import (
	"testing"

	apiv1 "k8s.io/api/core/v1"
)

func TestGenericPropagationStatusUpdateChanged(t *testing.T) {
	testCases := map[string]struct {
		generation        int64
		reason            AggregateReason
		statusMap         PropagationStatusMap
		resourceStatusMap map[string]interface{}
		resourcesUpdated  bool
		expectedChanged   bool
		resourceStatusCollection bool
	}{
		"Cluster not propagated indicates changed with status collected enabled": {
			statusMap: PropagationStatusMap{
				"cluster1": ClusterNotReady,
			},
			resourceStatusMap: map[string]interface{}{
				"cluster1": map[string]interface{}{},
			},
			reason: AggregateSuccess,
			resourcesUpdated: false,
			resourceStatusCollection: true,
			expectedChanged:  true,
		},
		"Cluster not propagated indicates changed with status collected disabled": {
			statusMap: PropagationStatusMap{
				"cluster1": ClusterNotReady,
			},
			resourceStatusMap: map[string]interface{}{
				"cluster1": map[string]interface{}{},
			},
			reason: AggregateSuccess,
			resourcesUpdated: false,
			resourceStatusCollection: false,
			expectedChanged:  true,
		},
		"Cluster status not retrieved indicates changed with status collected enabled": {
			statusMap: PropagationStatusMap{
				"cluster1": ClusterPropagationOK,
				"cluster2": ClusterPropagationOK,
			},
			resourceStatusMap: map[string]interface{}{
				"cluster1": map[string]interface{}{},
			},
			reason: AggregateSuccess,
			resourcesUpdated: false,
			resourceStatusCollection: true,
			expectedChanged:  true,
		},
		"Cluster status not retrieved indicates changed with status collected disabled": {
			statusMap: PropagationStatusMap{
				"cluster1": ClusterPropagationOK,
				"cluster2": ClusterPropagationOK,
			},
			resourceStatusMap: map[string]interface{}{
				"cluster1": map[string]interface{}{},
			},
			reason: AggregateSuccess,
			resourcesUpdated: false,
			resourceStatusCollection: false,
			expectedChanged:  true,
		},
		//"No change in clusters indicates unchanged": {
		//	statusMap: PropagationStatusMap{
		//		"cluster1": ClusterPropagationOK,
		//	},
		//	resourceStatusMap: map[string]interface{}{
		//		"cluster1": map[string]interface{}{},
		//	},
		//	resourcesUpdated: false,
		//	resourceStatusCollection: true,
		//	expectedChanged:  true,
		//},
		//"No change in clusters with update indicates changed": {
		//	statusMap: PropagationStatusMap{
		//		"cluster1": ClusterPropagationOK,
		//	},
		//	resourceStatusMap: map[string]interface{}{
		//		"ready": false,
		//		"stage": "absent",
		//	},
		//	resourcesUpdated: true,
		//	resourceStatusCollection: true,
		//	expectedChanged:  true,
		//},
		//"Change in clusters indicates changed": {
		//	statusMap: PropagationStatusMap{
		//		"cluster1": ClusterPropagationOK,
		//	},
		//	resourceStatusMap: map[string]interface{}{
		//		"ready": true,
		//		"stage": "deployed",
		//	},
		//	resourceStatusCollection: true,
		//	expectedChanged: true,
		//},
		//"Transition indicates changed with remote status collection enabled": {
		//	reason:          NamespaceNotFederated,
		//	resourceStatusCollection: true,
		//	expectedChanged: true,
		//},
		//"Transition indicates changed with remote status collection disabled": {
		//	reason:          NamespaceNotFederated,
		//	resourceStatusCollection: false,
		//	expectedChanged: true,
		//},
		//"Changed generation indicates changed with remote status collection enabled": {
		//	generation:      1,
		//	resourceStatusCollection: true,
		//	expectedChanged: true,
		//},
		//"Changed generation indicates changed with remote status collection disabled": {
		//	generation:      1,
		//	resourceStatusCollection: false,
		//	expectedChanged: true,
		//},
	}
	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			fedStatus := &GenericFederatedStatus{
				Clusters: []GenericClusterStatus{
					{
						Name: "cluster1",
					},
				},
				Conditions: []*GenericCondition{
					{
						Type:   PropagationConditionType,
						Status: apiv1.ConditionTrue,
					},
				},
			}
			collectedStatus := CollectedPropagationStatus{
				StatusMap:        tc.statusMap,
				ResourcesUpdated: tc.resourcesUpdated,
			}
			collectedResourceStatus := CollectedResourceStatus{
				StatusMap:        tc.resourceStatusMap,
				ResourcesUpdated: tc.resourcesUpdated,
			}
			changed := fedStatus.update(tc.generation, tc.reason, collectedStatus, collectedResourceStatus, tc.resourceStatusCollection)
			if tc.expectedChanged != changed {
				t.Fatalf("Expected changed to be %v, got %v", tc.expectedChanged, changed)
			}
		})
	}
}

