// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/spf13/cobra"
)

// Schema cli for Resource

// register flags to command
func registerModelResourceFlags(depth int, cmdPrefix string, cmd *cobra.Command) error {

	if err := registerResourcePropDescription(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerResourcePropFile(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerResourcePropID(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerResourcePropName(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerResourcePropWorkspaceInfo(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	return nil
}

func registerResourcePropDescription(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	descriptionDescription := `Description of resource`

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

func registerResourcePropFile(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	fileDescription := `Required. Boolean value which indicates that resource is file`

	var fileFlagName string
	if cmdPrefix == "" {
		fileFlagName = "file"
	} else {
		fileFlagName = fmt.Sprintf("%v.file", cmdPrefix)
	}

	var fileFlagDefault bool

	_ = cmd.PersistentFlags().Bool(fileFlagName, fileFlagDefault, fileDescription)

	return nil
}

func registerResourcePropID(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	idDescription := `Required. Id of resource`

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

func registerResourcePropName(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	nameDescription := `Required. Name of resource`

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

func registerResourcePropWorkspaceInfo(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	var workspaceInfoFlagName string
	if cmdPrefix == "" {
		workspaceInfoFlagName = "workspaceInfo"
	} else {
		workspaceInfoFlagName = fmt.Sprintf("%v.workspaceInfo", cmdPrefix)
	}

	if err := registerModelWorkspaceInfoFlags(depth+1, workspaceInfoFlagName, cmd); err != nil {
		return err
	}

	return nil
}

// retrieve flags from commands, and set value in model. Return true if any flag is passed by user to fill model field.
func retrieveModelResourceFlags(depth int, m *models.Resource, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false

	err, descriptionAdded := retrieveResourcePropDescriptionFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || descriptionAdded

	err, fileAdded := retrieveResourcePropFileFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || fileAdded

	err, idAdded := retrieveResourcePropIDFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || idAdded

	err, nameAdded := retrieveResourcePropNameFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || nameAdded

	err, workspaceInfoAdded := retrieveResourcePropWorkspaceInfoFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || workspaceInfoAdded

	return nil, retAdded
}

func retrieveResourcePropDescriptionFlags(depth int, m *models.Resource, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrieveResourcePropFileFlags(depth int, m *models.Resource, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	fileFlagName := fmt.Sprintf("%v.file", cmdPrefix)
	if cmd.Flags().Changed(fileFlagName) {

		var fileFlagName string
		if cmdPrefix == "" {
			fileFlagName = "file"
		} else {
			fileFlagName = fmt.Sprintf("%v.file", cmdPrefix)
		}

		fileFlagValue, err := cmd.Flags().GetBool(fileFlagName)
		if err != nil {
			return err, false
		}
		m.File = &fileFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveResourcePropIDFlags(depth int, m *models.Resource, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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
		m.ID = &idFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveResourcePropNameFlags(depth int, m *models.Resource, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrieveResourcePropWorkspaceInfoFlags(depth int, m *models.Resource, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	workspaceInfoFlagName := fmt.Sprintf("%v.workspaceInfo", cmdPrefix)
	if cmd.Flags().Changed(workspaceInfoFlagName) {
		// info: complex object workspaceInfo WorkspaceInfo is retrieved outside this Changed() block
	}
	workspaceInfoFlagValue := m.WorkspaceInfo
	if swag.IsZero(workspaceInfoFlagValue) {
		workspaceInfoFlagValue = &models.WorkspaceInfo{}
	}

	err, workspaceInfoAdded := retrieveModelWorkspaceInfoFlags(depth+1, workspaceInfoFlagValue, workspaceInfoFlagName, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || workspaceInfoAdded
	if workspaceInfoAdded {
		m.WorkspaceInfo = workspaceInfoFlagValue
	}

	return nil, retAdded
}
