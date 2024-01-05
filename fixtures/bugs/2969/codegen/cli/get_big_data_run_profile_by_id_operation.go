// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"

	"github.com/go-swagger/go-swagger/fixtures/bugs/2969/codegen/client/runtime_pipeline_engines_run_profiles_big_data"

	"github.com/go-openapi/swag"
	"github.com/spf13/cobra"
)

// makeOperationRuntimePipelineEnginesRunProfilesBigDataGetBigDataRunProfileByIDCmd returns a command to handle operation getBigDataRunProfileById
func makeOperationRuntimePipelineEnginesRunProfilesBigDataGetBigDataRunProfileByIDCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "getBigDataRunProfileById",
		Short: ``,
		RunE:  runOperationRuntimePipelineEnginesRunProfilesBigDataGetBigDataRunProfileByID,
	}

	if err := registerOperationRuntimePipelineEnginesRunProfilesBigDataGetBigDataRunProfileByIDParamFlags(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}

// runOperationRuntimePipelineEnginesRunProfilesBigDataGetBigDataRunProfileByID uses cmd flags to call endpoint api
func runOperationRuntimePipelineEnginesRunProfilesBigDataGetBigDataRunProfileByID(cmd *cobra.Command, args []string) error {
	appCli, err := makeClient(cmd, args)
	if err != nil {
		return err
	}
	// retrieve flag values from cmd and fill params
	params := runtime_pipeline_engines_run_profiles_big_data.NewGetBigDataRunProfileByIDParams()
	if err, _ = retrieveOperationRuntimePipelineEnginesRunProfilesBigDataGetBigDataRunProfileByIDEngineIDFlag(params, "", cmd); err != nil {
		return err
	}
	if err, _ = retrieveOperationRuntimePipelineEnginesRunProfilesBigDataGetBigDataRunProfileByIDRunProfileIDFlag(params, "", cmd); err != nil {
		return err
	}
	if dryRun {
		logDebugf("dry-run flag specified. Skip sending request.")
		return nil
	}
	// make request and then print result
	msgStr, err := parseOperationRuntimePipelineEnginesRunProfilesBigDataGetBigDataRunProfileByIDResult(appCli.RuntimePipelineEnginesRunProfilesBigData.GetBigDataRunProfileByID(params, nil))
	if err != nil {
		return err
	}

	if !debug {
		fmt.Println(msgStr)
	}

	return nil
}

// registerOperationRuntimePipelineEnginesRunProfilesBigDataGetBigDataRunProfileByIDParamFlags registers all flags needed to fill params
func registerOperationRuntimePipelineEnginesRunProfilesBigDataGetBigDataRunProfileByIDParamFlags(cmd *cobra.Command) error {
	if err := registerOperationRuntimePipelineEnginesRunProfilesBigDataGetBigDataRunProfileByIDEngineIDParamFlags("", cmd); err != nil {
		return err
	}
	if err := registerOperationRuntimePipelineEnginesRunProfilesBigDataGetBigDataRunProfileByIDRunProfileIDParamFlags("", cmd); err != nil {
		return err
	}
	return nil
}

func registerOperationRuntimePipelineEnginesRunProfilesBigDataGetBigDataRunProfileByIDEngineIDParamFlags(cmdPrefix string, cmd *cobra.Command) error {

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

func registerOperationRuntimePipelineEnginesRunProfilesBigDataGetBigDataRunProfileByIDRunProfileIDParamFlags(cmdPrefix string, cmd *cobra.Command) error {

	runProfileIdDescription := `Required. run profile id`

	var runProfileIdFlagName string
	if cmdPrefix == "" {
		runProfileIdFlagName = "runProfileId"
	} else {
		runProfileIdFlagName = fmt.Sprintf("%v.runProfileId", cmdPrefix)
	}

	var runProfileIdFlagDefault string

	_ = cmd.PersistentFlags().String(runProfileIdFlagName, runProfileIdFlagDefault, runProfileIdDescription)

	return nil
}

func retrieveOperationRuntimePipelineEnginesRunProfilesBigDataGetBigDataRunProfileByIDEngineIDFlag(m *runtime_pipeline_engines_run_profiles_big_data.GetBigDataRunProfileByIDParams, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

func retrieveOperationRuntimePipelineEnginesRunProfilesBigDataGetBigDataRunProfileByIDRunProfileIDFlag(m *runtime_pipeline_engines_run_profiles_big_data.GetBigDataRunProfileByIDParams, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false
	if cmd.Flags().Changed("runProfileId") {

		var runProfileIdFlagName string
		if cmdPrefix == "" {
			runProfileIdFlagName = "runProfileId"
		} else {
			runProfileIdFlagName = fmt.Sprintf("%v.runProfileId", cmdPrefix)
		}

		runProfileIdFlagValue, err := cmd.Flags().GetString(runProfileIdFlagName)
		if err != nil {
			return err, false
		}
		m.RunProfileID = runProfileIdFlagValue

	}

	return nil, retAdded
}

// parseOperationRuntimePipelineEnginesRunProfilesBigDataGetBigDataRunProfileByIDResult parses request result and return the string content
func parseOperationRuntimePipelineEnginesRunProfilesBigDataGetBigDataRunProfileByIDResult(resp0 *runtime_pipeline_engines_run_profiles_big_data.GetBigDataRunProfileByIDOK, respErr error) (string, error) {
	if respErr != nil {

		var iResp0 interface{} = respErr
		resp0, ok := iResp0.(*runtime_pipeline_engines_run_profiles_big_data.GetBigDataRunProfileByIDOK)
		if ok {
			if !swag.IsZero(resp0) && !swag.IsZero(resp0.Payload) {
				msgStr, err := json.Marshal(resp0.Payload)
				if err != nil {
					return "", err
				}
				return string(msgStr), nil
			}
		}

		var iResp1 interface{} = respErr
		resp1, ok := iResp1.(*runtime_pipeline_engines_run_profiles_big_data.GetBigDataRunProfileByIDBadRequest)
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
		resp2, ok := iResp2.(*runtime_pipeline_engines_run_profiles_big_data.GetBigDataRunProfileByIDUnauthorized)
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
		resp3, ok := iResp3.(*runtime_pipeline_engines_run_profiles_big_data.GetBigDataRunProfileByIDForbidden)
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
		resp4, ok := iResp4.(*runtime_pipeline_engines_run_profiles_big_data.GetBigDataRunProfileByIDNotFound)
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
		resp5, ok := iResp5.(*runtime_pipeline_engines_run_profiles_big_data.GetBigDataRunProfileByIDInternalServerError)
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

	if !swag.IsZero(resp0) && !swag.IsZero(resp0.Payload) {
		msgStr, err := json.Marshal(resp0.Payload)
		if err != nil {
			return "", err
		}
		return string(msgStr), nil
	}

	return "", nil
}
