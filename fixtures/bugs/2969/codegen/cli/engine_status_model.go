// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-swagger/go-swagger/fixtures/bugs/2969/codegen/models"
	"github.com/spf13/cobra"
)

// Schema cli for EngineStatus

// register flags to command
func registerModelEngineStatusFlags(depth int, cmdPrefix string, cmd *cobra.Command) error {

	if err := registerEngineStatusPropConsumed(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerEngineStatusPropType(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	return nil
}

func registerEngineStatusPropConsumed(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	consumedDescription := `Required. Number of engines used by type`

	var consumedFlagName string
	if cmdPrefix == "" {
		consumedFlagName = "consumed"
	} else {
		consumedFlagName = fmt.Sprintf("%v.consumed", cmdPrefix)
	}

	var consumedFlagDefault int32

	_ = cmd.PersistentFlags().Int32(consumedFlagName, consumedFlagDefault, consumedDescription)

	return nil
}

func registerEngineStatusPropType(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	typeDescription := `Required. Engine type`

	var typeFlagName string
	if cmdPrefix == "" {
		typeFlagName = "type"
	} else {
		typeFlagName = fmt.Sprintf("%v.type", cmdPrefix)
	}

	var typeFlagDefault string

	_ = cmd.PersistentFlags().String(typeFlagName, typeFlagDefault, typeDescription)

	return nil
}

// retrieve flags from commands, and set value in model. Return true if any flag is passed by user to fill model field.
func retrieveModelEngineStatusFlags(depth int, m *models.EngineStatus, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false

	err, consumedAdded := retrieveEngineStatusPropConsumedFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || consumedAdded

	err, typeAdded := retrieveEngineStatusPropTypeFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || typeAdded

	return nil, retAdded
}

func retrieveEngineStatusPropConsumedFlags(depth int, m *models.EngineStatus, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	consumedFlagName := fmt.Sprintf("%v.consumed", cmdPrefix)
	if cmd.Flags().Changed(consumedFlagName) {

		var consumedFlagName string
		if cmdPrefix == "" {
			consumedFlagName = "consumed"
		} else {
			consumedFlagName = fmt.Sprintf("%v.consumed", cmdPrefix)
		}

		consumedFlagValue, err := cmd.Flags().GetInt32(consumedFlagName)
		if err != nil {
			return err, false
		}
		m.Consumed = &consumedFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveEngineStatusPropTypeFlags(depth int, m *models.EngineStatus, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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
