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

// cdCmd represents the cd command
var cdCmd = &cobra.Command{
	Use:   "cd [directory]",
	Short: "Change the working directory",
	Long: `The cd utility change the working directory of the current
session. Supports absolute or relative paths

Examples:
	cd dir1
	cd ../dir1/dir2
	cd /dir1/dir2/dir3`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("invalid arguments")
		}
		req := &fsservice.Request{
			Request: &fsservice.Request_ChangeWorkingDirectory{
				ChangeWorkingDirectory: &fsservice.ChangeWorkingDirectoryRequest{
					Path: args[0],
				},
			},
		}
		fsclient.Session.DoRequest(req, fsclient.Session.ChangeWorkingDirectory, updateWokingDirectory)
		return nil
	},
}

func updateWokingDirectory(resp *fsservice.Response) {
	fsclient.Session.SetWorkingDirPath(resp.GetChangeWorkingDirectory().GetPath())
}

func init() {
	rootCmd.AddCommand(cdCmd)
}
