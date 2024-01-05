// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/go-swagger/go-swagger/fixtures/bugs/2969/codegen/models"
	"github.com/spf13/cobra"
)

// Schema cli for TaskAutoUpgradeRequest

// register flags to command
func registerModelTaskAutoUpgradeRequestFlags(depth int, cmdPrefix string, cmd *cobra.Command) error {

	if err := registerTaskAutoUpgradeRequestPropArtifact(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerTaskAutoUpgradeRequestPropAutoUpgradeInfo(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerTaskAutoUpgradeRequestPropConnections(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerTaskAutoUpgradeRequestPropDescription(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerTaskAutoUpgradeRequestPropName(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerTaskAutoUpgradeRequestPropParameters(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerTaskAutoUpgradeRequestPropResources(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerTaskAutoUpgradeRequestPropTags(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerTaskAutoUpgradeRequestPropWorkspaceID(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	return nil
}

func registerTaskAutoUpgradeRequestPropArtifact(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	var artifactFlagName string
	if cmdPrefix == "" {
		artifactFlagName = "artifact"
	} else {
		artifactFlagName = fmt.Sprintf("%v.artifact", cmdPrefix)
	}

	if err := registerModelArtifactRequestFlags(depth+1, artifactFlagName, cmd); err != nil {
		return err
	}

	return nil
}

func registerTaskAutoUpgradeRequestPropAutoUpgradeInfo(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	var autoUpgradeInfoFlagName string
	if cmdPrefix == "" {
		autoUpgradeInfoFlagName = "autoUpgradeInfo"
	} else {
		autoUpgradeInfoFlagName = fmt.Sprintf("%v.autoUpgradeInfo", cmdPrefix)
	}

	if err := registerModelAutoUpgradeInfoFlags(depth+1, autoUpgradeInfoFlagName, cmd); err != nil {
		return err
	}

	return nil
}

func registerTaskAutoUpgradeRequestPropConnections(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	// warning: connections map[string]string map type is not supported by go-swagger cli yet

	return nil
}

func registerTaskAutoUpgradeRequestPropDescription(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	descriptionDescription := `Required. Task description`

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

func registerTaskAutoUpgradeRequestPropName(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	nameDescription := `Required. Task name`

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

func registerTaskAutoUpgradeRequestPropParameters(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	// warning: parameters map[string]string map type is not supported by go-swagger cli yet

	return nil
}

func registerTaskAutoUpgradeRequestPropResources(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	// warning: resources map[string]string map type is not supported by go-swagger cli yet

	return nil
}

func registerTaskAutoUpgradeRequestPropTags(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	// warning: tags []string array type is not supported by go-swagger cli yet

	return nil
}

func registerTaskAutoUpgradeRequestPropWorkspaceID(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	workspaceIdDescription := `Required. Workspace id of task to create`

	var workspaceIdFlagName string
	if cmdPrefix == "" {
		workspaceIdFlagName = "workspaceId"
	} else {
		workspaceIdFlagName = fmt.Sprintf("%v.workspaceId", cmdPrefix)
	}

	var workspaceIdFlagDefault string

	_ = cmd.PersistentFlags().String(workspaceIdFlagName, workspaceIdFlagDefault, workspaceIdDescription)

	return nil
}

// retrieve flags from commands, and set value in model. Return true if any flag is passed by user to fill model field.
func retrieveModelTaskAutoUpgradeRequestFlags(depth int, m *models.TaskAutoUpgradeRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false

	err, artifactAdded := retrieveTaskAutoUpgradeRequestPropArtifactFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || artifactAdded

	err, autoUpgradeInfoAdded := retrieveTaskAutoUpgradeRequestPropAutoUpgradeInfoFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || autoUpgradeInfoAdded

	err, connectionsAdded := retrieveTaskAutoUpgradeRequestPropConnectionsFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || connectionsAdded

	err, descriptionAdded := retrieveTaskAutoUpgradeRequestPropDescriptionFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || descriptionAdded

	err, nameAdded := retrieveTaskAutoUpgradeRequestPropNameFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || nameAdded

	err, parametersAdded := retrieveTaskAutoUpgradeRequestPropParametersFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || parametersAdded

	err, resourcesAdded := retrieveTaskAutoUpgradeRequestPropResourcesFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || resourcesAdded

	err, tagsAdded := retrieveTaskAutoUpgradeRequestPropTagsFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || tagsAdded

	err, workspaceIdAdded := retrieveTaskAutoUpgradeRequestPropWorkspaceIDFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || workspaceIdAdded

	return nil, retAdded
}

func retrieveTaskAutoUpgradeRequestPropArtifactFlags(depth int, m *models.TaskAutoUpgradeRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	artifactFlagName := fmt.Sprintf("%v.artifact", cmdPrefix)
	if cmd.Flags().Changed(artifactFlagName) {
		// info: complex object artifact ArtifactRequest is retrieved outside this Changed() block
	}
	artifactFlagValue := m.Artifact
	if swag.IsZero(artifactFlagValue) {
		artifactFlagValue = &models.ArtifactRequest{}
	}

	err, artifactAdded := retrieveModelArtifactRequestFlags(depth+1, artifactFlagValue, artifactFlagName, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || artifactAdded
	if artifactAdded {
		m.Artifact = artifactFlagValue
	}

	return nil, retAdded
}

func retrieveTaskAutoUpgradeRequestPropAutoUpgradeInfoFlags(depth int, m *models.TaskAutoUpgradeRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	autoUpgradeInfoFlagName := fmt.Sprintf("%v.autoUpgradeInfo", cmdPrefix)
	if cmd.Flags().Changed(autoUpgradeInfoFlagName) {
		// info: complex object autoUpgradeInfo AutoUpgradeInfo is retrieved outside this Changed() block
	}
	autoUpgradeInfoFlagValue := m.AutoUpgradeInfo
	if swag.IsZero(autoUpgradeInfoFlagValue) {
		autoUpgradeInfoFlagValue = &models.AutoUpgradeInfo{}
	}

	err, autoUpgradeInfoAdded := retrieveModelAutoUpgradeInfoFlags(depth+1, autoUpgradeInfoFlagValue, autoUpgradeInfoFlagName, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || autoUpgradeInfoAdded
	if autoUpgradeInfoAdded {
		m.AutoUpgradeInfo = autoUpgradeInfoFlagValue
	}

	return nil, retAdded
}

func retrieveTaskAutoUpgradeRequestPropConnectionsFlags(depth int, m *models.TaskAutoUpgradeRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	connectionsFlagName := fmt.Sprintf("%v.connections", cmdPrefix)
	if cmd.Flags().Changed(connectionsFlagName) {
		// warning: connections map type map[string]string is not supported by go-swagger cli yet
	}

	return nil, retAdded
}

func retrieveTaskAutoUpgradeRequestPropDescriptionFlags(depth int, m *models.TaskAutoUpgradeRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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
		m.Description = &descriptionFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveTaskAutoUpgradeRequestPropNameFlags(depth int, m *models.TaskAutoUpgradeRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrieveTaskAutoUpgradeRequestPropParametersFlags(depth int, m *models.TaskAutoUpgradeRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	parametersFlagName := fmt.Sprintf("%v.parameters", cmdPrefix)
	if cmd.Flags().Changed(parametersFlagName) {
		// warning: parameters map type map[string]string is not supported by go-swagger cli yet
	}

	return nil, retAdded
}

func retrieveTaskAutoUpgradeRequestPropResourcesFlags(depth int, m *models.TaskAutoUpgradeRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	resourcesFlagName := fmt.Sprintf("%v.resources", cmdPrefix)
	if cmd.Flags().Changed(resourcesFlagName) {
		// warning: resources map type map[string]string is not supported by go-swagger cli yet
	}

	return nil, retAdded
}

func retrieveTaskAutoUpgradeRequestPropTagsFlags(depth int, m *models.TaskAutoUpgradeRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	tagsFlagName := fmt.Sprintf("%v.tags", cmdPrefix)
	if cmd.Flags().Changed(tagsFlagName) {
		// warning: tags array type []string is not supported by go-swagger cli yet
	}

	return nil, retAdded
}

func retrieveTaskAutoUpgradeRequestPropWorkspaceIDFlags(depth int, m *models.TaskAutoUpgradeRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	workspaceIdFlagName := fmt.Sprintf("%v.workspaceId", cmdPrefix)
	if cmd.Flags().Changed(workspaceIdFlagName) {

		var workspaceIdFlagName string
		if cmdPrefix == "" {
			workspaceIdFlagName = "workspaceId"
		} else {
			workspaceIdFlagName = fmt.Sprintf("%v.workspaceId", cmdPrefix)
		}

		workspaceIdFlagValue, err := cmd.Flags().GetString(workspaceIdFlagName)
		if err != nil {
			return err, false
		}
		m.WorkspaceID = &workspaceIdFlagValue

		retAdded = true
	}

	return nil, retAdded
}
