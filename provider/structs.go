package provider

type CommandOptions struct {
	AwsRegion string
	Path      string
	Output    string
	Overwrite bool
	Combined  bool
	Backup    bool
	DryRun    bool
	Version   bool
}

// Azure
type subs struct {
	Name string `json:"name,omitempty"`
	Id   string `json:"id,omitempty"`
}

// Standard of Cluster Info Output
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
type LocalConfig struct {
	Authinfo User     `json:"authinfo,omitempty"`
	Context  Context  `json:"context,omitempty"`
	Cluster  CCluster `json:"cluster,omitempty"`
}

type AllConfig struct {
	auth     []Users
	context  []Contexts
	clusters []Clusters
}

// merged config struct
type Config struct {
	APIVersion     string      `yaml:"apiVersion,omitempty"`
	Clusters       []Clusters  `yaml:"clusters,omitempty"`
	Contexts       []Contexts  `yaml:"contexts,omitempty"`
	CurrentContext string      `yaml:"current-context,omitempty"`
	Kind           string      `yaml:"kind,omitempty"`
	Preferences    Preferences `yaml:"preferences,omitempty"`
	Users          []Users     `yaml:"users,omitempty"`
}

type Clusters struct {
	Cluster CCluster `yaml:"cluster,omitempty"`
	Name    string   `yaml:"name,omitempty"`
}

type CCluster struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data,omitempty"`
	Server                   string `yaml:"server,omitempty"`
}

type Context struct {
	Cluster string `yaml:"cluster,omitempty"`
	User    string `yaml:"user,omitempty"`
}

type Contexts struct {
	Context Context `yaml:"context,omitempty"`
	Name    string  `yaml:"name,omitempty"`
}

type Preferences struct {
}

type Exec struct {
	APIVersion         string      `yaml:"apiVersion,omitempty"`
	Args               []string    `yaml:"args,omitempty"`
	Command            string      `yaml:"command,omitempty"`
	Env                interface{} `yaml:"env,omitempty"`
	ProvideClusterInfo bool        `yaml:"provideClusterInfo,omitempty"`
}

type User struct {
	Exec                  Exec   `yaml:"exec,omitempty"`
	ClientCertificateData string `yaml:"client-certificate-data,omitempty"`
	ClientKeyData         string `yaml:"client-key-data,omitempty"`
	Token                 string `yaml:"token,omitempty"`
}

type Users struct {
	Name string `yaml:"name,omitempty"`
	User User   `yaml:"user,omitempty"`
}
