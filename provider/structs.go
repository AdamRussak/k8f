package provider

import clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

// Azure
type subs struct {
	Name string `json:"name,omitempty"`
	Id   string `json:"id,omitempty"`
}

// Standard of Cluster Output

type Cluster struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Latest  string `json:"latest,omitempty"`
	Region  string `json:"region,omitempty"`
	Id      string `json:"id,omitempty"`
}

type Account struct {
	Name       string    `json:"name,omitempty"`
	Clusters   []Cluster `json:"clusters,omitempty"`
	TotalCount int       `json:"totalCount,omitempty"`
}

type Provider struct {
	Provider   string    `json:"provider,omitempty"`
	Accounts   []Account `json:"accounts,omitempty"`
	TotalCount int       `json:"totalCount,omitempty"`
}

// AWS Kubeconfig
type clusterConfigInfo struct {
	Cluster     string      `json:"cluster,omitempty"`
	Version     string      `json:"version,omitempty"`
	LocalConfig LocalConfig `json:"localConfig,omitempty"`
}

type LocalConfig struct {
	Authinfo *clientcmdapi.AuthInfo `json:"authinfo,omitempty"`
	Context  *clientcmdapi.Context  `json:"context,omitempty"`
	Cluster  *clientcmdapi.Cluster  `json:"cluster,omitempty"`
}

type AllConfig struct {
	Authinfos map[string]*clientcmdapi.AuthInfo `json:"authinfos,omitempty"`
	Contexts  map[string]*clientcmdapi.Context  `json:"contexts,omitempty"`
	Clusters  map[string]*clientcmdapi.Cluster  `json:"clusters,omitempty"`
}
