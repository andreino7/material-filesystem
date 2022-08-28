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

var isSymbolicLink *bool

// lnCmd represents the ln command
var lnCmd = &cobra.Command{
	Use:   "ln",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("invalid argument")
		}

		req := &fsservice.Request{}
		if *isSymbolicLink {
			req.Request = &fsservice.Request_SymLink{
				SymLink: &fsservice.CreateSymbolicLinkRequest{
					SrcPath:  args[0],
					DestPath: args[1],
				},
			}
			fsclient.Session.DoRequest(req, fsclient.Session.CreateSymbolicLink, noop)
		} else {
			req.Request = &fsservice.Request_HardLink{
				HardLink: &fsservice.CreateHardLinkRequest{
					SrcPath:  args[0],
					DestPath: args[1],
				},
			}
			fsclient.Session.DoRequest(req, fsclient.Session.CreateHardLink, noop)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(lnCmd)
	lnCmd.PostRun = lnPostRun
	lnPostRun(nil, nil)
}

func lnPostRun(cmd *cobra.Command, args []string) {
	lnCmd.ResetFlags()
	isSymbolicLink = lnCmd.Flags().BoolP("symbolic", "s", false, "make parent directories as needed")
}
