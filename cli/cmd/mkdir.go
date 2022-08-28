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
	Use:   "mkdir [DIRECTORY]",
	Short: "Make directories",
	Long: `Create the DIRECTORY(ies), if they do not already exist.
Supports absolute and relative paths.

Examples:
mkdir dir
mkdir /dir
mkdir -p /dir1/dir2/dir3`,
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
