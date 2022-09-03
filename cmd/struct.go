package cmd

type FlagsOptions struct {
	Path          string `json:"path,omitempty"`
	Output        string `json:"output,omitempty"`
	Overwrite     bool   `json:"overwrite,omitempty"`
	Backup        bool   `json:"backup,omitempty"`
	DryRun        bool   `json:"dry-run,omitempty"`
	AwsAuth       bool   `json:"aws_auth,omitempty"`
	AwsAssumeRole bool   `json:"aws_assume_role,omitempty"`
	AwsRoleString string `json:"aws_role_string,omitempty"`
	AwsEnvProfile bool   `json:"aws_profile,omitempty"`
}
