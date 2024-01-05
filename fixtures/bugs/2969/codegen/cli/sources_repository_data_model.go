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

// Schema cli for SourcesRepositoryData

// register flags to command
func registerModelSourcesRepositoryDataFlags(depth int, cmdPrefix string, cmd *cobra.Command) error {

	if err := registerSourcesRepositoryDataPropBranch(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerSourcesRepositoryDataPropCommit(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerSourcesRepositoryDataPropProject(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	return nil
}

func registerSourcesRepositoryDataPropBranch(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	branchDescription := `repository branch`

	var branchFlagName string
	if cmdPrefix == "" {
		branchFlagName = "branch"
	} else {
		branchFlagName = fmt.Sprintf("%v.branch", cmdPrefix)
	}

	var branchFlagDefault string

	_ = cmd.PersistentFlags().String(branchFlagName, branchFlagDefault, branchDescription)

	return nil
}

func registerSourcesRepositoryDataPropCommit(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	var commitFlagName string
	if cmdPrefix == "" {
		commitFlagName = "commit"
	} else {
		commitFlagName = fmt.Sprintf("%v.commit", cmdPrefix)
	}

	if err := registerModelRepositoryCommitDataFlags(depth+1, commitFlagName, cmd); err != nil {
		return err
	}

	return nil
}

func registerSourcesRepositoryDataPropProject(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	projectDescription := `Required. project name`

	var projectFlagName string
	if cmdPrefix == "" {
		projectFlagName = "project"
	} else {
		projectFlagName = fmt.Sprintf("%v.project", cmdPrefix)
	}

	var projectFlagDefault string

	_ = cmd.PersistentFlags().String(projectFlagName, projectFlagDefault, projectDescription)

	return nil
}

// retrieve flags from commands, and set value in model. Return true if any flag is passed by user to fill model field.
func retrieveModelSourcesRepositoryDataFlags(depth int, m *models.SourcesRepositoryData, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false

	err, branchAdded := retrieveSourcesRepositoryDataPropBranchFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || branchAdded

	err, commitAdded := retrieveSourcesRepositoryDataPropCommitFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || commitAdded

	err, projectAdded := retrieveSourcesRepositoryDataPropProjectFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || projectAdded

	return nil, retAdded
}

func retrieveSourcesRepositoryDataPropBranchFlags(depth int, m *models.SourcesRepositoryData, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	branchFlagName := fmt.Sprintf("%v.branch", cmdPrefix)
	if cmd.Flags().Changed(branchFlagName) {

		var branchFlagName string
		if cmdPrefix == "" {
			branchFlagName = "branch"
		} else {
			branchFlagName = fmt.Sprintf("%v.branch", cmdPrefix)
		}

		branchFlagValue, err := cmd.Flags().GetString(branchFlagName)
		if err != nil {
			return err, false
		}
		m.Branch = branchFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveSourcesRepositoryDataPropCommitFlags(depth int, m *models.SourcesRepositoryData, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	commitFlagName := fmt.Sprintf("%v.commit", cmdPrefix)
	if cmd.Flags().Changed(commitFlagName) {
		// info: complex object commit RepositoryCommitData is retrieved outside this Changed() block
	}
	commitFlagValue := m.Commit
	if swag.IsZero(commitFlagValue) {
		commitFlagValue = &models.RepositoryCommitData{}
	}

	err, commitAdded := retrieveModelRepositoryCommitDataFlags(depth+1, commitFlagValue, commitFlagName, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || commitAdded
	if commitAdded {
		m.Commit = commitFlagValue
	}

	return nil, retAdded
}

func retrieveSourcesRepositoryDataPropProjectFlags(depth int, m *models.SourcesRepositoryData, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	projectFlagName := fmt.Sprintf("%v.project", cmdPrefix)
	if cmd.Flags().Changed(projectFlagName) {

		var projectFlagName string
		if cmdPrefix == "" {
			projectFlagName = "project"
		} else {
			projectFlagName = fmt.Sprintf("%v.project", cmdPrefix)
		}

		projectFlagValue, err := cmd.Flags().GetString(projectFlagName)
		if err != nil {
			return err, false
		}
		m.Project = &projectFlagValue

		retAdded = true
	}

	return nil, retAdded
}
