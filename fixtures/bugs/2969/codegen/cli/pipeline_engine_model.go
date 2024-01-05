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

// Schema cli for PipelineEngine

// register flags to command
func registerModelPipelineEngineFlags(depth int, cmdPrefix string, cmd *cobra.Command) error {

	if err := registerPipelineEnginePropAvailability(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPipelineEnginePropCloudRunner(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPipelineEnginePropCreateDate(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPipelineEnginePropDescription(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPipelineEnginePropID(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPipelineEnginePropManaged(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPipelineEnginePropName(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPipelineEnginePropPreAuthorizedKey(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPipelineEnginePropRunStatus(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPipelineEnginePropRuntimeID(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPipelineEnginePropStatus(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPipelineEnginePropUpdateDate(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPipelineEnginePropWorkspace(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	return nil
}

func registerPipelineEnginePropAvailability(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	availabilityDescription := `Enum: ["AVAILABLE","NOT_AVAILABLE","RETIRED"]. Availability status of engine|cluster`

	var availabilityFlagName string
	if cmdPrefix == "" {
		availabilityFlagName = "availability"
	} else {
		availabilityFlagName = fmt.Sprintf("%v.availability", cmdPrefix)
	}

	var availabilityFlagDefault string

	_ = cmd.PersistentFlags().String(availabilityFlagName, availabilityFlagDefault, availabilityDescription)

	if err := cmd.RegisterFlagCompletionFunc(availabilityFlagName,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			var res []string
			if err := json.Unmarshal([]byte(`["AVAILABLE","NOT_AVAILABLE","RETIRED"]`), &res); err != nil {
				panic(err)
			}
			return res, cobra.ShellCompDirectiveDefault
		}); err != nil {
		return err
	}

	return nil
}

func registerPipelineEnginePropCloudRunner(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	cloudRunnerDescription := ``

	var cloudRunnerFlagName string
	if cmdPrefix == "" {
		cloudRunnerFlagName = "cloudRunner"
	} else {
		cloudRunnerFlagName = fmt.Sprintf("%v.cloudRunner", cmdPrefix)
	}

	var cloudRunnerFlagDefault bool

	_ = cmd.PersistentFlags().Bool(cloudRunnerFlagName, cloudRunnerFlagDefault, cloudRunnerDescription)

	return nil
}

func registerPipelineEnginePropCreateDate(depth int, cmdPrefix string, cmd *cobra.Command) error {
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

func registerPipelineEnginePropDescription(depth int, cmdPrefix string, cmd *cobra.Command) error {
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

func registerPipelineEnginePropID(depth int, cmdPrefix string, cmd *cobra.Command) error {
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

func registerPipelineEnginePropManaged(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	managedDescription := `Indicates whether target runtime (engine/cluster) is managed or not`

	var managedFlagName string
	if cmdPrefix == "" {
		managedFlagName = "managed"
	} else {
		managedFlagName = fmt.Sprintf("%v.managed", cmdPrefix)
	}

	var managedFlagDefault bool

	_ = cmd.PersistentFlags().Bool(managedFlagName, managedFlagDefault, managedDescription)

	return nil
}

func registerPipelineEnginePropName(depth int, cmdPrefix string, cmd *cobra.Command) error {
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

func registerPipelineEnginePropPreAuthorizedKey(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	preAuthorizedKeyDescription := ``

	var preAuthorizedKeyFlagName string
	if cmdPrefix == "" {
		preAuthorizedKeyFlagName = "preAuthorizedKey"
	} else {
		preAuthorizedKeyFlagName = fmt.Sprintf("%v.preAuthorizedKey", cmdPrefix)
	}

	var preAuthorizedKeyFlagDefault string

	_ = cmd.PersistentFlags().String(preAuthorizedKeyFlagName, preAuthorizedKeyFlagDefault, preAuthorizedKeyDescription)

	return nil
}

func registerPipelineEnginePropRunStatus(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	runStatusDescription := `Enum: ["STARTING","RUNNING","STOPPING","STOPPED","START_FAILED","STOP_FAILED","ERROR","NOT_FOUND"]. `

	var runStatusFlagName string
	if cmdPrefix == "" {
		runStatusFlagName = "runStatus"
	} else {
		runStatusFlagName = fmt.Sprintf("%v.runStatus", cmdPrefix)
	}

	var runStatusFlagDefault string

	_ = cmd.PersistentFlags().String(runStatusFlagName, runStatusFlagDefault, runStatusDescription)

	if err := cmd.RegisterFlagCompletionFunc(runStatusFlagName,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			var res []string
			if err := json.Unmarshal([]byte(`["STARTING","RUNNING","STOPPING","STOPPED","START_FAILED","STOP_FAILED","ERROR","NOT_FOUND"]`), &res); err != nil {
				panic(err)
			}
			return res, cobra.ShellCompDirectiveDefault
		}); err != nil {
		return err
	}

	return nil
}

func registerPipelineEnginePropRuntimeID(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	runtimeIdDescription := `Required. Resource runtime id`

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

func registerPipelineEnginePropStatus(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	statusDescription := `Enum: ["PAIRED","NOT_PAIRED"]. `

	var statusFlagName string
	if cmdPrefix == "" {
		statusFlagName = "status"
	} else {
		statusFlagName = fmt.Sprintf("%v.status", cmdPrefix)
	}

	var statusFlagDefault string

	_ = cmd.PersistentFlags().String(statusFlagName, statusFlagDefault, statusDescription)

	if err := cmd.RegisterFlagCompletionFunc(statusFlagName,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			var res []string
			if err := json.Unmarshal([]byte(`["PAIRED","NOT_PAIRED"]`), &res); err != nil {
				panic(err)
			}
			return res, cobra.ShellCompDirectiveDefault
		}); err != nil {
		return err
	}

	return nil
}

func registerPipelineEnginePropUpdateDate(depth int, cmdPrefix string, cmd *cobra.Command) error {
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

func registerPipelineEnginePropWorkspace(depth int, cmdPrefix string, cmd *cobra.Command) error {
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
func retrieveModelPipelineEngineFlags(depth int, m *models.PipelineEngine, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false

	err, availabilityAdded := retrievePipelineEnginePropAvailabilityFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || availabilityAdded

	err, cloudRunnerAdded := retrievePipelineEnginePropCloudRunnerFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || cloudRunnerAdded

	err, createDateAdded := retrievePipelineEnginePropCreateDateFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || createDateAdded

	err, descriptionAdded := retrievePipelineEnginePropDescriptionFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || descriptionAdded

	err, idAdded := retrievePipelineEnginePropIDFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || idAdded

	err, managedAdded := retrievePipelineEnginePropManagedFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || managedAdded

	err, nameAdded := retrievePipelineEnginePropNameFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || nameAdded

	err, preAuthorizedKeyAdded := retrievePipelineEnginePropPreAuthorizedKeyFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || preAuthorizedKeyAdded

	err, runStatusAdded := retrievePipelineEnginePropRunStatusFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || runStatusAdded

	err, runtimeIdAdded := retrievePipelineEnginePropRuntimeIDFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || runtimeIdAdded

	err, statusAdded := retrievePipelineEnginePropStatusFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || statusAdded

	err, updateDateAdded := retrievePipelineEnginePropUpdateDateFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || updateDateAdded

	err, workspaceAdded := retrievePipelineEnginePropWorkspaceFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || workspaceAdded

	return nil, retAdded
}

func retrievePipelineEnginePropAvailabilityFlags(depth int, m *models.PipelineEngine, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	availabilityFlagName := fmt.Sprintf("%v.availability", cmdPrefix)
	if cmd.Flags().Changed(availabilityFlagName) {

		var availabilityFlagName string
		if cmdPrefix == "" {
			availabilityFlagName = "availability"
		} else {
			availabilityFlagName = fmt.Sprintf("%v.availability", cmdPrefix)
		}

		availabilityFlagValue, err := cmd.Flags().GetString(availabilityFlagName)
		if err != nil {
			return err, false
		}
		m.Availability = availabilityFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrievePipelineEnginePropCloudRunnerFlags(depth int, m *models.PipelineEngine, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	cloudRunnerFlagName := fmt.Sprintf("%v.cloudRunner", cmdPrefix)
	if cmd.Flags().Changed(cloudRunnerFlagName) {

		var cloudRunnerFlagName string
		if cmdPrefix == "" {
			cloudRunnerFlagName = "cloudRunner"
		} else {
			cloudRunnerFlagName = fmt.Sprintf("%v.cloudRunner", cmdPrefix)
		}

		cloudRunnerFlagValue, err := cmd.Flags().GetBool(cloudRunnerFlagName)
		if err != nil {
			return err, false
		}
		m.CloudRunner = cloudRunnerFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrievePipelineEnginePropCreateDateFlags(depth int, m *models.PipelineEngine, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrievePipelineEnginePropDescriptionFlags(depth int, m *models.PipelineEngine, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrievePipelineEnginePropIDFlags(depth int, m *models.PipelineEngine, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrievePipelineEnginePropManagedFlags(depth int, m *models.PipelineEngine, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	managedFlagName := fmt.Sprintf("%v.managed", cmdPrefix)
	if cmd.Flags().Changed(managedFlagName) {

		var managedFlagName string
		if cmdPrefix == "" {
			managedFlagName = "managed"
		} else {
			managedFlagName = fmt.Sprintf("%v.managed", cmdPrefix)
		}

		managedFlagValue, err := cmd.Flags().GetBool(managedFlagName)
		if err != nil {
			return err, false
		}
		m.Managed = managedFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrievePipelineEnginePropNameFlags(depth int, m *models.PipelineEngine, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrievePipelineEnginePropPreAuthorizedKeyFlags(depth int, m *models.PipelineEngine, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	preAuthorizedKeyFlagName := fmt.Sprintf("%v.preAuthorizedKey", cmdPrefix)
	if cmd.Flags().Changed(preAuthorizedKeyFlagName) {

		var preAuthorizedKeyFlagName string
		if cmdPrefix == "" {
			preAuthorizedKeyFlagName = "preAuthorizedKey"
		} else {
			preAuthorizedKeyFlagName = fmt.Sprintf("%v.preAuthorizedKey", cmdPrefix)
		}

		preAuthorizedKeyFlagValue, err := cmd.Flags().GetString(preAuthorizedKeyFlagName)
		if err != nil {
			return err, false
		}
		m.PreAuthorizedKey = preAuthorizedKeyFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrievePipelineEnginePropRunStatusFlags(depth int, m *models.PipelineEngine, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	runStatusFlagName := fmt.Sprintf("%v.runStatus", cmdPrefix)
	if cmd.Flags().Changed(runStatusFlagName) {

		var runStatusFlagName string
		if cmdPrefix == "" {
			runStatusFlagName = "runStatus"
		} else {
			runStatusFlagName = fmt.Sprintf("%v.runStatus", cmdPrefix)
		}

		runStatusFlagValue, err := cmd.Flags().GetString(runStatusFlagName)
		if err != nil {
			return err, false
		}
		m.RunStatus = runStatusFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrievePipelineEnginePropRuntimeIDFlags(depth int, m *models.PipelineEngine, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrievePipelineEnginePropStatusFlags(depth int, m *models.PipelineEngine, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

		statusFlagValue, err := cmd.Flags().GetString(statusFlagName)
		if err != nil {
			return err, false
		}
		m.Status = statusFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrievePipelineEnginePropUpdateDateFlags(depth int, m *models.PipelineEngine, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrievePipelineEnginePropWorkspaceFlags(depth int, m *models.PipelineEngine, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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
