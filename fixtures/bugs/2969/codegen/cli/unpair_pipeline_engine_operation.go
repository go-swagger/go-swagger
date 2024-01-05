// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"

	"github.com/go-swagger/go-swagger/fixtures/bugs/2969/codegen/client/runtime_pipeline_engines"

	"github.com/go-openapi/swag"
	"github.com/spf13/cobra"
)

// makeOperationRuntimePipelineEnginesUnpairPipelineEngineCmd returns a command to handle operation unpairPipelineEngine
func makeOperationRuntimePipelineEnginesUnpairPipelineEngineCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "unpairPipelineEngine",
		Short: `Unpair Pipeline Engine`,
		RunE:  runOperationRuntimePipelineEnginesUnpairPipelineEngine,
	}

	if err := registerOperationRuntimePipelineEnginesUnpairPipelineEngineParamFlags(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}

// runOperationRuntimePipelineEnginesUnpairPipelineEngine uses cmd flags to call endpoint api
func runOperationRuntimePipelineEnginesUnpairPipelineEngine(cmd *cobra.Command, args []string) error {
	appCli, err := makeClient(cmd, args)
	if err != nil {
		return err
	}
	// retrieve flag values from cmd and fill params
	params := runtime_pipeline_engines.NewUnpairPipelineEngineParams()
	if err, _ = retrieveOperationRuntimePipelineEnginesUnpairPipelineEngineEngineIDFlag(params, "", cmd); err != nil {
		return err
	}
	if dryRun {
		logDebugf("dry-run flag specified. Skip sending request.")
		return nil
	}
	// make request and then print result
	msgStr, err := parseOperationRuntimePipelineEnginesUnpairPipelineEngineResult(appCli.RuntimePipelineEngines.UnpairPipelineEngine(params, nil))
	if err != nil {
		return err
	}

	if !debug {
		fmt.Println(msgStr)
	}

	return nil
}

// registerOperationRuntimePipelineEnginesUnpairPipelineEngineParamFlags registers all flags needed to fill params
func registerOperationRuntimePipelineEnginesUnpairPipelineEngineParamFlags(cmd *cobra.Command) error {
	if err := registerOperationRuntimePipelineEnginesUnpairPipelineEngineEngineIDParamFlags("", cmd); err != nil {
		return err
	}
	return nil
}

func registerOperationRuntimePipelineEnginesUnpairPipelineEngineEngineIDParamFlags(cmdPrefix string, cmd *cobra.Command) error {

	engineIdDescription := `Required. pipeline engine id`

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

func retrieveOperationRuntimePipelineEnginesUnpairPipelineEngineEngineIDFlag(m *runtime_pipeline_engines.UnpairPipelineEngineParams, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

// parseOperationRuntimePipelineEnginesUnpairPipelineEngineResult parses request result and return the string content
func parseOperationRuntimePipelineEnginesUnpairPipelineEngineResult(resp0 *runtime_pipeline_engines.UnpairPipelineEngineOK, respErr error) (string, error) {
	if respErr != nil {

		var iResp0 interface{} = respErr
		resp0, ok := iResp0.(*runtime_pipeline_engines.UnpairPipelineEngineOK)
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
		resp1, ok := iResp1.(*runtime_pipeline_engines.UnpairPipelineEngineBadRequest)
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
		resp2, ok := iResp2.(*runtime_pipeline_engines.UnpairPipelineEngineUnauthorized)
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
		resp3, ok := iResp3.(*runtime_pipeline_engines.UnpairPipelineEngineForbidden)
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
		resp4, ok := iResp4.(*runtime_pipeline_engines.UnpairPipelineEngineNotFound)
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
		resp5, ok := iResp5.(*runtime_pipeline_engines.UnpairPipelineEngineInternalServerError)
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
		msgStr := fmt.Sprintf("%v", resp0.Payload)
		return string(msgStr), nil
	}

	return "", nil
}
