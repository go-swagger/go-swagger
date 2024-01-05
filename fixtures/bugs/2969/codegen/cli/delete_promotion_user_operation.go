// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"

	"github.com/go-swagger/go-swagger/fixtures/bugs/2969/codegen/client/promotions_executables_authorization"

	"github.com/go-openapi/swag"
	"github.com/spf13/cobra"
)

// makeOperationPromotionsExecutablesAuthorizationDeletePromotionUserCmd returns a command to handle operation deletePromotionUser
func makeOperationPromotionsExecutablesAuthorizationDeletePromotionUserCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "deletePromotionUser",
		Short: `Unassign user from Promotion`,
		RunE:  runOperationPromotionsExecutablesAuthorizationDeletePromotionUser,
	}

	if err := registerOperationPromotionsExecutablesAuthorizationDeletePromotionUserParamFlags(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}

// runOperationPromotionsExecutablesAuthorizationDeletePromotionUser uses cmd flags to call endpoint api
func runOperationPromotionsExecutablesAuthorizationDeletePromotionUser(cmd *cobra.Command, args []string) error {
	appCli, err := makeClient(cmd, args)
	if err != nil {
		return err
	}
	// retrieve flag values from cmd and fill params
	params := promotions_executables_authorization.NewDeletePromotionUserParams()
	if err, _ = retrieveOperationPromotionsExecutablesAuthorizationDeletePromotionUserPromotionIDFlag(params, "", cmd); err != nil {
		return err
	}
	if err, _ = retrieveOperationPromotionsExecutablesAuthorizationDeletePromotionUserUserIDFlag(params, "", cmd); err != nil {
		return err
	}
	if dryRun {
		logDebugf("dry-run flag specified. Skip sending request.")
		return nil
	}
	// make request and then print result
	msgStr, err := parseOperationPromotionsExecutablesAuthorizationDeletePromotionUserResult(appCli.PromotionsExecutablesAuthorization.DeletePromotionUser(params, nil))
	if err != nil {
		return err
	}

	if !debug {
		fmt.Println(msgStr)
	}

	return nil
}

// registerOperationPromotionsExecutablesAuthorizationDeletePromotionUserParamFlags registers all flags needed to fill params
func registerOperationPromotionsExecutablesAuthorizationDeletePromotionUserParamFlags(cmd *cobra.Command) error {
	if err := registerOperationPromotionsExecutablesAuthorizationDeletePromotionUserPromotionIDParamFlags("", cmd); err != nil {
		return err
	}
	if err := registerOperationPromotionsExecutablesAuthorizationDeletePromotionUserUserIDParamFlags("", cmd); err != nil {
		return err
	}
	return nil
}

func registerOperationPromotionsExecutablesAuthorizationDeletePromotionUserPromotionIDParamFlags(cmdPrefix string, cmd *cobra.Command) error {

	promotionIdDescription := `Required. Promotion ID`

	var promotionIdFlagName string
	if cmdPrefix == "" {
		promotionIdFlagName = "promotionId"
	} else {
		promotionIdFlagName = fmt.Sprintf("%v.promotionId", cmdPrefix)
	}

	var promotionIdFlagDefault string

	_ = cmd.PersistentFlags().String(promotionIdFlagName, promotionIdFlagDefault, promotionIdDescription)

	return nil
}

func registerOperationPromotionsExecutablesAuthorizationDeletePromotionUserUserIDParamFlags(cmdPrefix string, cmd *cobra.Command) error {

	userIdDescription := `Required. User ID`

	var userIdFlagName string
	if cmdPrefix == "" {
		userIdFlagName = "userId"
	} else {
		userIdFlagName = fmt.Sprintf("%v.userId", cmdPrefix)
	}

	var userIdFlagDefault string

	_ = cmd.PersistentFlags().String(userIdFlagName, userIdFlagDefault, userIdDescription)

	return nil
}

func retrieveOperationPromotionsExecutablesAuthorizationDeletePromotionUserPromotionIDFlag(m *promotions_executables_authorization.DeletePromotionUserParams, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false
	if cmd.Flags().Changed("promotionId") {

		var promotionIdFlagName string
		if cmdPrefix == "" {
			promotionIdFlagName = "promotionId"
		} else {
			promotionIdFlagName = fmt.Sprintf("%v.promotionId", cmdPrefix)
		}

		promotionIdFlagValue, err := cmd.Flags().GetString(promotionIdFlagName)
		if err != nil {
			return err, false
		}
		m.PromotionID = promotionIdFlagValue

	}

	return nil, retAdded
}

func retrieveOperationPromotionsExecutablesAuthorizationDeletePromotionUserUserIDFlag(m *promotions_executables_authorization.DeletePromotionUserParams, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false
	if cmd.Flags().Changed("userId") {

		var userIdFlagName string
		if cmdPrefix == "" {
			userIdFlagName = "userId"
		} else {
			userIdFlagName = fmt.Sprintf("%v.userId", cmdPrefix)
		}

		userIdFlagValue, err := cmd.Flags().GetString(userIdFlagName)
		if err != nil {
			return err, false
		}
		m.UserID = userIdFlagValue

	}

	return nil, retAdded
}

// parseOperationPromotionsExecutablesAuthorizationDeletePromotionUserResult parses request result and return the string content
func parseOperationPromotionsExecutablesAuthorizationDeletePromotionUserResult(resp0 *promotions_executables_authorization.DeletePromotionUserNoContent, respErr error) (string, error) {
	if respErr != nil {

		// Non schema case: warning deletePromotionUserNoContent is not supported

		var iResp1 interface{} = respErr
		resp1, ok := iResp1.(*promotions_executables_authorization.DeletePromotionUserBadRequest)
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
		resp2, ok := iResp2.(*promotions_executables_authorization.DeletePromotionUserUnauthorized)
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
		resp3, ok := iResp3.(*promotions_executables_authorization.DeletePromotionUserForbidden)
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
		resp4, ok := iResp4.(*promotions_executables_authorization.DeletePromotionUserNotFound)
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
		resp5, ok := iResp5.(*promotions_executables_authorization.DeletePromotionUserInternalServerError)
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

	// warning: non schema response deletePromotionUserNoContent is not supported by go-swagger cli yet.

	return "", nil
}
