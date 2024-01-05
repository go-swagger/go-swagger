// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"

	"github.com/go-swagger/go-swagger/fixtures/bugs/2969/codegen/models"
	"github.com/spf13/cobra"
)

// Schema cli for DownloadLogsStatus

// register flags to command
func registerModelDownloadLogsStatusFlags(depth int, cmdPrefix string, cmd *cobra.Command) error {

	if err := registerDownloadLogsStatusPropPresignedURL(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	if err := registerDownloadLogsStatusPropStatus(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	return nil
}

func registerDownloadLogsStatusPropPresignedURL(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	presignedUrlDescription := ``

	var presignedUrlFlagName string
	if cmdPrefix == "" {
		presignedUrlFlagName = "presignedURL"
	} else {
		presignedUrlFlagName = fmt.Sprintf("%v.presignedURL", cmdPrefix)
	}

	var presignedUrlFlagDefault string

	_ = cmd.PersistentFlags().String(presignedUrlFlagName, presignedUrlFlagDefault, presignedUrlDescription)

	return nil
}

func registerDownloadLogsStatusPropStatus(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	statusDescription := `Enum: ["READY","IN_PROGRESS"]. `

	var statusFlagName string
	if cmdPrefix == "" {
		statusFlagName = "status"
	} else {
		statusFlagName = fmt.Sprintf("%v.status", cmdPrefix)
	}

	var statusFlagDefault string

	_ = cmd.PersistentFlags().String(statusFlagName, statusFlagDefault, statusDescription)

	if err := cmd.RegisterFlagCompletionFunc(statusFlagName,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			var res []string
			if err := json.Unmarshal([]byte(`["READY","IN_PROGRESS"]`), &res); err != nil {
				panic(err)
			}
			return res, cobra.ShellCompDirectiveDefault
		}); err != nil {
		return err
	}

	return nil
}

// retrieve flags from commands, and set value in model. Return true if any flag is passed by user to fill model field.
func retrieveModelDownloadLogsStatusFlags(depth int, m *models.DownloadLogsStatus, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false

	err, presignedUrlAdded := retrieveDownloadLogsStatusPropPresignedURLFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || presignedUrlAdded

	err, statusAdded := retrieveDownloadLogsStatusPropStatusFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || statusAdded

	return nil, retAdded
}

func retrieveDownloadLogsStatusPropPresignedURLFlags(depth int, m *models.DownloadLogsStatus, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	presignedUrlFlagName := fmt.Sprintf("%v.presignedURL", cmdPrefix)
	if cmd.Flags().Changed(presignedUrlFlagName) {

		var presignedUrlFlagName string
		if cmdPrefix == "" {
			presignedUrlFlagName = "presignedURL"
		} else {
			presignedUrlFlagName = fmt.Sprintf("%v.presignedURL", cmdPrefix)
		}

		presignedUrlFlagValue, err := cmd.Flags().GetString(presignedUrlFlagName)
		if err != nil {
			return err, false
		}
		m.PresignedURL = presignedUrlFlagValue

		retAdded = true
	}

	return nil, retAdded
}

func retrieveDownloadLogsStatusPropStatusFlags(depth int, m *models.DownloadLogsStatus, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	statusFlagName := fmt.Sprintf("%v.status", cmdPrefix)
	if cmd.Flags().Changed(statusFlagName) {

		var statusFlagName string
		if cmdPrefix == "" {
			statusFlagName = "status"
		} else {
			statusFlagName = fmt.Sprintf("%v.status", cmdPrefix)
		}

		statusFlagValue, err := cmd.Flags().GetString(statusFlagName)
		if err != nil {
			return err, false
		}
		m.Status = statusFlagValue

		retAdded = true
	}

	return nil, retAdded
}
