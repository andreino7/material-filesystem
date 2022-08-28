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
	Use:   "rm [FILE]",
	Short: "Remove files or directories",
	Long: `rm removes each specified file. By default, it does not remove directories.
Supports absolute and relative paths

Examples:
rm /file1
rm -r /dir1/dir2
rm file1
`,
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
	removeChildren = rmCmd.Flags().BoolP("recursive", "r", false, "remove directories and their contents recursively")

}
