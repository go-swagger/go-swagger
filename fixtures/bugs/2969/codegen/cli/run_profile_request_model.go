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

// Schema cli for RunProfileRequest

// register flags to command
func registerModelRunProfileRequestFlags(depth int, cmdPrefix string, cmd *cobra.Command) error {

	if err := registerRunProfileRequestPropDescription(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerRunProfileRequestPropJvmArguments(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerRunProfileRequestPropName(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerRunProfileRequestPropRuntimeID(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerRunProfileRequestPropType(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	return nil
}

func registerRunProfileRequestPropDescription(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	descriptionDescription := `Description`

	var descriptionFlagName string
	if cmdPrefix == "" {
		descriptionFlagName = "description"
	} else {
		descriptionFlagName = fmt.Sprintf("%v.description", cmdPrefix)
	}

	var descriptionFlagDefault string

	_ = cmd.PersistentFlags().String(descriptionFlagName, descriptionFlagDefault, descriptionDescription)

	return nil
}

func registerRunProfileRequestPropJvmArguments(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	// warning: jvmArguments []string array type is not supported by go-swagger cli yet

	return nil
}

func registerRunProfileRequestPropName(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	nameDescription := `Required. Run profile name`

	var nameFlagName string
	if cmdPrefix == "" {
		nameFlagName = "name"
	} else {
		nameFlagName = fmt.Sprintf("%v.name", cmdPrefix)
	}

	var nameFlagDefault string

	_ = cmd.PersistentFlags().String(nameFlagName, nameFlagDefault, nameDescription)

	return nil
}

func registerRunProfileRequestPropRuntimeID(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	runtimeIdDescription := `Required. Runtime id`

	var runtimeIdFlagName string
	if cmdPrefix == "" {
		runtimeIdFlagName = "runtimeId"
	} else {
		runtimeIdFlagName = fmt.Sprintf("%v.runtimeId", cmdPrefix)
	}

	var runtimeIdFlagDefault string

	_ = cmd.PersistentFlags().String(runtimeIdFlagName, runtimeIdFlagDefault, runtimeIdDescription)

	return nil
}

func registerRunProfileRequestPropType(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	typeDescription := `Enum: ["JOB_SERVER","MICROSERVICE","TALEND_RUNTIME"]. Required. Run profile type`

	var typeFlagName string
	if cmdPrefix == "" {
		typeFlagName = "type"
	} else {
		typeFlagName = fmt.Sprintf("%v.type", cmdPrefix)
	}

	var typeFlagDefault string

	_ = cmd.PersistentFlags().String(typeFlagName, typeFlagDefault, typeDescription)

	if err := cmd.RegisterFlagCompletionFunc(typeFlagName,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			var res []string
			if err := json.Unmarshal([]byte(`["JOB_SERVER","MICROSERVICE","TALEND_RUNTIME"]`), &res); err != nil {
				panic(err)
			}
			return res, cobra.ShellCompDirectiveDefault
		}); err != nil {
		return err
	}

	return nil
}

// retrieve flags from commands, and set value in model. Return true if any flag is passed by user to fill model field.
func retrieveModelRunProfileRequestFlags(depth int, m *models.RunProfileRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false

	err, descriptionAdded := retrieveRunProfileRequestPropDescriptionFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || descriptionAdded

	err, jvmArgumentsAdded := retrieveRunProfileRequestPropJvmArgumentsFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || jvmArgumentsAdded

	err, nameAdded := retrieveRunProfileRequestPropNameFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || nameAdded

	err, runtimeIdAdded := retrieveRunProfileRequestPropRuntimeIDFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || runtimeIdAdded

	err, typeAdded := retrieveRunProfileRequestPropTypeFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || typeAdded

	return nil, retAdded
}

func retrieveRunProfileRequestPropDescriptionFlags(depth int, m *models.RunProfileRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	descriptionFlagName := fmt.Sprintf("%v.description", cmdPrefix)
	if cmd.Flags().Changed(descriptionFlagName) {

		var descriptionFlagName string
		if cmdPrefix == "" {
			descriptionFlagName = "description"
		} else {
			descriptionFlagName = fmt.Sprintf("%v.description", cmdPrefix)
		}

		descriptionFlagValue, err := cmd.Flags().GetString(descriptionFlagName)
		if err != nil {
			return err, false
		}
		m.Description = descriptionFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveRunProfileRequestPropJvmArgumentsFlags(depth int, m *models.RunProfileRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	jvmArgumentsFlagName := fmt.Sprintf("%v.jvmArguments", cmdPrefix)
	if cmd.Flags().Changed(jvmArgumentsFlagName) {
		// warning: jvmArguments array type []string is not supported by go-swagger cli yet
	}

	return nil, retAdded
}

func retrieveRunProfileRequestPropNameFlags(depth int, m *models.RunProfileRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	nameFlagName := fmt.Sprintf("%v.name", cmdPrefix)
	if cmd.Flags().Changed(nameFlagName) {

		var nameFlagName string
		if cmdPrefix == "" {
			nameFlagName = "name"
		} else {
			nameFlagName = fmt.Sprintf("%v.name", cmdPrefix)
		}

		nameFlagValue, err := cmd.Flags().GetString(nameFlagName)
		if err != nil {
			return err, false
		}
		m.Name = &nameFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveRunProfileRequestPropRuntimeIDFlags(depth int, m *models.RunProfileRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	runtimeIdFlagName := fmt.Sprintf("%v.runtimeId", cmdPrefix)
	if cmd.Flags().Changed(runtimeIdFlagName) {

		var runtimeIdFlagName string
		if cmdPrefix == "" {
			runtimeIdFlagName = "runtimeId"
		} else {
			runtimeIdFlagName = fmt.Sprintf("%v.runtimeId", cmdPrefix)
		}

		runtimeIdFlagValue, err := cmd.Flags().GetString(runtimeIdFlagName)
		if err != nil {
			return err, false
		}
		m.RuntimeID = &runtimeIdFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveRunProfileRequestPropTypeFlags(depth int, m *models.RunProfileRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	typeFlagName := fmt.Sprintf("%v.type", cmdPrefix)
	if cmd.Flags().Changed(typeFlagName) {

		var typeFlagName string
		if cmdPrefix == "" {
			typeFlagName = "type"
		} else {
			typeFlagName = fmt.Sprintf("%v.type", cmdPrefix)
		}

		typeFlagValue, err := cmd.Flags().GetString(typeFlagName)
		if err != nil {
			return err, false
		}
		m.Type = &typeFlagValue

		retAdded = true
	}

	return nil, retAdded
}
