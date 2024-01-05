// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-swagger/go-swagger/fixtures/bugs/2969/codegen/models"
	"github.com/spf13/cobra"
)

// Schema cli for ExecutionLogs

// register flags to command
func registerModelExecutionLogsFlags(depth int, cmdPrefix string, cmd *cobra.Command) error {

	if err := registerExecutionLogsPropData(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerExecutionLogsPropNextIndex(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerExecutionLogsPropSize(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerExecutionLogsPropTotalSize(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerExecutionLogsPropTotalSizeFiltered(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	return nil
}

func registerExecutionLogsPropData(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	// warning: data []*CloudStorageLog array type is not supported by go-swagger cli yet

	return nil
}

func registerExecutionLogsPropNextIndex(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	nextIndexDescription := ``

	var nextIndexFlagName string
	if cmdPrefix == "" {
		nextIndexFlagName = "nextIndex"
	} else {
		nextIndexFlagName = fmt.Sprintf("%v.nextIndex", cmdPrefix)
	}

	var nextIndexFlagDefault int32

	_ = cmd.PersistentFlags().Int32(nextIndexFlagName, nextIndexFlagDefault, nextIndexDescription)

	return nil
}

func registerExecutionLogsPropSize(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	sizeDescription := ``

	var sizeFlagName string
	if cmdPrefix == "" {
		sizeFlagName = "size"
	} else {
		sizeFlagName = fmt.Sprintf("%v.size", cmdPrefix)
	}

	var sizeFlagDefault int32

	_ = cmd.PersistentFlags().Int32(sizeFlagName, sizeFlagDefault, sizeDescription)

	return nil
}

func registerExecutionLogsPropTotalSize(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	totalSizeDescription := ``

	var totalSizeFlagName string
	if cmdPrefix == "" {
		totalSizeFlagName = "totalSize"
	} else {
		totalSizeFlagName = fmt.Sprintf("%v.totalSize", cmdPrefix)
	}

	var totalSizeFlagDefault int32

	_ = cmd.PersistentFlags().Int32(totalSizeFlagName, totalSizeFlagDefault, totalSizeDescription)

	return nil
}

func registerExecutionLogsPropTotalSizeFiltered(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	totalSizeFilteredDescription := ``

	var totalSizeFilteredFlagName string
	if cmdPrefix == "" {
		totalSizeFilteredFlagName = "totalSizeFiltered"
	} else {
		totalSizeFilteredFlagName = fmt.Sprintf("%v.totalSizeFiltered", cmdPrefix)
	}

	var totalSizeFilteredFlagDefault int32

	_ = cmd.PersistentFlags().Int32(totalSizeFilteredFlagName, totalSizeFilteredFlagDefault, totalSizeFilteredDescription)

	return nil
}

// retrieve flags from commands, and set value in model. Return true if any flag is passed by user to fill model field.
func retrieveModelExecutionLogsFlags(depth int, m *models.ExecutionLogs, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false

	err, dataAdded := retrieveExecutionLogsPropDataFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || dataAdded

	err, nextIndexAdded := retrieveExecutionLogsPropNextIndexFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || nextIndexAdded

	err, sizeAdded := retrieveExecutionLogsPropSizeFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || sizeAdded

	err, totalSizeAdded := retrieveExecutionLogsPropTotalSizeFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || totalSizeAdded

	err, totalSizeFilteredAdded := retrieveExecutionLogsPropTotalSizeFilteredFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || totalSizeFilteredAdded

	return nil, retAdded
}

func retrieveExecutionLogsPropDataFlags(depth int, m *models.ExecutionLogs, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	dataFlagName := fmt.Sprintf("%v.data", cmdPrefix)
	if cmd.Flags().Changed(dataFlagName) {
		// warning: data array type []*CloudStorageLog is not supported by go-swagger cli yet
	}

	return nil, retAdded
}

func retrieveExecutionLogsPropNextIndexFlags(depth int, m *models.ExecutionLogs, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	nextIndexFlagName := fmt.Sprintf("%v.nextIndex", cmdPrefix)
	if cmd.Flags().Changed(nextIndexFlagName) {

		var nextIndexFlagName string
		if cmdPrefix == "" {
			nextIndexFlagName = "nextIndex"
		} else {
			nextIndexFlagName = fmt.Sprintf("%v.nextIndex", cmdPrefix)
		}

		nextIndexFlagValue, err := cmd.Flags().GetInt32(nextIndexFlagName)
		if err != nil {
			return err, false
		}
		m.NextIndex = nextIndexFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveExecutionLogsPropSizeFlags(depth int, m *models.ExecutionLogs, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	sizeFlagName := fmt.Sprintf("%v.size", cmdPrefix)
	if cmd.Flags().Changed(sizeFlagName) {

		var sizeFlagName string
		if cmdPrefix == "" {
			sizeFlagName = "size"
		} else {
			sizeFlagName = fmt.Sprintf("%v.size", cmdPrefix)
		}

		sizeFlagValue, err := cmd.Flags().GetInt32(sizeFlagName)
		if err != nil {
			return err, false
		}
		m.Size = sizeFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveExecutionLogsPropTotalSizeFlags(depth int, m *models.ExecutionLogs, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	totalSizeFlagName := fmt.Sprintf("%v.totalSize", cmdPrefix)
	if cmd.Flags().Changed(totalSizeFlagName) {

		var totalSizeFlagName string
		if cmdPrefix == "" {
			totalSizeFlagName = "totalSize"
		} else {
			totalSizeFlagName = fmt.Sprintf("%v.totalSize", cmdPrefix)
		}

		totalSizeFlagValue, err := cmd.Flags().GetInt32(totalSizeFlagName)
		if err != nil {
			return err, false
		}
		m.TotalSize = totalSizeFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveExecutionLogsPropTotalSizeFilteredFlags(depth int, m *models.ExecutionLogs, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	totalSizeFilteredFlagName := fmt.Sprintf("%v.totalSizeFiltered", cmdPrefix)
	if cmd.Flags().Changed(totalSizeFilteredFlagName) {

		var totalSizeFilteredFlagName string
		if cmdPrefix == "" {
			totalSizeFilteredFlagName = "totalSizeFiltered"
		} else {
			totalSizeFilteredFlagName = fmt.Sprintf("%v.totalSizeFiltered", cmdPrefix)
		}

		totalSizeFilteredFlagValue, err := cmd.Flags().GetInt32(totalSizeFilteredFlagName)
		if err != nil {
			return err, false
		}
		m.TotalSizeFiltered = totalSizeFilteredFlagValue

		retAdded = true
	}

	return nil, retAdded
}
