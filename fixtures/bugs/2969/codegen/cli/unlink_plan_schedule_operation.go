// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"

	"github.com/go-swagger/go-swagger/fixtures/bugs/2969/codegen/client/plans_executables"

	"github.com/go-openapi/swag"
	"github.com/spf13/cobra"
)

// makeOperationPlansExecutablesUnlinkPlanScheduleCmd returns a command to handle operation unlinkPlanSchedule
func makeOperationPlansExecutablesUnlinkPlanScheduleCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "unlinkPlanSchedule",
		Short: `Unlink schedule from plan`,
		RunE:  runOperationPlansExecutablesUnlinkPlanSchedule,
	}

	if err := registerOperationPlansExecutablesUnlinkPlanScheduleParamFlags(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}

// runOperationPlansExecutablesUnlinkPlanSchedule uses cmd flags to call endpoint api
func runOperationPlansExecutablesUnlinkPlanSchedule(cmd *cobra.Command, args []string) error {
	appCli, err := makeClient(cmd, args)
	if err != nil {
		return err
	}
	// retrieve flag values from cmd and fill params
	params := plans_executables.NewUnlinkPlanScheduleParams()
	if err, _ = retrieveOperationPlansExecutablesUnlinkPlanSchedulePlanIDFlag(params, "", cmd); err != nil {
		return err
	}
	if dryRun {
		logDebugf("dry-run flag specified. Skip sending request.")
		return nil
	}
	// make request and then print result
	msgStr, err := parseOperationPlansExecutablesUnlinkPlanScheduleResult(appCli.PlansExecutables.UnlinkPlanSchedule(params, nil))
	if err != nil {
		return err
	}

	if !debug {
		fmt.Println(msgStr)
	}

	return nil
}

// registerOperationPlansExecutablesUnlinkPlanScheduleParamFlags registers all flags needed to fill params
func registerOperationPlansExecutablesUnlinkPlanScheduleParamFlags(cmd *cobra.Command) error {
	if err := registerOperationPlansExecutablesUnlinkPlanSchedulePlanIDParamFlags("", cmd); err != nil {
		return err
	}
	return nil
}

func registerOperationPlansExecutablesUnlinkPlanSchedulePlanIDParamFlags(cmdPrefix string, cmd *cobra.Command) error {

	planIdDescription := `Required. plan id`

	var planIdFlagName string
	if cmdPrefix == "" {
		planIdFlagName = "planId"
	} else {
		planIdFlagName = fmt.Sprintf("%v.planId", cmdPrefix)
	}

	var planIdFlagDefault string

	_ = cmd.PersistentFlags().String(planIdFlagName, planIdFlagDefault, planIdDescription)

	return nil
}

func retrieveOperationPlansExecutablesUnlinkPlanSchedulePlanIDFlag(m *plans_executables.UnlinkPlanScheduleParams, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false
	if cmd.Flags().Changed("planId") {

		var planIdFlagName string
		if cmdPrefix == "" {
			planIdFlagName = "planId"
		} else {
			planIdFlagName = fmt.Sprintf("%v.planId", cmdPrefix)
		}

		planIdFlagValue, err := cmd.Flags().GetString(planIdFlagName)
		if err != nil {
			return err, false
		}
		m.PlanID = planIdFlagValue

	}

	return nil, retAdded
}

// parseOperationPlansExecutablesUnlinkPlanScheduleResult parses request result and return the string content
func parseOperationPlansExecutablesUnlinkPlanScheduleResult(resp0 *plans_executables.UnlinkPlanScheduleNoContent, respErr error) (string, error) {
	if respErr != nil {

		var iResp0 interface{} = respErr
		resp0, ok := iResp0.(*plans_executables.UnlinkPlanScheduleNoContent)
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
		resp1, ok := iResp1.(*plans_executables.UnlinkPlanScheduleBadRequest)
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
		resp2, ok := iResp2.(*plans_executables.UnlinkPlanScheduleUnauthorized)
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
		resp3, ok := iResp3.(*plans_executables.UnlinkPlanScheduleForbidden)
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
		resp4, ok := iResp4.(*plans_executables.UnlinkPlanScheduleNotFound)
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
		resp5, ok := iResp5.(*plans_executables.UnlinkPlanScheduleInternalServerError)
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
