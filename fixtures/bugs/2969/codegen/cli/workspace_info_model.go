// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/go-swagger/go-swagger/fixtures/bugs/2969/codegen/models"
	"github.com/spf13/cobra"
)

// Schema cli for WorkspaceInfo

// register flags to command
func registerModelWorkspaceInfoFlags(depth int, cmdPrefix string, cmd *cobra.Command) error {

	if err := registerWorkspaceInfoPropDescription(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerWorkspaceInfoPropEnvironment(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerWorkspaceInfoPropID(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerWorkspaceInfoPropName(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerWorkspaceInfoPropOwner(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerWorkspaceInfoPropProtectedArtifactUpdate(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerWorkspaceInfoPropType(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	return nil
}

func registerWorkspaceInfoPropDescription(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	descriptionDescription := `Workspace description`

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

func registerWorkspaceInfoPropEnvironment(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	var environmentFlagName string
	if cmdPrefix == "" {
		environmentFlagName = "environment"
	} else {
		environmentFlagName = fmt.Sprintf("%v.environment", cmdPrefix)
	}

	if err := registerModelEnvironmentInfoFlags(depth+1, environmentFlagName, cmd); err != nil {
		return err
	}

	return nil
}

func registerWorkspaceInfoPropID(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	idDescription := `Workspace identifier`

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

func registerWorkspaceInfoPropName(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	nameDescription := `Workspace name`

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

func registerWorkspaceInfoPropOwner(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	ownerDescription := `Workspace owner`

	var ownerFlagName string
	if cmdPrefix == "" {
		ownerFlagName = "owner"
	} else {
		ownerFlagName = fmt.Sprintf("%v.owner", cmdPrefix)
	}

	var ownerFlagDefault string

	_ = cmd.PersistentFlags().String(ownerFlagName, ownerFlagDefault, ownerDescription)

	return nil
}

func registerWorkspaceInfoPropProtectedArtifactUpdate(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	protectedArtifactUpdateDescription := `Task update with workspace artifact only`

	var protectedArtifactUpdateFlagName string
	if cmdPrefix == "" {
		protectedArtifactUpdateFlagName = "protectedArtifactUpdate"
	} else {
		protectedArtifactUpdateFlagName = fmt.Sprintf("%v.protectedArtifactUpdate", cmdPrefix)
	}

	var protectedArtifactUpdateFlagDefault bool

	_ = cmd.PersistentFlags().Bool(protectedArtifactUpdateFlagName, protectedArtifactUpdateFlagDefault, protectedArtifactUpdateDescription)

	return nil
}

func registerWorkspaceInfoPropType(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	typeDescription := `Enum: ["shared","personal","custom"]. Workspace type`

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
			if err := json.Unmarshal([]byte(`["shared","personal","custom"]`), &res); err != nil {
				panic(err)
			}
			return res, cobra.ShellCompDirectiveDefault
		}); err != nil {
		return err
	}

	return nil
}

// retrieve flags from commands, and set value in model. Return true if any flag is passed by user to fill model field.
func retrieveModelWorkspaceInfoFlags(depth int, m *models.WorkspaceInfo, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false

	err, descriptionAdded := retrieveWorkspaceInfoPropDescriptionFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || descriptionAdded

	err, environmentAdded := retrieveWorkspaceInfoPropEnvironmentFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || environmentAdded

	err, idAdded := retrieveWorkspaceInfoPropIDFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || idAdded

	err, nameAdded := retrieveWorkspaceInfoPropNameFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || nameAdded

	err, ownerAdded := retrieveWorkspaceInfoPropOwnerFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || ownerAdded

	err, protectedArtifactUpdateAdded := retrieveWorkspaceInfoPropProtectedArtifactUpdateFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || protectedArtifactUpdateAdded

	err, typeAdded := retrieveWorkspaceInfoPropTypeFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || typeAdded

	return nil, retAdded
}

func retrieveWorkspaceInfoPropDescriptionFlags(depth int, m *models.WorkspaceInfo, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrieveWorkspaceInfoPropEnvironmentFlags(depth int, m *models.WorkspaceInfo, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	environmentFlagName := fmt.Sprintf("%v.environment", cmdPrefix)
	if cmd.Flags().Changed(environmentFlagName) {
		// info: complex object environment EnvironmentInfo is retrieved outside this Changed() block
	}
	environmentFlagValue := m.Environment
	if swag.IsZero(environmentFlagValue) {
		environmentFlagValue = &models.EnvironmentInfo{}
	}

	err, environmentAdded := retrieveModelEnvironmentInfoFlags(depth+1, environmentFlagValue, environmentFlagName, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || environmentAdded
	if environmentAdded {
		m.Environment = environmentFlagValue
	}

	return nil, retAdded
}

func retrieveWorkspaceInfoPropIDFlags(depth int, m *models.WorkspaceInfo, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrieveWorkspaceInfoPropNameFlags(depth int, m *models.WorkspaceInfo, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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
		m.Name = nameFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveWorkspaceInfoPropOwnerFlags(depth int, m *models.WorkspaceInfo, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	ownerFlagName := fmt.Sprintf("%v.owner", cmdPrefix)
	if cmd.Flags().Changed(ownerFlagName) {

		var ownerFlagName string
		if cmdPrefix == "" {
			ownerFlagName = "owner"
		} else {
			ownerFlagName = fmt.Sprintf("%v.owner", cmdPrefix)
		}

		ownerFlagValue, err := cmd.Flags().GetString(ownerFlagName)
		if err != nil {
			return err, false
		}
		m.Owner = ownerFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveWorkspaceInfoPropProtectedArtifactUpdateFlags(depth int, m *models.WorkspaceInfo, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	protectedArtifactUpdateFlagName := fmt.Sprintf("%v.protectedArtifactUpdate", cmdPrefix)
	if cmd.Flags().Changed(protectedArtifactUpdateFlagName) {

		var protectedArtifactUpdateFlagName string
		if cmdPrefix == "" {
			protectedArtifactUpdateFlagName = "protectedArtifactUpdate"
		} else {
			protectedArtifactUpdateFlagName = fmt.Sprintf("%v.protectedArtifactUpdate", cmdPrefix)
		}

		protectedArtifactUpdateFlagValue, err := cmd.Flags().GetBool(protectedArtifactUpdateFlagName)
		if err != nil {
			return err, false
		}
		m.ProtectedArtifactUpdate = protectedArtifactUpdateFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveWorkspaceInfoPropTypeFlags(depth int, m *models.WorkspaceInfo, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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
		m.Type = typeFlagValue

		retAdded = true
	}

	return nil, retAdded
}
