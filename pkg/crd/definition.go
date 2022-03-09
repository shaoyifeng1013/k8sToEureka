package crd

import "k8s.io/apimachinery/pkg/runtime/schema"

type Build struct {
	ApiVersion           string
	Kind                 string
	GroupVersionResource schema.GroupVersionResource
}

type MeshServer struct {
	Server `json:"server"`
	Zone   string `json:"zone"`
}

type Server struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

var MeshServers = Build{
	ApiVersion: "runtime.syf.io/v1alpha",
	Kind:       "MeshServer",
	GroupVersionResource: schema.GroupVersionResource{
		Group:    "runtime.syf.io",
		Version:  "v1alpha",
		Resource: "meshservers",
	},
}
