package provider

// type account struct {
// 	AwsProfile string `json:"accountName,omitempty"`
// 	Report     report `json:"accounts,omitempty"`
// }

// type report struct {
// 	K8s        []region `json:"regions,omitempty"`
// 	TotalCount int      `json:"totalCount,omitempty"`
// }

// type region struct {
// 	Region     string    `json:"region,omitempty"`
// 	Clusters   []cluster `json:"clusters,omitempty"`
// 	TotalCount int       `json:"totalCount,omitempty"`
// }

// type cluster struct {
// 	ClusterName    string `json:"clusterName,omitempty"`
// 	CurrentVersion string `json:"currentVersion,omitempty"`
// 	LatestVersion  string `json:"latestVersion,omitempty"`
// }

type resource struct {
	Id            string `json:"id,omitempty"`
	Location      string `json:"location,omitempty"`
	Name          string `json:"name,omitempty"`
	Type          string `json:"type,omitempty"`
	Version       string `json:"version,omitempty"`
	LatestVersion string `json:"latest_version,omitempty"`
}

type subs struct {
	Name string `json:"name,omitempty"`
	Id   string `json:"id,omitempty"`
}

type subAks struct {
	Resources  []resource
	TotalCount int `json:"totalCount,omitempty"`
}
type AzAKSList struct {
	Subscription string `json:"subscription,omitempty"`
	Aks          subAks
}

// Standard of Cluster Output

type Cluster struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Latest  string `json:"latest,omitempty"`
	Region  string `json:"region,omitempty"`
}

type Account struct {
	Name     string    `json:"name,omitempty"`
	Clusters []Cluster `json:"clusters,omitempty"`
}

type Provider struct {
	Provider string    `json:"provider,omitempty"`
	Accounts []Account `json:"accounts,omitempty"`
}

type FullK8S struct {
	Providers []Provider `json:"providers,omitempty"`
}
