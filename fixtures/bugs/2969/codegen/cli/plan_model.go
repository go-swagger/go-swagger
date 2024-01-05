// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/spf13/cobra"
)

// Schema cli for Plan

// register flags to command
func registerModelPlanFlags(depth int, cmdPrefix string, cmd *cobra.Command) error {

	if err := registerPlanPropAuthor(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPlanPropAuthorType(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPlanPropCreateDate(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPlanPropDescription(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPlanPropID(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPlanPropName(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPlanPropSteps(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPlanPropUpdateDate(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPlanPropWorkspace(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	return nil
}

func registerPlanPropAuthor(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	authorDescription := `Required. Plan author`

	var authorFlagName string
	if cmdPrefix == "" {
		authorFlagName = "author"
	} else {
		authorFlagName = fmt.Sprintf("%v.author", cmdPrefix)
	}

	var authorFlagDefault string

	_ = cmd.PersistentFlags().String(authorFlagName, authorFlagDefault, authorDescription)

	return nil
}

func registerPlanPropAuthorType(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	authorTypeDescription := `Enum: ["HUMAN","SERVICE"]. Required. Plan author type (HUMAN or SERVICE)`

	var authorTypeFlagName string
	if cmdPrefix == "" {
		authorTypeFlagName = "authorType"
	} else {
		authorTypeFlagName = fmt.Sprintf("%v.authorType", cmdPrefix)
	}

	var authorTypeFlagDefault string

	_ = cmd.PersistentFlags().String(authorTypeFlagName, authorTypeFlagDefault, authorTypeDescription)

	if err := cmd.RegisterFlagCompletionFunc(authorTypeFlagName,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			var res []string
			if err := json.Unmarshal([]byte(`["HUMAN","SERVICE"]`), &res); err != nil {
				panic(err)
			}
			return res, cobra.ShellCompDirectiveDefault
		}); err != nil {
		return err
	}

	return nil
}

func registerPlanPropCreateDate(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	createDateDescription := `Created on`

	var createDateFlagName string
	if cmdPrefix == "" {
		createDateFlagName = "createDate"
	} else {
		createDateFlagName = fmt.Sprintf("%v.createDate", cmdPrefix)
	}

	_ = cmd.PersistentFlags().String(createDateFlagName, "", createDateDescription)

	return nil
}

func registerPlanPropDescription(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	descriptionDescription := `Resource description`

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

func registerPlanPropID(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	idDescription := `Required. Resource id`

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

func registerPlanPropName(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	nameDescription := `Required. Resource name`

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

func registerPlanPropSteps(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	// warning: steps []*Step array type is not supported by go-swagger cli yet

	return nil
}

func registerPlanPropUpdateDate(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	updateDateDescription := `Updated on`

	var updateDateFlagName string
	if cmdPrefix == "" {
		updateDateFlagName = "updateDate"
	} else {
		updateDateFlagName = fmt.Sprintf("%v.updateDate", cmdPrefix)
	}

	_ = cmd.PersistentFlags().String(updateDateFlagName, "", updateDateDescription)

	return nil
}

func registerPlanPropWorkspace(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	var workspaceFlagName string
	if cmdPrefix == "" {
		workspaceFlagName = "workspace"
	} else {
		workspaceFlagName = fmt.Sprintf("%v.workspace", cmdPrefix)
	}

	if err := registerModelWorkspaceInfoFlags(depth+1, workspaceFlagName, cmd); err != nil {
		return err
	}

	return nil
}

// retrieve flags from commands, and set value in model. Return true if any flag is passed by user to fill model field.
func retrieveModelPlanFlags(depth int, m *models.Plan, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false

	err, authorAdded := retrievePlanPropAuthorFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || authorAdded

	err, authorTypeAdded := retrievePlanPropAuthorTypeFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || authorTypeAdded

	err, createDateAdded := retrievePlanPropCreateDateFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || createDateAdded

	err, descriptionAdded := retrievePlanPropDescriptionFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || descriptionAdded

	err, idAdded := retrievePlanPropIDFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || idAdded

	err, nameAdded := retrievePlanPropNameFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || nameAdded

	err, stepsAdded := retrievePlanPropStepsFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || stepsAdded

	err, updateDateAdded := retrievePlanPropUpdateDateFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || updateDateAdded

	err, workspaceAdded := retrievePlanPropWorkspaceFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || workspaceAdded

	return nil, retAdded
}

func retrievePlanPropAuthorFlags(depth int, m *models.Plan, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	authorFlagName := fmt.Sprintf("%v.author", cmdPrefix)
	if cmd.Flags().Changed(authorFlagName) {

		var authorFlagName string
		if cmdPrefix == "" {
			authorFlagName = "author"
		} else {
			authorFlagName = fmt.Sprintf("%v.author", cmdPrefix)
		}

		authorFlagValue, err := cmd.Flags().GetString(authorFlagName)
		if err != nil {
			return err, false
		}
		m.Author = &authorFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrievePlanPropAuthorTypeFlags(depth int, m *models.Plan, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	authorTypeFlagName := fmt.Sprintf("%v.authorType", cmdPrefix)
	if cmd.Flags().Changed(authorTypeFlagName) {

		var authorTypeFlagName string
		if cmdPrefix == "" {
			authorTypeFlagName = "authorType"
		} else {
			authorTypeFlagName = fmt.Sprintf("%v.authorType", cmdPrefix)
		}

		authorTypeFlagValue, err := cmd.Flags().GetString(authorTypeFlagName)
		if err != nil {
			return err, false
		}
		m.AuthorType = &authorTypeFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrievePlanPropCreateDateFlags(depth int, m *models.Plan, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	createDateFlagName := fmt.Sprintf("%v.createDate", cmdPrefix)
	if cmd.Flags().Changed(createDateFlagName) {

		var createDateFlagName string
		if cmdPrefix == "" {
			createDateFlagName = "createDate"
		} else {
			createDateFlagName = fmt.Sprintf("%v.createDate", cmdPrefix)
		}

		createDateFlagValueStr, err := cmd.Flags().GetString(createDateFlagName)
		if err != nil {
			return err, false
		}
		var createDateFlagValue strfmt.DateTime
		if err := createDateFlagValue.UnmarshalText([]byte(createDateFlagValueStr)); err != nil {
			return err, false
		}
		m.CreateDate = createDateFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrievePlanPropDescriptionFlags(depth int, m *models.Plan, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrievePlanPropIDFlags(depth int, m *models.Plan, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrievePlanPropNameFlags(depth int, m *models.Plan, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrievePlanPropStepsFlags(depth int, m *models.Plan, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	stepsFlagName := fmt.Sprintf("%v.steps", cmdPrefix)
	if cmd.Flags().Changed(stepsFlagName) {
		// warning: steps array type []*Step is not supported by go-swagger cli yet
	}

	return nil, retAdded
}

func retrievePlanPropUpdateDateFlags(depth int, m *models.Plan, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	updateDateFlagName := fmt.Sprintf("%v.updateDate", cmdPrefix)
	if cmd.Flags().Changed(updateDateFlagName) {

		var updateDateFlagName string
		if cmdPrefix == "" {
			updateDateFlagName = "updateDate"
		} else {
			updateDateFlagName = fmt.Sprintf("%v.updateDate", cmdPrefix)
		}

		updateDateFlagValueStr, err := cmd.Flags().GetString(updateDateFlagName)
		if err != nil {
			return err, false
		}
		var updateDateFlagValue strfmt.DateTime
		if err := updateDateFlagValue.UnmarshalText([]byte(updateDateFlagValueStr)); err != nil {
			return err, false
		}
		m.UpdateDate = updateDateFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrievePlanPropWorkspaceFlags(depth int, m *models.Plan, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	workspaceFlagName := fmt.Sprintf("%v.workspace", cmdPrefix)
	if cmd.Flags().Changed(workspaceFlagName) {
		// info: complex object workspace WorkspaceInfo is retrieved outside this Changed() block
	}
	workspaceFlagValue := m.Workspace
	if swag.IsZero(workspaceFlagValue) {
		workspaceFlagValue = &models.WorkspaceInfo{}
	}

	err, workspaceAdded := retrieveModelWorkspaceInfoFlags(depth+1, workspaceFlagValue, workspaceFlagName, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || workspaceAdded
	if workspaceAdded {
		m.Workspace = workspaceFlagValue
	}

	return nil, retAdded
}
