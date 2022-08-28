/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"material/filesystem/cli/fsclient"
	"material/filesystem/pb/proto/fsservice"

	"github.com/spf13/cobra"
)

// cpCmd represents the cp command
var cpCmd = &cobra.Command{
	Use:   "cp",
	Short: "Copy files and directories",
	Long: `Copy SOURCE to DEST and resolves any name conflict
by merging directories and renaming files`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("invalid argument")
		}

		req := &fsservice.Request{
			Request: &fsservice.Request_Copy{
				Copy: &fsservice.CopyRequest{
					SrcPath:  args[0],
					DestPath: args[1],
				},
			},
		}
		fsclient.Session.DoRequest(req, fsclient.Session.Copy, noop)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cpCmd)
}
