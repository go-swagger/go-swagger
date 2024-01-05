// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"

	"github.com/go-swagger/go-swagger/fixtures/bugs/2969/codegen/client/schedules"
	"github.com/go-swagger/go-swagger/fixtures/bugs/2969/codegen/models"

	"github.com/go-openapi/swag"
	"github.com/spf13/cobra"
)

// makeOperationSchedulesAppendTriggerToScheduleMultipleCmd returns a command to handle operation appendTriggerToScheduleMultiple
func makeOperationSchedulesAppendTriggerToScheduleMultipleCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "appendTriggerToScheduleMultiple",
		Short: ``,
		RunE:  runOperationSchedulesAppendTriggerToScheduleMultiple,
	}

	if err := registerOperationSchedulesAppendTriggerToScheduleMultipleParamFlags(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}

// runOperationSchedulesAppendTriggerToScheduleMultiple uses cmd flags to call endpoint api
func runOperationSchedulesAppendTriggerToScheduleMultiple(cmd *cobra.Command, args []string) error {
	appCli, err := makeClient(cmd, args)
	if err != nil {
		return err
	}
	// retrieve flag values from cmd and fill params
	params := schedules.NewAppendTriggerToScheduleMultipleParams()
	if err, _ = retrieveOperationSchedulesAppendTriggerToScheduleMultipleBodyFlag(params, "", cmd); err != nil {
		return err
	}
	if err, _ = retrieveOperationSchedulesAppendTriggerToScheduleMultipleScheduleIDFlag(params, "", cmd); err != nil {
		return err
	}
	if dryRun {
		logDebugf("dry-run flag specified. Skip sending request.")
		return nil
	}
	// make request and then print result
	msgStr, err := parseOperationSchedulesAppendTriggerToScheduleMultipleResult(appCli.Schedules.AppendTriggerToScheduleMultiple(params, nil))
	if err != nil {
		return err
	}

	if !debug {
		fmt.Println(msgStr)
	}

	return nil
}

// registerOperationSchedulesAppendTriggerToScheduleMultipleParamFlags registers all flags needed to fill params
func registerOperationSchedulesAppendTriggerToScheduleMultipleParamFlags(cmd *cobra.Command) error {
	if err := registerOperationSchedulesAppendTriggerToScheduleMultipleBodyParamFlags("", cmd); err != nil {
		return err
	}
	if err := registerOperationSchedulesAppendTriggerToScheduleMultipleScheduleIDParamFlags("", cmd); err != nil {
		return err
	}
	return nil
}

func registerOperationSchedulesAppendTriggerToScheduleMultipleBodyParamFlags(cmdPrefix string, cmd *cobra.Command) error {

	var bodyFlagName string
	if cmdPrefix == "" {
		bodyFlagName = "body"
	} else {
		bodyFlagName = fmt.Sprintf("%v.body", cmdPrefix)
	}

	_ = cmd.PersistentFlags().String(bodyFlagName, "", "Optional json string for [body]. Trigger")

	// add flags for body
	if err := registerModelMultipleTriggerFlags(0, "multipleTrigger", cmd); err != nil {
		return err
	}

	return nil
}

func registerOperationSchedulesAppendTriggerToScheduleMultipleScheduleIDParamFlags(cmdPrefix string, cmd *cobra.Command) error {

	scheduleIdDescription := `Required. schedules id`

	var scheduleIdFlagName string
	if cmdPrefix == "" {
		scheduleIdFlagName = "scheduleId"
	} else {
		scheduleIdFlagName = fmt.Sprintf("%v.scheduleId", cmdPrefix)
	}

	var scheduleIdFlagDefault string

	_ = cmd.PersistentFlags().String(scheduleIdFlagName, scheduleIdFlagDefault, scheduleIdDescription)

	return nil
}

func retrieveOperationSchedulesAppendTriggerToScheduleMultipleBodyFlag(m *schedules.AppendTriggerToScheduleMultipleParams, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false
	if cmd.Flags().Changed("body") {
		// Read body string from cmd and unmarshal
		bodyValueStr, err := cmd.Flags().GetString("body")
		if err != nil {
			return err, false
		}

		bodyValue := models.MultipleTrigger{}
		if err := json.Unmarshal([]byte(bodyValueStr), &bodyValue); err != nil {
			return fmt.Errorf("cannot unmarshal body string in models.MultipleTrigger: %v", err), false
		}
		m.Body = &bodyValue
	}
	bodyValueModel := m.Body
	if swag.IsZero(bodyValueModel) {
		bodyValueModel = &models.MultipleTrigger{}
	}
	err, added := retrieveModelMultipleTriggerFlags(0, bodyValueModel, "multipleTrigger", cmd)
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

func retrieveOperationSchedulesAppendTriggerToScheduleMultipleScheduleIDFlag(m *schedules.AppendTriggerToScheduleMultipleParams, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false
	if cmd.Flags().Changed("scheduleId") {

		var scheduleIdFlagName string
		if cmdPrefix == "" {
			scheduleIdFlagName = "scheduleId"
		} else {
			scheduleIdFlagName = fmt.Sprintf("%v.scheduleId", cmdPrefix)
		}

		scheduleIdFlagValue, err := cmd.Flags().GetString(scheduleIdFlagName)
		if err != nil {
			return err, false
		}
		m.ScheduleID = scheduleIdFlagValue

	}

	return nil, retAdded
}

// parseOperationSchedulesAppendTriggerToScheduleMultipleResult parses request result and return the string content
func parseOperationSchedulesAppendTriggerToScheduleMultipleResult(resp0 *schedules.AppendTriggerToScheduleMultipleOK, respErr error) (string, error) {
	if respErr != nil {

		// Non schema case: warning appendTriggerToScheduleMultipleOK is not supported

		var iResp1 interface{} = respErr
		resp1, ok := iResp1.(*schedules.AppendTriggerToScheduleMultipleBadRequest)
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
		resp2, ok := iResp2.(*schedules.AppendTriggerToScheduleMultipleUnauthorized)
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
		resp3, ok := iResp3.(*schedules.AppendTriggerToScheduleMultipleForbidden)
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
		resp4, ok := iResp4.(*schedules.AppendTriggerToScheduleMultipleNotFound)
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
		resp5, ok := iResp5.(*schedules.AppendTriggerToScheduleMultipleConflict)
		if ok {
			if !swag.IsZero(resp5) && !swag.IsZero(resp5.Payload) {
				msgStr, err := json.Marshal(resp5.Payload)
				if err != nil {
					return "", err
				}
				return string(msgStr), nil
			}
		}

		var iResp6 interface{} = respErr
		resp6, ok := iResp6.(*schedules.AppendTriggerToScheduleMultiplePreconditionRequired)
		if ok {
			if !swag.IsZero(resp6) && !swag.IsZero(resp6.Payload) {
				msgStr, err := json.Marshal(resp6.Payload)
				if err != nil {
					return "", err
				}
				return string(msgStr), nil
			}
		}

		var iResp7 interface{} = respErr
		resp7, ok := iResp7.(*schedules.AppendTriggerToScheduleMultipleInternalServerError)
		if ok {
			if !swag.IsZero(resp7) && !swag.IsZero(resp7.Payload) {
				msgStr, err := json.Marshal(resp7.Payload)
				if err != nil {
					return "", err
				}
				return string(msgStr), nil
			}
		}

		return "", respErr
	}

	// warning: non schema response appendTriggerToScheduleMultipleOK is not supported by go-swagger cli yet.

	return "", nil
}
