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

var createParentDirs *bool

// mkdirCmd represents the mkdir command
var mkdirCmd = &cobra.Command{
	Use:   "mkdir",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("invalid argument")
		}
		req := &fsservice.Request{
			Request: &fsservice.Request_Mkdir{
				Mkdir: &fsservice.MkdirRequest{
					Path:      args[0],
					Recursive: createParentDirs,
				},
			},
		}
		fsclient.Session.DoRequest(req, fsclient.Session.Mkdir, noop)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(mkdirCmd)
	mkdirCmd.PostRun = postRun
	postRun(nil, nil)
}

func postRun(cmd *cobra.Command, args []string) {
	mkdirCmd.ResetFlags()
	createParentDirs = mkdirCmd.Flags().BoolP("parents", "p", false, "make parent directories as needed")
}
