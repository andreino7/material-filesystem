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

var removeChildren *bool

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
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
			Request: &fsservice.Request_Remove{
				Remove: &fsservice.RemoveRequest{
					Path:      args[0],
					Recursive: removeChildren,
				},
			},
		}
		fsclient.Session.DoRequest(req, fsclient.Session.Remove, noop)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.PostRun = rmPostRun
	rmPostRun(nil, nil)
}

func rmPostRun(cmd *cobra.Command, args []string) {
	rmCmd.ResetFlags()
	removeChildren = mkdirCmd.Flags().BoolP("recursive", "r", false, "make parent directories as needed")

}
