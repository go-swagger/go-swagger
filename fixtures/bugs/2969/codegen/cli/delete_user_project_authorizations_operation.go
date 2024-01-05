// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"

	"github.com/go-swagger/go-swagger/fixtures/bugs/2969/codegen/client/projects_authorization"

	"github.com/go-openapi/swag"
	"github.com/spf13/cobra"
)

// makeOperationProjectsAuthorizationDeleteUserProjectAuthorizationsCmd returns a command to handle operation deleteUserProjectAuthorizations
func makeOperationProjectsAuthorizationDeleteUserProjectAuthorizationsCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "deleteUserProjectAuthorizations",
		Short: `Delete user project authorization`,
		RunE:  runOperationProjectsAuthorizationDeleteUserProjectAuthorizations,
	}

	if err := registerOperationProjectsAuthorizationDeleteUserProjectAuthorizationsParamFlags(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}

// runOperationProjectsAuthorizationDeleteUserProjectAuthorizations uses cmd flags to call endpoint api
func runOperationProjectsAuthorizationDeleteUserProjectAuthorizations(cmd *cobra.Command, args []string) error {
	appCli, err := makeClient(cmd, args)
	if err != nil {
		return err
	}
	// retrieve flag values from cmd and fill params
	params := projects_authorization.NewDeleteUserProjectAuthorizationsParams()
	if err, _ = retrieveOperationProjectsAuthorizationDeleteUserProjectAuthorizationsProjectIDFlag(params, "", cmd); err != nil {
		return err
	}
	if err, _ = retrieveOperationProjectsAuthorizationDeleteUserProjectAuthorizationsUserIDFlag(params, "", cmd); err != nil {
		return err
	}
	if dryRun {
		logDebugf("dry-run flag specified. Skip sending request.")
		return nil
	}
	// make request and then print result
	msgStr, err := parseOperationProjectsAuthorizationDeleteUserProjectAuthorizationsResult(appCli.ProjectsAuthorization.DeleteUserProjectAuthorizations(params, nil))
	if err != nil {
		return err
	}

	if !debug {
		fmt.Println(msgStr)
	}

	return nil
}

// registerOperationProjectsAuthorizationDeleteUserProjectAuthorizationsParamFlags registers all flags needed to fill params
func registerOperationProjectsAuthorizationDeleteUserProjectAuthorizationsParamFlags(cmd *cobra.Command) error {
	if err := registerOperationProjectsAuthorizationDeleteUserProjectAuthorizationsProjectIDParamFlags("", cmd); err != nil {
		return err
	}
	if err := registerOperationProjectsAuthorizationDeleteUserProjectAuthorizationsUserIDParamFlags("", cmd); err != nil {
		return err
	}
	return nil
}

func registerOperationProjectsAuthorizationDeleteUserProjectAuthorizationsProjectIDParamFlags(cmdPrefix string, cmd *cobra.Command) error {

	projectIdDescription := `Required. project ID`

	var projectIdFlagName string
	if cmdPrefix == "" {
		projectIdFlagName = "projectId"
	} else {
		projectIdFlagName = fmt.Sprintf("%v.projectId", cmdPrefix)
	}

	var projectIdFlagDefault string

	_ = cmd.PersistentFlags().String(projectIdFlagName, projectIdFlagDefault, projectIdDescription)

	return nil
}

func registerOperationProjectsAuthorizationDeleteUserProjectAuthorizationsUserIDParamFlags(cmdPrefix string, cmd *cobra.Command) error {

	userIdDescription := `Required. project ID`

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

func retrieveOperationProjectsAuthorizationDeleteUserProjectAuthorizationsProjectIDFlag(m *projects_authorization.DeleteUserProjectAuthorizationsParams, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false
	if cmd.Flags().Changed("projectId") {

		var projectIdFlagName string
		if cmdPrefix == "" {
			projectIdFlagName = "projectId"
		} else {
			projectIdFlagName = fmt.Sprintf("%v.projectId", cmdPrefix)
		}

		projectIdFlagValue, err := cmd.Flags().GetString(projectIdFlagName)
		if err != nil {
			return err, false
		}
		m.ProjectID = projectIdFlagValue

	}

	return nil, retAdded
}

func retrieveOperationProjectsAuthorizationDeleteUserProjectAuthorizationsUserIDFlag(m *projects_authorization.DeleteUserProjectAuthorizationsParams, cmdPrefix string, cmd *cobra.Command) (error, bool) {
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

// parseOperationProjectsAuthorizationDeleteUserProjectAuthorizationsResult parses request result and return the string content
func parseOperationProjectsAuthorizationDeleteUserProjectAuthorizationsResult(resp0 *projects_authorization.DeleteUserProjectAuthorizationsNoContent, respErr error) (string, error) {
	if respErr != nil {

		// Non schema case: warning deleteUserProjectAuthorizationsNoContent is not supported

		var iResp1 interface{} = respErr
		resp1, ok := iResp1.(*projects_authorization.DeleteUserProjectAuthorizationsBadRequest)
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
		resp2, ok := iResp2.(*projects_authorization.DeleteUserProjectAuthorizationsUnauthorized)
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
		resp3, ok := iResp3.(*projects_authorization.DeleteUserProjectAuthorizationsForbidden)
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
		resp4, ok := iResp4.(*projects_authorization.DeleteUserProjectAuthorizationsNotFound)
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
		resp5, ok := iResp5.(*projects_authorization.DeleteUserProjectAuthorizationsInternalServerError)
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

	// warning: non schema response deleteUserProjectAuthorizationsNoContent is not supported by go-swagger cli yet.

	return "", nil
}
