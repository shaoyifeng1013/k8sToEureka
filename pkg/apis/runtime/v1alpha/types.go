/*
Copyright 2017 The Kubernetes Authors.

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

package v1alpha

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MeshServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec MeshServerSpec `json:"spec"`
}

type MeshServerSpec struct {
	Server `json:"server"`
	Zone   string `json:"zone"`
}

type Server struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	IpAddr string `json:"ipAddr"`
}

func (ms MeshServer) IsValid() bool {
	if ms.Spec.Host == "" || ms.Spec.IpAddr == "" || ms.Spec.Port == 0 {
		klog.Info(ms.Name, "is invalid")
		return false
	}
	return true
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MeshServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []MeshServer `json:"items"`
}
