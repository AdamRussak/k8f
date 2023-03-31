package cmd

import (
	"k8f/core"
	"k8f/provider"
)

func newCommandStruct(o FlagsOptions, args []string) provider.CommandOptions {
	commandOptions := provider.CommandOptions{
		AwsRegion:     AwsRegion,
		ForceMerge:    o.ForceMerge,
		UiSize:        o.UiSize,
		Path:          o.Path,
		Output:        o.Output,
		Overwrite:     o.Overwrite,
		Combined:      core.BoolCombine(args[0], supportedProvider),
		Merge:         o.Merge,
		Backup:        o.Backup,
		DryRun:        o.DryRun,
		AwsAuth:       o.AwsAuth,
		AwsRoleString: o.AwsRoleString,
		AwsEnvProfile: o.AwsEnvProfile,
	}

	return commandOptions
}
