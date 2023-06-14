package cmd

import (
	"errors"
	"fmt"
	"k8f/core"
	"k8f/provider"

	"github.com/spf13/cobra"
)

func newCommandStruct(o FlagsOptions, args []string) provider.CommandOptions {
	commandOptions := provider.CommandOptions{
		AwsRegion:       AwsRegion,
		ForceMerge:      o.ForceMerge,
		UiSize:          o.UiSize,
		Path:            o.Path,
		Output:          o.Output,
		Overwrite:       o.Overwrite,
		Combined:        core.BoolCombine(args[0], supportedProvider),
		Merge:           o.Merge,
		Backup:          o.Backup,
		DryRun:          o.DryRun,
		AwsAuth:         o.AwsAuth,
		AwsRoleString:   o.AwsRoleString,
		AwsEnvProfile:   o.AwsEnvProfile,
		ProfileSelector: o.ProfileSelector,
	}

	return commandOptions
}

func argValidator(cmd *cobra.Command, args []string) error {
	var err error

	err = checkArgsCount(args)
	core.FailOnError(err, validationError)
	err = providerValidator(args)
	core.FailOnError(err, validationError)
	err = uiSelectValidator(args)
	core.FailOnError(err, validationError)
	return err
}

// check amounts of args in the command
func checkArgsCount(args []string) error {
	if len(args) < 1 {
		return errors.New(providerError)
	}
	return nil
}

// check the args for supported Provider
func providerValidator(args []string) error {
	argouments = append(argouments, supportedProvider...)
	if len(args) > 0 && len(args) <= len(argouments) {
		for a := range args {
			if !core.IfXinY(args[a], argouments) {
				return fmt.Errorf("invalid cloud provider specified - %s", args[a])
			}
		}
	}
	return nil
}

// check if aws UI select was set
func uiSelectValidator(args []string) error {
	if o.ProfileSelector && args[0] != "aws" {
		return fmt.Errorf("profile selector supports only AWS and is not supporting - %s", args[0])
	}
	return nil
}
