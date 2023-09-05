package cmd

type FlagsOptions struct {
	Path            string `json:"path,omitempty"`
	ListPath        string `json:"list_path,omitempty"`
	Output          string `json:"output,omitempty"`
	ListOutput      string `json:"list_output,omitempty"`
	Overwrite       bool   `json:"overwrite,omitempty"`
	Backup          bool   `json:"backup,omitempty"`
	Merge           bool   `json:"merge,omitempty"`
	ForceMerge      bool   `json:"force_merge,omitempty"`
	UiSize          int    `json:"uiSize,omitempty"`
	DryRun          bool   `json:"dry-run,omitempty"`
	AwsAuth         bool   `json:"aws_auth,omitempty"`
	AwsAssumeRole   bool   `json:"aws_assume_role,omitempty"`
	AwsRoleString   string `json:"aws_role_string,omitempty"`
	AwsEnvProfile   bool   `json:"aws_profile,omitempty"`
	AwsClusterName  bool   `json:"aws_cluster_name,omitempty"`
	ProfileSelector bool   `json:"profile_selector,omitempty"`
	SaveOutput      bool   `json:"save_output,omitempty"`
	Validate        bool   `json:"validate,omitempty"`
}
