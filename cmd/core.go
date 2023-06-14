package cmd

import (
	"errors"
	"fmt"
	"k8f/core"

	"github.com/spf13/cobra"
)

func argValidator(cmd *cobra.Command, args []string) error {
	var err error

	err = checkArgsCount(args)
	core.OnErrorFail(err, "validation failed")
	err = providerValidator(args)
	core.OnErrorFail(err, "validation failed")
	err = uiSelectValidator(args)
	core.OnErrorFail(err, "validation failed")
	return err
}

// check amounts of args in the command
func checkArgsCount(args []string) error {
	if len(args) < 1 {
		return errors.New("requires cloud provider")
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
