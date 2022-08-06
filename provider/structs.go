package provider

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
	Location  *string            `json:"location,omitempty"`
	ManagedBy *string            `json:"managedBy,omitempty"`
	Tags      map[string]*string `json:"tags,omitempty"`
	ID        *string            `json:"id,omitempty" azure:"ro"`
	Name      *string            `json:"name,omitempty" azure:"ro"`
	Type      *string            `json:"type,omitempty" azure:"ro"`
}

type resource struct {
	Id            string `json:"id,omitempty"`
	Location      string `json:"location,omitempty"`
	Name          string `json:"name,omitempty"`
	Type          string `json:"type,omitempty"`
	Version       string `json:"version,omitempty"`
	LatestVersion string `json:"latest_version,omitempty"`
}

type rgAndResouce struct {
	RGName    string     `json:"resource_group_name,omitempty"`
	Resources []resource `json:"resources,omitempty"`
}

type subs struct {
	Name string `json:"name,omitempty"`
	Id   string `json:"id,omitempty"`
}
