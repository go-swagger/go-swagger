// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"

	"github.com/go-swagger/go-swagger/fixtures/bugs/2969/codegen/models"
	"github.com/spf13/cobra"
)

// Schema cli for RunProfileRescheduleOptions

// register flags to command
func registerModelRunProfileRescheduleOptionsFlags(depth int, cmdPrefix string, cmd *cobra.Command) error {

	if err := registerRunProfileRescheduleOptionsPropMethod(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerRunProfileRescheduleOptionsPropRunProfileID(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	return nil
}

func registerRunProfileRescheduleOptionsPropMethod(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	methodDescription := `Enum: ["DEFAULT","DUPLICATE_RUN_PROFILES","CLUSTER_RUN_PROFILE"]. Run profile reschedule method`

	var methodFlagName string
	if cmdPrefix == "" {
		methodFlagName = "method"
	} else {
		methodFlagName = fmt.Sprintf("%v.method", cmdPrefix)
	}

	var methodFlagDefault string

	_ = cmd.PersistentFlags().String(methodFlagName, methodFlagDefault, methodDescription)

	if err := cmd.RegisterFlagCompletionFunc(methodFlagName,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			var res []string
			if err := json.Unmarshal([]byte(`["DEFAULT","DUPLICATE_RUN_PROFILES","CLUSTER_RUN_PROFILE"]`), &res); err != nil {
				panic(err)
			}
			return res, cobra.ShellCompDirectiveDefault
		}); err != nil {
		return err
	}

	return nil
}

func registerRunProfileRescheduleOptionsPropRunProfileID(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	runProfileIdDescription := `Run profile to be used from cluster`

	var runProfileIdFlagName string
	if cmdPrefix == "" {
		runProfileIdFlagName = "runProfileId"
	} else {
		runProfileIdFlagName = fmt.Sprintf("%v.runProfileId", cmdPrefix)
	}

	var runProfileIdFlagDefault string

	_ = cmd.PersistentFlags().String(runProfileIdFlagName, runProfileIdFlagDefault, runProfileIdDescription)

	return nil
}

// retrieve flags from commands, and set value in model. Return true if any flag is passed by user to fill model field.
func retrieveModelRunProfileRescheduleOptionsFlags(depth int, m *models.RunProfileRescheduleOptions, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false

	err, methodAdded := retrieveRunProfileRescheduleOptionsPropMethodFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || methodAdded

	err, runProfileIdAdded := retrieveRunProfileRescheduleOptionsPropRunProfileIDFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || runProfileIdAdded

	return nil, retAdded
}

func retrieveRunProfileRescheduleOptionsPropMethodFlags(depth int, m *models.RunProfileRescheduleOptions, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	methodFlagName := fmt.Sprintf("%v.method", cmdPrefix)
	if cmd.Flags().Changed(methodFlagName) {

		var methodFlagName string
		if cmdPrefix == "" {
			methodFlagName = "method"
		} else {
			methodFlagName = fmt.Sprintf("%v.method", cmdPrefix)
		}

		methodFlagValue, err := cmd.Flags().GetString(methodFlagName)
		if err != nil {
			return err, false
		}
		m.Method = methodFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveRunProfileRescheduleOptionsPropRunProfileIDFlags(depth int, m *models.RunProfileRescheduleOptions, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	runProfileIdFlagName := fmt.Sprintf("%v.runProfileId", cmdPrefix)
	if cmd.Flags().Changed(runProfileIdFlagName) {

		var runProfileIdFlagName string
		if cmdPrefix == "" {
			runProfileIdFlagName = "runProfileId"
		} else {
			runProfileIdFlagName = fmt.Sprintf("%v.runProfileId", cmdPrefix)
		}

		runProfileIdFlagValue, err := cmd.Flags().GetString(runProfileIdFlagName)
		if err != nil {
			return err, false
		}
		m.RunProfileID = runProfileIdFlagValue

		retAdded = true
	}

	return nil, retAdded
}
