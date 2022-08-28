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
	Use:   "cd",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	fsclient.Session.SetWorkingDirName(resp.GetChangeWorkingDirectory().GetName())
}

func init() {
	rootCmd.AddCommand(cdCmd)
}
