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
type ResourceGroup struct {
	// REQUIRED; The location of the resource group. It cannot be changed after the resource group has been created. It must be
	// one of the supported Azure locations.
	Location *string `json:"location,omitempty"`

	// The ID of the resource that manages this resource group.
	ManagedBy *string `json:"managedBy,omitempty"`

	// The tags attached to the resource group.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; The ID of the resource group.
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource group.
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource group.
	Type *string `json:"type,omitempty" azure:"ro"`
}
