package provider

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
