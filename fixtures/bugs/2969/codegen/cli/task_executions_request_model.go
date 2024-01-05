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

// Schema cli for TaskExecutionsRequest

// register flags to command
func registerModelTaskExecutionsRequestFlags(depth int, cmdPrefix string, cmd *cobra.Command) error {

	if err := registerTaskExecutionsRequestPropEnvironmentID(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerTaskExecutionsRequestPropFrom(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerTaskExecutionsRequestPropLastDays(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerTaskExecutionsRequestPropLimit(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerTaskExecutionsRequestPropOffset(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerTaskExecutionsRequestPropStatus(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerTaskExecutionsRequestPropTags(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerTaskExecutionsRequestPropTo(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerTaskExecutionsRequestPropWorkspaceID(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	return nil
}

func registerTaskExecutionsRequestPropEnvironmentID(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	environmentIdDescription := `tasks environment id`

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

func registerTaskExecutionsRequestPropFrom(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	fromDescription := `from date time (milliseconds)`

	var fromFlagName string
	if cmdPrefix == "" {
		fromFlagName = "from"
	} else {
		fromFlagName = fmt.Sprintf("%v.from", cmdPrefix)
	}

	var fromFlagDefault int64

	_ = cmd.PersistentFlags().Int64(fromFlagName, fromFlagDefault, fromDescription)

	return nil
}

func registerTaskExecutionsRequestPropLastDays(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	lastDaysDescription := `lastDays`

	var lastDaysFlagName string
	if cmdPrefix == "" {
		lastDaysFlagName = "lastDays"
	} else {
		lastDaysFlagName = fmt.Sprintf("%v.lastDays", cmdPrefix)
	}

	var lastDaysFlagDefault int32

	_ = cmd.PersistentFlags().Int32(lastDaysFlagName, lastDaysFlagDefault, lastDaysDescription)

	return nil
}

func registerTaskExecutionsRequestPropLimit(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	limitDescription := `tasks limit`

	var limitFlagName string
	if cmdPrefix == "" {
		limitFlagName = "limit"
	} else {
		limitFlagName = fmt.Sprintf("%v.limit", cmdPrefix)
	}

	var limitFlagDefault int32

	_ = cmd.PersistentFlags().Int32(limitFlagName, limitFlagDefault, limitDescription)

	return nil
}

func registerTaskExecutionsRequestPropOffset(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	offsetDescription := `tasks offset`

	var offsetFlagName string
	if cmdPrefix == "" {
		offsetFlagName = "offset"
	} else {
		offsetFlagName = fmt.Sprintf("%v.offset", cmdPrefix)
	}

	var offsetFlagDefault int32

	_ = cmd.PersistentFlags().Int32(offsetFlagName, offsetFlagDefault, offsetDescription)

	return nil
}

func registerTaskExecutionsRequestPropStatus(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	statusDescription := `Enum: ["dispatching","deploy_failed","executing","execution_successful","execution_rejected","execution_failed","terminated","terminated_timeout"]. tasks execution status`

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
			if err := json.Unmarshal([]byte(`["dispatching","deploy_failed","executing","execution_successful","execution_rejected","execution_failed","terminated","terminated_timeout"]`), &res); err != nil {
				panic(err)
			}
			return res, cobra.ShellCompDirectiveDefault
		}); err != nil {
		return err
	}

	return nil
}

func registerTaskExecutionsRequestPropTags(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	// warning: tags []string array type is not supported by go-swagger cli yet

	return nil
}

func registerTaskExecutionsRequestPropTo(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	toDescription := `to date time (milliseconds)`

	var toFlagName string
	if cmdPrefix == "" {
		toFlagName = "to"
	} else {
		toFlagName = fmt.Sprintf("%v.to", cmdPrefix)
	}

	var toFlagDefault int64

	_ = cmd.PersistentFlags().Int64(toFlagName, toFlagDefault, toDescription)

	return nil
}

func registerTaskExecutionsRequestPropWorkspaceID(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	workspaceIdDescription := `tasks workspace id`

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
func retrieveModelTaskExecutionsRequestFlags(depth int, m *models.TaskExecutionsRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false

	err, environmentIdAdded := retrieveTaskExecutionsRequestPropEnvironmentIDFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || environmentIdAdded

	err, fromAdded := retrieveTaskExecutionsRequestPropFromFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || fromAdded

	err, lastDaysAdded := retrieveTaskExecutionsRequestPropLastDaysFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || lastDaysAdded

	err, limitAdded := retrieveTaskExecutionsRequestPropLimitFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || limitAdded

	err, offsetAdded := retrieveTaskExecutionsRequestPropOffsetFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || offsetAdded

	err, statusAdded := retrieveTaskExecutionsRequestPropStatusFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || statusAdded

	err, tagsAdded := retrieveTaskExecutionsRequestPropTagsFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || tagsAdded

	err, toAdded := retrieveTaskExecutionsRequestPropToFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || toAdded

	err, workspaceIdAdded := retrieveTaskExecutionsRequestPropWorkspaceIDFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || workspaceIdAdded

	return nil, retAdded
}

func retrieveTaskExecutionsRequestPropEnvironmentIDFlags(depth int, m *models.TaskExecutionsRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrieveTaskExecutionsRequestPropFromFlags(depth int, m *models.TaskExecutionsRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	fromFlagName := fmt.Sprintf("%v.from", cmdPrefix)
	if cmd.Flags().Changed(fromFlagName) {

		var fromFlagName string
		if cmdPrefix == "" {
			fromFlagName = "from"
		} else {
			fromFlagName = fmt.Sprintf("%v.from", cmdPrefix)
		}

		fromFlagValue, err := cmd.Flags().GetInt64(fromFlagName)
		if err != nil {
			return err, false
		}
		m.From = fromFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveTaskExecutionsRequestPropLastDaysFlags(depth int, m *models.TaskExecutionsRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	lastDaysFlagName := fmt.Sprintf("%v.lastDays", cmdPrefix)
	if cmd.Flags().Changed(lastDaysFlagName) {

		var lastDaysFlagName string
		if cmdPrefix == "" {
			lastDaysFlagName = "lastDays"
		} else {
			lastDaysFlagName = fmt.Sprintf("%v.lastDays", cmdPrefix)
		}

		lastDaysFlagValue, err := cmd.Flags().GetInt32(lastDaysFlagName)
		if err != nil {
			return err, false
		}
		m.LastDays = lastDaysFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveTaskExecutionsRequestPropLimitFlags(depth int, m *models.TaskExecutionsRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	limitFlagName := fmt.Sprintf("%v.limit", cmdPrefix)
	if cmd.Flags().Changed(limitFlagName) {

		var limitFlagName string
		if cmdPrefix == "" {
			limitFlagName = "limit"
		} else {
			limitFlagName = fmt.Sprintf("%v.limit", cmdPrefix)
		}

		limitFlagValue, err := cmd.Flags().GetInt32(limitFlagName)
		if err != nil {
			return err, false
		}
		m.Limit = limitFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveTaskExecutionsRequestPropOffsetFlags(depth int, m *models.TaskExecutionsRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	offsetFlagName := fmt.Sprintf("%v.offset", cmdPrefix)
	if cmd.Flags().Changed(offsetFlagName) {

		var offsetFlagName string
		if cmdPrefix == "" {
			offsetFlagName = "offset"
		} else {
			offsetFlagName = fmt.Sprintf("%v.offset", cmdPrefix)
		}

		offsetFlagValue, err := cmd.Flags().GetInt32(offsetFlagName)
		if err != nil {
			return err, false
		}
		m.Offset = offsetFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveTaskExecutionsRequestPropStatusFlags(depth int, m *models.TaskExecutionsRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrieveTaskExecutionsRequestPropTagsFlags(depth int, m *models.TaskExecutionsRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrieveTaskExecutionsRequestPropToFlags(depth int, m *models.TaskExecutionsRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	toFlagName := fmt.Sprintf("%v.to", cmdPrefix)
	if cmd.Flags().Changed(toFlagName) {

		var toFlagName string
		if cmdPrefix == "" {
			toFlagName = "to"
		} else {
			toFlagName = fmt.Sprintf("%v.to", cmdPrefix)
		}

		toFlagValue, err := cmd.Flags().GetInt64(toFlagName)
		if err != nil {
			return err, false
		}
		m.To = toFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveTaskExecutionsRequestPropWorkspaceIDFlags(depth int, m *models.TaskExecutionsRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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
		m.WorkspaceID = workspaceIdFlagValue

		retAdded = true
	}

	return nil, retAdded
}
