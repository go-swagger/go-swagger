// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-swagger/go-swagger/fixtures/bugs/2969/codegen/models"
	"github.com/spf13/cobra"
)

// Schema cli for Schedule

// register flags to command
func registerModelScheduleFlags(depth int, cmdPrefix string, cmd *cobra.Command) error {

	if err := registerSchedulePropDescription(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerSchedulePropEnvironmentID(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerSchedulePropExecutableID(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerSchedulePropExecutableType(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerSchedulePropID(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerSchedulePropTriggers(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	return nil
}

func registerSchedulePropDescription(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	descriptionDescription := ``

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

func registerSchedulePropEnvironmentID(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	environmentIdDescription := ``

	var environmentIdFlagName string
	if cmdPrefix == "" {
		environmentIdFlagName = "environmentId"
	} else {
		environmentIdFlagName = fmt.Sprintf("%v.environmentId", cmdPrefix)
	}

	var environmentIdFlagDefault string

	_ = cmd.PersistentFlags().String(environmentIdFlagName, environmentIdFlagDefault, environmentIdDescription)

	return nil
}

func registerSchedulePropExecutableID(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	executableIdDescription := ``

	var executableIdFlagName string
	if cmdPrefix == "" {
		executableIdFlagName = "executableId"
	} else {
		executableIdFlagName = fmt.Sprintf("%v.executableId", cmdPrefix)
	}

	var executableIdFlagDefault string

	_ = cmd.PersistentFlags().String(executableIdFlagName, executableIdFlagDefault, executableIdDescription)

	return nil
}

func registerSchedulePropExecutableType(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	executableTypeDescription := ``

	var executableTypeFlagName string
	if cmdPrefix == "" {
		executableTypeFlagName = "executableType"
	} else {
		executableTypeFlagName = fmt.Sprintf("%v.executableType", cmdPrefix)
	}

	var executableTypeFlagDefault string

	_ = cmd.PersistentFlags().String(executableTypeFlagName, executableTypeFlagDefault, executableTypeDescription)

	return nil
}

func registerSchedulePropID(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	idDescription := ``

	var idFlagName string
	if cmdPrefix == "" {
		idFlagName = "id"
	} else {
		idFlagName = fmt.Sprintf("%v.id", cmdPrefix)
	}

	var idFlagDefault string

	_ = cmd.PersistentFlags().String(idFlagName, idFlagDefault, idDescription)

	return nil
}

func registerSchedulePropTriggers(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	// warning: triggers []*MultipleTrigger array type is not supported by go-swagger cli yet

	return nil
}

// retrieve flags from commands, and set value in model. Return true if any flag is passed by user to fill model field.
func retrieveModelScheduleFlags(depth int, m *models.Schedule, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false

	err, descriptionAdded := retrieveSchedulePropDescriptionFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || descriptionAdded

	err, environmentIdAdded := retrieveSchedulePropEnvironmentIDFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || environmentIdAdded

	err, executableIdAdded := retrieveSchedulePropExecutableIDFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || executableIdAdded

	err, executableTypeAdded := retrieveSchedulePropExecutableTypeFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || executableTypeAdded

	err, idAdded := retrieveSchedulePropIDFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || idAdded

	err, triggersAdded := retrieveSchedulePropTriggersFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || triggersAdded

	return nil, retAdded
}

func retrieveSchedulePropDescriptionFlags(depth int, m *models.Schedule, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrieveSchedulePropEnvironmentIDFlags(depth int, m *models.Schedule, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	environmentIdFlagName := fmt.Sprintf("%v.environmentId", cmdPrefix)
	if cmd.Flags().Changed(environmentIdFlagName) {

		var environmentIdFlagName string
		if cmdPrefix == "" {
			environmentIdFlagName = "environmentId"
		} else {
			environmentIdFlagName = fmt.Sprintf("%v.environmentId", cmdPrefix)
		}

		environmentIdFlagValue, err := cmd.Flags().GetString(environmentIdFlagName)
		if err != nil {
			return err, false
		}
		m.EnvironmentID = environmentIdFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveSchedulePropExecutableIDFlags(depth int, m *models.Schedule, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	executableIdFlagName := fmt.Sprintf("%v.executableId", cmdPrefix)
	if cmd.Flags().Changed(executableIdFlagName) {

		var executableIdFlagName string
		if cmdPrefix == "" {
			executableIdFlagName = "executableId"
		} else {
			executableIdFlagName = fmt.Sprintf("%v.executableId", cmdPrefix)
		}

		executableIdFlagValue, err := cmd.Flags().GetString(executableIdFlagName)
		if err != nil {
			return err, false
		}
		m.ExecutableID = executableIdFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveSchedulePropExecutableTypeFlags(depth int, m *models.Schedule, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	executableTypeFlagName := fmt.Sprintf("%v.executableType", cmdPrefix)
	if cmd.Flags().Changed(executableTypeFlagName) {

		var executableTypeFlagName string
		if cmdPrefix == "" {
			executableTypeFlagName = "executableType"
		} else {
			executableTypeFlagName = fmt.Sprintf("%v.executableType", cmdPrefix)
		}

		executableTypeFlagValue, err := cmd.Flags().GetString(executableTypeFlagName)
		if err != nil {
			return err, false
		}
		m.ExecutableType = executableTypeFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveSchedulePropIDFlags(depth int, m *models.Schedule, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	idFlagName := fmt.Sprintf("%v.id", cmdPrefix)
	if cmd.Flags().Changed(idFlagName) {

		var idFlagName string
		if cmdPrefix == "" {
			idFlagName = "id"
		} else {
			idFlagName = fmt.Sprintf("%v.id", cmdPrefix)
		}

		idFlagValue, err := cmd.Flags().GetString(idFlagName)
		if err != nil {
			return err, false
		}
		m.ID = idFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveSchedulePropTriggersFlags(depth int, m *models.Schedule, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	triggersFlagName := fmt.Sprintf("%v.triggers", cmdPrefix)
	if cmd.Flags().Changed(triggersFlagName) {
		// warning: triggers array type []*MultipleTrigger is not supported by go-swagger cli yet
	}

	return nil, retAdded
}
