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

// Schema cli for RuntimePromotionResult

// register flags to command
func registerModelRuntimePromotionResultFlags(depth int, cmdPrefix string, cmd *cobra.Command) error {

	if err := registerRuntimePromotionResultPropAnalyzeReport(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerRuntimePromotionResultPropID(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerRuntimePromotionResultPropJobType(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerRuntimePromotionResultPropName(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerRuntimePromotionResultPropPromotionReport(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerRuntimePromotionResultPropRunProfiles(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerRuntimePromotionResultPropTargetID(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerRuntimePromotionResultPropTargetVersion(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerRuntimePromotionResultPropType(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerRuntimePromotionResultPropUsedBy(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerRuntimePromotionResultPropVersion(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerRuntimePromotionResultPropVersions(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerRuntimePromotionResultPropWorkspaceID(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	return nil
}

func registerRuntimePromotionResultPropAnalyzeReport(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	var analyzeReportFlagName string
	if cmdPrefix == "" {
		analyzeReportFlagName = "analyzeReport"
	} else {
		analyzeReportFlagName = fmt.Sprintf("%v.analyzeReport", cmdPrefix)
	}

	if err := registerModelReportFlags(depth+1, analyzeReportFlagName, cmd); err != nil {
		return err
	}

	return nil
}

func registerRuntimePromotionResultPropID(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	idDescription := `Artifact ID`

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

func registerRuntimePromotionResultPropJobType(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	jobTypeDescription := `Job type`

	var jobTypeFlagName string
	if cmdPrefix == "" {
		jobTypeFlagName = "jobType"
	} else {
		jobTypeFlagName = fmt.Sprintf("%v.jobType", cmdPrefix)
	}

	var jobTypeFlagDefault string

	_ = cmd.PersistentFlags().String(jobTypeFlagName, jobTypeFlagDefault, jobTypeDescription)

	return nil
}

func registerRuntimePromotionResultPropName(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	nameDescription := `Artifact Name`

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

func registerRuntimePromotionResultPropPromotionReport(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	var promotionReportFlagName string
	if cmdPrefix == "" {
		promotionReportFlagName = "promotionReport"
	} else {
		promotionReportFlagName = fmt.Sprintf("%v.promotionReport", cmdPrefix)
	}

	if err := registerModelReportFlags(depth+1, promotionReportFlagName, cmd); err != nil {
		return err
	}

	return nil
}

func registerRuntimePromotionResultPropRunProfiles(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	// warning: runProfiles []*ArtifactPromotionResult array type is not supported by go-swagger cli yet

	return nil
}

func registerRuntimePromotionResultPropTargetID(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	targetIdDescription := `Target Artifact ID`

	var targetIdFlagName string
	if cmdPrefix == "" {
		targetIdFlagName = "targetId"
	} else {
		targetIdFlagName = fmt.Sprintf("%v.targetId", cmdPrefix)
	}

	var targetIdFlagDefault string

	_ = cmd.PersistentFlags().String(targetIdFlagName, targetIdFlagDefault, targetIdDescription)

	return nil
}

func registerRuntimePromotionResultPropTargetVersion(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	targetVersionDescription := `Artifact target version`

	var targetVersionFlagName string
	if cmdPrefix == "" {
		targetVersionFlagName = "targetVersion"
	} else {
		targetVersionFlagName = fmt.Sprintf("%v.targetVersion", cmdPrefix)
	}

	var targetVersionFlagDefault string

	_ = cmd.PersistentFlags().String(targetVersionFlagName, targetVersionFlagDefault, targetVersionDescription)

	return nil
}

func registerRuntimePromotionResultPropType(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	typeDescription := `Enum: ["WORKSPACE","PLAN","FLOW","ACTION","CONNECTION","RESOURCE","ENGINE","CLUSTER"]. Artifact Type`

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
			if err := json.Unmarshal([]byte(`["WORKSPACE","PLAN","FLOW","ACTION","CONNECTION","RESOURCE","ENGINE","CLUSTER"]`), &res); err != nil {
				panic(err)
			}
			return res, cobra.ShellCompDirectiveDefault
		}); err != nil {
		return err
	}

	return nil
}

func registerRuntimePromotionResultPropUsedBy(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	// warning: usedBy []*PromotionResultInfo array type is not supported by go-swagger cli yet

	return nil
}

func registerRuntimePromotionResultPropVersion(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	versionDescription := `Artifact version`

	var versionFlagName string
	if cmdPrefix == "" {
		versionFlagName = "version"
	} else {
		versionFlagName = fmt.Sprintf("%v.version", cmdPrefix)
	}

	var versionFlagDefault string

	_ = cmd.PersistentFlags().String(versionFlagName, versionFlagDefault, versionDescription)

	return nil
}

func registerRuntimePromotionResultPropVersions(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	// warning: versions []*ArtifactVersionPromotionResult array type is not supported by go-swagger cli yet

	return nil
}

func registerRuntimePromotionResultPropWorkspaceID(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	workspaceIdDescription := `Workspace id`

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
func retrieveModelRuntimePromotionResultFlags(depth int, m *models.RuntimePromotionResult, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false

	err, analyzeReportAdded := retrieveRuntimePromotionResultPropAnalyzeReportFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || analyzeReportAdded

	err, idAdded := retrieveRuntimePromotionResultPropIDFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || idAdded

	err, jobTypeAdded := retrieveRuntimePromotionResultPropJobTypeFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || jobTypeAdded

	err, nameAdded := retrieveRuntimePromotionResultPropNameFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || nameAdded

	err, promotionReportAdded := retrieveRuntimePromotionResultPropPromotionReportFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || promotionReportAdded

	err, runProfilesAdded := retrieveRuntimePromotionResultPropRunProfilesFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || runProfilesAdded

	err, targetIdAdded := retrieveRuntimePromotionResultPropTargetIDFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || targetIdAdded

	err, targetVersionAdded := retrieveRuntimePromotionResultPropTargetVersionFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || targetVersionAdded

	err, typeAdded := retrieveRuntimePromotionResultPropTypeFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || typeAdded

	err, usedByAdded := retrieveRuntimePromotionResultPropUsedByFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || usedByAdded

	err, versionAdded := retrieveRuntimePromotionResultPropVersionFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || versionAdded

	err, versionsAdded := retrieveRuntimePromotionResultPropVersionsFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || versionsAdded

	err, workspaceIdAdded := retrieveRuntimePromotionResultPropWorkspaceIDFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || workspaceIdAdded

	return nil, retAdded
}

func retrieveRuntimePromotionResultPropAnalyzeReportFlags(depth int, m *models.RuntimePromotionResult, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	analyzeReportFlagName := fmt.Sprintf("%v.analyzeReport", cmdPrefix)
	if cmd.Flags().Changed(analyzeReportFlagName) {
		// info: complex object analyzeReport Report is retrieved outside this Changed() block
	}
	analyzeReportFlagValue := m.AnalyzeReport
	if swag.IsZero(analyzeReportFlagValue) {
		analyzeReportFlagValue = &models.Report{}
	}

	err, analyzeReportAdded := retrieveModelReportFlags(depth+1, analyzeReportFlagValue, analyzeReportFlagName, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || analyzeReportAdded
	if analyzeReportAdded {
		m.AnalyzeReport = analyzeReportFlagValue
	}

	return nil, retAdded
}

func retrieveRuntimePromotionResultPropIDFlags(depth int, m *models.RuntimePromotionResult, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrieveRuntimePromotionResultPropJobTypeFlags(depth int, m *models.RuntimePromotionResult, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	jobTypeFlagName := fmt.Sprintf("%v.jobType", cmdPrefix)
	if cmd.Flags().Changed(jobTypeFlagName) {

		var jobTypeFlagName string
		if cmdPrefix == "" {
			jobTypeFlagName = "jobType"
		} else {
			jobTypeFlagName = fmt.Sprintf("%v.jobType", cmdPrefix)
		}

		jobTypeFlagValue, err := cmd.Flags().GetString(jobTypeFlagName)
		if err != nil {
			return err, false
		}
		m.JobType = jobTypeFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveRuntimePromotionResultPropNameFlags(depth int, m *models.RuntimePromotionResult, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrieveRuntimePromotionResultPropPromotionReportFlags(depth int, m *models.RuntimePromotionResult, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	promotionReportFlagName := fmt.Sprintf("%v.promotionReport", cmdPrefix)
	if cmd.Flags().Changed(promotionReportFlagName) {
		// info: complex object promotionReport Report is retrieved outside this Changed() block
	}
	promotionReportFlagValue := m.PromotionReport
	if swag.IsZero(promotionReportFlagValue) {
		promotionReportFlagValue = &models.Report{}
	}

	err, promotionReportAdded := retrieveModelReportFlags(depth+1, promotionReportFlagValue, promotionReportFlagName, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || promotionReportAdded
	if promotionReportAdded {
		m.PromotionReport = promotionReportFlagValue
	}

	return nil, retAdded
}

func retrieveRuntimePromotionResultPropRunProfilesFlags(depth int, m *models.RuntimePromotionResult, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	runProfilesFlagName := fmt.Sprintf("%v.runProfiles", cmdPrefix)
	if cmd.Flags().Changed(runProfilesFlagName) {
		// warning: runProfiles array type []*ArtifactPromotionResult is not supported by go-swagger cli yet
	}

	return nil, retAdded
}

func retrieveRuntimePromotionResultPropTargetIDFlags(depth int, m *models.RuntimePromotionResult, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	targetIdFlagName := fmt.Sprintf("%v.targetId", cmdPrefix)
	if cmd.Flags().Changed(targetIdFlagName) {

		var targetIdFlagName string
		if cmdPrefix == "" {
			targetIdFlagName = "targetId"
		} else {
			targetIdFlagName = fmt.Sprintf("%v.targetId", cmdPrefix)
		}

		targetIdFlagValue, err := cmd.Flags().GetString(targetIdFlagName)
		if err != nil {
			return err, false
		}
		m.TargetID = targetIdFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveRuntimePromotionResultPropTargetVersionFlags(depth int, m *models.RuntimePromotionResult, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	targetVersionFlagName := fmt.Sprintf("%v.targetVersion", cmdPrefix)
	if cmd.Flags().Changed(targetVersionFlagName) {

		var targetVersionFlagName string
		if cmdPrefix == "" {
			targetVersionFlagName = "targetVersion"
		} else {
			targetVersionFlagName = fmt.Sprintf("%v.targetVersion", cmdPrefix)
		}

		targetVersionFlagValue, err := cmd.Flags().GetString(targetVersionFlagName)
		if err != nil {
			return err, false
		}
		m.TargetVersion = targetVersionFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveRuntimePromotionResultPropTypeFlags(depth int, m *models.RuntimePromotionResult, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrieveRuntimePromotionResultPropUsedByFlags(depth int, m *models.RuntimePromotionResult, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	usedByFlagName := fmt.Sprintf("%v.usedBy", cmdPrefix)
	if cmd.Flags().Changed(usedByFlagName) {
		// warning: usedBy array type []*PromotionResultInfo is not supported by go-swagger cli yet
	}

	return nil, retAdded
}

func retrieveRuntimePromotionResultPropVersionFlags(depth int, m *models.RuntimePromotionResult, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	versionFlagName := fmt.Sprintf("%v.version", cmdPrefix)
	if cmd.Flags().Changed(versionFlagName) {

		var versionFlagName string
		if cmdPrefix == "" {
			versionFlagName = "version"
		} else {
			versionFlagName = fmt.Sprintf("%v.version", cmdPrefix)
		}

		versionFlagValue, err := cmd.Flags().GetString(versionFlagName)
		if err != nil {
			return err, false
		}
		m.Version = versionFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveRuntimePromotionResultPropVersionsFlags(depth int, m *models.RuntimePromotionResult, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	versionsFlagName := fmt.Sprintf("%v.versions", cmdPrefix)
	if cmd.Flags().Changed(versionsFlagName) {
		// warning: versions array type []*ArtifactVersionPromotionResult is not supported by go-swagger cli yet
	}

	return nil, retAdded
}

func retrieveRuntimePromotionResultPropWorkspaceIDFlags(depth int, m *models.RuntimePromotionResult, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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
