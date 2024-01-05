// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-swagger/go-swagger/fixtures/bugs/2969/codegen/models"
	"github.com/spf13/cobra"
)

// Schema cli for PageArtifact

// register flags to command
func registerModelPageArtifactFlags(depth int, cmdPrefix string, cmd *cobra.Command) error {

	if err := registerPageArtifactPropItems(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPageArtifactPropLimit(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPageArtifactPropOffset(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerPageArtifactPropTotal(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	return nil
}

func registerPageArtifactPropItems(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	// warning: items []*Artifact array type is not supported by go-swagger cli yet

	return nil
}

func registerPageArtifactPropLimit(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	limitDescription := ``

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

func registerPageArtifactPropOffset(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	offsetDescription := ``

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

func registerPageArtifactPropTotal(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	totalDescription := ``

	var totalFlagName string
	if cmdPrefix == "" {
		totalFlagName = "total"
	} else {
		totalFlagName = fmt.Sprintf("%v.total", cmdPrefix)
	}

	var totalFlagDefault int32

	_ = cmd.PersistentFlags().Int32(totalFlagName, totalFlagDefault, totalDescription)

	return nil
}

// retrieve flags from commands, and set value in model. Return true if any flag is passed by user to fill model field.
func retrieveModelPageArtifactFlags(depth int, m *models.PageArtifact, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false

	err, itemsAdded := retrievePageArtifactPropItemsFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || itemsAdded

	err, limitAdded := retrievePageArtifactPropLimitFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || limitAdded

	err, offsetAdded := retrievePageArtifactPropOffsetFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || offsetAdded

	err, totalAdded := retrievePageArtifactPropTotalFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || totalAdded

	return nil, retAdded
}

func retrievePageArtifactPropItemsFlags(depth int, m *models.PageArtifact, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	itemsFlagName := fmt.Sprintf("%v.items", cmdPrefix)
	if cmd.Flags().Changed(itemsFlagName) {
		// warning: items array type []*Artifact is not supported by go-swagger cli yet
	}

	return nil, retAdded
}

func retrievePageArtifactPropLimitFlags(depth int, m *models.PageArtifact, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrievePageArtifactPropOffsetFlags(depth int, m *models.PageArtifact, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrievePageArtifactPropTotalFlags(depth int, m *models.PageArtifact, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	totalFlagName := fmt.Sprintf("%v.total", cmdPrefix)
	if cmd.Flags().Changed(totalFlagName) {

		var totalFlagName string
		if cmdPrefix == "" {
			totalFlagName = "total"
		} else {
			totalFlagName = fmt.Sprintf("%v.total", cmdPrefix)
		}

		totalFlagValue, err := cmd.Flags().GetInt32(totalFlagName)
		if err != nil {
			return err, false
		}
		m.Total = totalFlagValue

		retAdded = true
	}

	return nil, retAdded
}
