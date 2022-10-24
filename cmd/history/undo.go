// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package history

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"storj.io/storj-up/cmd"
	"storj.io/storj-up/pkg/common"
)

func undoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "undo",
		Short: "revert to a previous version of the generated docker compose file",
		RunE: func(cmd *cobra.Command, args []string) error {
			newTemplateBytes, err := common.Store.RestoreLatestVersion()
			if err != nil {
				return err
			}
			if newTemplateBytes == nil {
				return fmt.Errorf("no previous version of the compose file found")
			}
			newTemplate, err := common.LoadComposeFromBytes(newTemplateBytes)
			if err != nil {
				return err
			}
			pwd, _ := os.Getwd()
			err = common.WriteComposeFileNoHistory(pwd, newTemplate)
			if err != nil {
				return err
			}
			return nil
		},
	}
}

func init() {
	cmd.RootCmd.AddCommand(undoCmd())
}
