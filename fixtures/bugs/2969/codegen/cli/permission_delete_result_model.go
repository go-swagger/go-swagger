// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/spf13/cobra"
)

// Schema cli for PermissionDeleteResult

// register flags to command
func registerModelPermissionDeleteResultFlags(depth int, cmdPrefix string, cmd *cobra.Command) error {

	if err := registerPermissionDeleteResultPropMessage(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPermissionDeleteResultPropPermissionID(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPermissionDeleteResultPropStatus(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	return nil
}

func registerPermissionDeleteResultPropMessage(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	messageDescription := ``

	var messageFlagName string
	if cmdPrefix == "" {
		messageFlagName = "message"
	} else {
		messageFlagName = fmt.Sprintf("%v.message", cmdPrefix)
	}

	var messageFlagDefault string

	_ = cmd.PersistentFlags().String(messageFlagName, messageFlagDefault, messageDescription)

	return nil
}

func registerPermissionDeleteResultPropPermissionID(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	var permissionIdFlagName string
	if cmdPrefix == "" {
		permissionIdFlagName = "permissionId"
	} else {
		permissionIdFlagName = fmt.Sprintf("%v.permissionId", cmdPrefix)
	}

	if err := registerModelPermissionIDFlags(depth+1, permissionIdFlagName, cmd); err != nil {
		return err
	}

	return nil
}

func registerPermissionDeleteResultPropStatus(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	statusDescription := ``

	var statusFlagName string
	if cmdPrefix == "" {
		statusFlagName = "status"
	} else {
		statusFlagName = fmt.Sprintf("%v.status", cmdPrefix)
	}

	var statusFlagDefault int32

	_ = cmd.PersistentFlags().Int32(statusFlagName, statusFlagDefault, statusDescription)

	return nil
}

// retrieve flags from commands, and set value in model. Return true if any flag is passed by user to fill model field.
func retrieveModelPermissionDeleteResultFlags(depth int, m *models.PermissionDeleteResult, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false

	err, messageAdded := retrievePermissionDeleteResultPropMessageFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || messageAdded

	err, permissionIdAdded := retrievePermissionDeleteResultPropPermissionIDFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || permissionIdAdded

	err, statusAdded := retrievePermissionDeleteResultPropStatusFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || statusAdded

	return nil, retAdded
}

func retrievePermissionDeleteResultPropMessageFlags(depth int, m *models.PermissionDeleteResult, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	messageFlagName := fmt.Sprintf("%v.message", cmdPrefix)
	if cmd.Flags().Changed(messageFlagName) {

		var messageFlagName string
		if cmdPrefix == "" {
			messageFlagName = "message"
		} else {
			messageFlagName = fmt.Sprintf("%v.message", cmdPrefix)
		}

		messageFlagValue, err := cmd.Flags().GetString(messageFlagName)
		if err != nil {
			return err, false
		}
		m.Message = messageFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrievePermissionDeleteResultPropPermissionIDFlags(depth int, m *models.PermissionDeleteResult, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	permissionIdFlagName := fmt.Sprintf("%v.permissionId", cmdPrefix)
	if cmd.Flags().Changed(permissionIdFlagName) {
		// info: complex object permissionId PermissionID is retrieved outside this Changed() block
	}
	permissionIdFlagValue := m.PermissionID
	if swag.IsZero(permissionIdFlagValue) {
		permissionIdFlagValue = &models.PermissionID{}
	}

	err, permissionIdAdded := retrieveModelPermissionIDFlags(depth+1, permissionIdFlagValue, permissionIdFlagName, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || permissionIdAdded
	if permissionIdAdded {
		m.PermissionID = permissionIdFlagValue
	}

	return nil, retAdded
}

func retrievePermissionDeleteResultPropStatusFlags(depth int, m *models.PermissionDeleteResult, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	statusFlagName := fmt.Sprintf("%v.status", cmdPrefix)
	if cmd.Flags().Changed(statusFlagName) {

		var statusFlagName string
		if cmdPrefix == "" {
			statusFlagName = "status"
		} else {
			statusFlagName = fmt.Sprintf("%v.status", cmdPrefix)
		}

		statusFlagValue, err := cmd.Flags().GetInt32(statusFlagName)
		if err != nil {
			return err, false
		}
		m.Status = statusFlagValue

		retAdded = true
	}

	return nil, retAdded
}
