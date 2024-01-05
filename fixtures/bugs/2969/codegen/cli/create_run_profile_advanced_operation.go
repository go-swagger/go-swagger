// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"

	"github.com/go-swagger/go-swagger/fixtures/bugs/2969/codegen/client/runtime_pipeline_engines_run_profiles_advanced"
	"github.com/go-swagger/go-swagger/fixtures/bugs/2969/codegen/models"

	"github.com/go-openapi/swag"
	"github.com/spf13/cobra"
)

// makeOperationRuntimePipelineEnginesRunProfilesAdvancedCreateRunProfileAdvancedCmd returns a command to handle operation createRunProfileAdvanced
func makeOperationRuntimePipelineEnginesRunProfilesAdvancedCreateRunProfileAdvancedCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "createRunProfileAdvanced",
		Short: ``,
		RunE:  runOperationRuntimePipelineEnginesRunProfilesAdvancedCreateRunProfileAdvanced,
	}

	if err := registerOperationRuntimePipelineEnginesRunProfilesAdvancedCreateRunProfileAdvancedParamFlags(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}

// runOperationRuntimePipelineEnginesRunProfilesAdvancedCreateRunProfileAdvanced uses cmd flags to call endpoint api
func runOperationRuntimePipelineEnginesRunProfilesAdvancedCreateRunProfileAdvanced(cmd *cobra.Command, args []string) error {
	appCli, err := makeClient(cmd, args)
	if err != nil {
		return err
	}
	// retrieve flag values from cmd and fill params
	params := runtime_pipeline_engines_run_profiles_advanced.NewCreateRunProfileAdvancedParams()
	if err, _ = retrieveOperationRuntimePipelineEnginesRunProfilesAdvancedCreateRunProfileAdvancedBodyFlag(params, "", cmd); err != nil {
		return err
	}
	if err, _ = retrieveOperationRuntimePipelineEnginesRunProfilesAdvancedCreateRunProfileAdvancedEngineIDFlag(params, "", cmd); err != nil {
		return err
	}
	if dryRun {
		logDebugf("dry-run flag specified. Skip sending request.")
		return nil
	}
	// make request and then print result
	msgStr, err := parseOperationRuntimePipelineEnginesRunProfilesAdvancedCreateRunProfileAdvancedResult(appCli.RuntimePipelineEnginesRunProfilesAdvanced.CreateRunProfileAdvanced(params, nil))
	if err != nil {
		return err
	}

	if !debug {
		fmt.Println(msgStr)
	}

	return nil
}

// registerOperationRuntimePipelineEnginesRunProfilesAdvancedCreateRunProfileAdvancedParamFlags registers all flags needed to fill params
func registerOperationRuntimePipelineEnginesRunProfilesAdvancedCreateRunProfileAdvancedParamFlags(cmd *cobra.Command) error {
	if err := registerOperationRuntimePipelineEnginesRunProfilesAdvancedCreateRunProfileAdvancedBodyParamFlags("", cmd); err != nil {
		return err
	}
	if err := registerOperationRuntimePipelineEnginesRunProfilesAdvancedCreateRunProfileAdvancedEngineIDParamFlags("", cmd); err != nil {
		return err
	}
	return nil
}

func registerOperationRuntimePipelineEnginesRunProfilesAdvancedCreateRunProfileAdvancedBodyParamFlags(cmdPrefix string, cmd *cobra.Command) error {

	var bodyFlagName string
	if cmdPrefix == "" {
		bodyFlagName = "body"
	} else {
		bodyFlagName = fmt.Sprintf("%v.body", cmdPrefix)
	}

	_ = cmd.PersistentFlags().String(bodyFlagName, "", "Optional json string for [body]. ")

	// add flags for body
	if err := registerModelAdvancedRunProfileFlags(0, "advancedRunProfile", cmd); err != nil {
		return err
	}

	return nil
}

func registerOperationRuntimePipelineEnginesRunProfilesAdvancedCreateRunProfileAdvancedEngineIDParamFlags(cmdPrefix string, cmd *cobra.Command) error {

	engineIdDescription := `Required. engine id`

	var engineIdFlagName string
	if cmdPrefix == "" {
		engineIdFlagName = "engineId"
	} else {
		engineIdFlagName = fmt.Sprintf("%v.engineId", cmdPrefix)
	}

	var engineIdFlagDefault string

	_ = cmd.PersistentFlags().String(engineIdFlagName, engineIdFlagDefault, engineIdDescription)

	return nil
}

func retrieveOperationRuntimePipelineEnginesRunProfilesAdvancedCreateRunProfileAdvancedBodyFlag(m *runtime_pipeline_engines_run_profiles_advanced.CreateRunProfileAdvancedParams, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false
	if cmd.Flags().Changed("body") {
		// Read body string from cmd and unmarshal
		bodyValueStr, err := cmd.Flags().GetString("body")
		if err != nil {
			return err, false
		}

		bodyValue := models.AdvancedRunProfile{}
		if err := json.Unmarshal([]byte(bodyValueStr), &bodyValue); err != nil {
			return fmt.Errorf("cannot unmarshal body string in models.AdvancedRunProfile: %v", err), false
		}
		m.Body = &bodyValue
	}
	bodyValueModel := m.Body
	if swag.IsZero(bodyValueModel) {
		bodyValueModel = &models.AdvancedRunProfile{}
	}
	err, added := retrieveModelAdvancedRunProfileFlags(0, bodyValueModel, "advancedRunProfile", cmd)
	if err != nil {
		return err, false
	}
	if added {
		m.Body = bodyValueModel
	}

	if dryRun && debug {
		bodyValueDebugBytes, err := json.Marshal(m.Body)
		if err != nil {
			return err, false
		}
		logDebugf("Body dry-run payload: %v", string(bodyValueDebugBytes))
	}

	retAdded = retAdded || added

	return nil, retAdded
}

func retrieveOperationRuntimePipelineEnginesRunProfilesAdvancedCreateRunProfileAdvancedEngineIDFlag(m *runtime_pipeline_engines_run_profiles_advanced.CreateRunProfileAdvancedParams, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false
	if cmd.Flags().Changed("engineId") {

		var engineIdFlagName string
		if cmdPrefix == "" {
			engineIdFlagName = "engineId"
		} else {
			engineIdFlagName = fmt.Sprintf("%v.engineId", cmdPrefix)
		}

		engineIdFlagValue, err := cmd.Flags().GetString(engineIdFlagName)
		if err != nil {
			return err, false
		}
		m.EngineID = engineIdFlagValue

	}

	return nil, retAdded
}

// parseOperationRuntimePipelineEnginesRunProfilesAdvancedCreateRunProfileAdvancedResult parses request result and return the string content
func parseOperationRuntimePipelineEnginesRunProfilesAdvancedCreateRunProfileAdvancedResult(resp0 *runtime_pipeline_engines_run_profiles_advanced.CreateRunProfileAdvancedCreated, respErr error) (string, error) {
	if respErr != nil {

		// Non schema case: warning createRunProfileAdvancedCreated is not supported

		var iResp1 interface{} = respErr
		resp1, ok := iResp1.(*runtime_pipeline_engines_run_profiles_advanced.CreateRunProfileAdvancedBadRequest)
		if ok {
			if !swag.IsZero(resp1) && !swag.IsZero(resp1.Payload) {
				msgStr, err := json.Marshal(resp1.Payload)
				if err != nil {
					return "", err
				}
				return string(msgStr), nil
			}
		}

		var iResp2 interface{} = respErr
		resp2, ok := iResp2.(*runtime_pipeline_engines_run_profiles_advanced.CreateRunProfileAdvancedUnauthorized)
		if ok {
			if !swag.IsZero(resp2) && !swag.IsZero(resp2.Payload) {
				msgStr, err := json.Marshal(resp2.Payload)
				if err != nil {
					return "", err
				}
				return string(msgStr), nil
			}
		}

		var iResp3 interface{} = respErr
		resp3, ok := iResp3.(*runtime_pipeline_engines_run_profiles_advanced.CreateRunProfileAdvancedForbidden)
		if ok {
			if !swag.IsZero(resp3) && !swag.IsZero(resp3.Payload) {
				msgStr, err := json.Marshal(resp3.Payload)
				if err != nil {
					return "", err
				}
				return string(msgStr), nil
			}
		}

		var iResp4 interface{} = respErr
		resp4, ok := iResp4.(*runtime_pipeline_engines_run_profiles_advanced.CreateRunProfileAdvancedNotFound)
		if ok {
			if !swag.IsZero(resp4) && !swag.IsZero(resp4.Payload) {
				msgStr, err := json.Marshal(resp4.Payload)
				if err != nil {
					return "", err
				}
				return string(msgStr), nil
			}
		}

		var iResp5 interface{} = respErr
		resp5, ok := iResp5.(*runtime_pipeline_engines_run_profiles_advanced.CreateRunProfileAdvancedInternalServerError)
		if ok {
			if !swag.IsZero(resp5) && !swag.IsZero(resp5.Payload) {
				msgStr, err := json.Marshal(resp5.Payload)
				if err != nil {
					return "", err
				}
				return string(msgStr), nil
			}
		}

		return "", respErr
	}

	// warning: non schema response createRunProfileAdvancedCreated is not supported by go-swagger cli yet.

	return "", nil
}
