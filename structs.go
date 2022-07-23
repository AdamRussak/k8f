package main

type account struct {
	AwsProfile string `json:"accountName,omitempty"`
	Report     report `json:"accounts,omitempty"`
}
type report struct {
	K8s        []region `json:"regions,omitempty"`
	TotalCount int      `json:"totalCount,omitempty"`
}
type region struct {
	Region     string    `json:"region,omitempty"`
	Clusters   []cluster `json:"clusters,omitempty"`
	TotalCount int       `json:"totalCount,omitempty"`
}
type cluster struct {
	ClusterName    string `json:"clusterName,omitempty"`
	CurrentVersion string `json:"currentVersion,omitempty"`
	LatestVersion  string `json:"latestVersion,omitempty"`
}
