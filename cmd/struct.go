package cmd

type FlagsOptions struct {
	Path      string `json:"path,omitempty"`
	Output    string `json:"output,omitempty"`
	Overwrite bool   `json:"overwrite,omitempty"`
	Backup    bool   `json:"backup,omitempty"`
	DryRun    bool   `json:"dry-run,omitempty"`
	Version   bool   `json:"version,omitempty"`
}
type loggingStruct struct {
	Provider string
}
