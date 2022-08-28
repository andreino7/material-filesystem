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

// touchCmd represents the touch command
var createCmd = &cobra.Command{
	Use:   "create [FILE]",
	Short: "Create a file",
	Long: `Create the FILE, if it does not already exist.
Supports absolute and relative paths.

Examples:
create file1
create /dir1/file1
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("invalid argument")
		}
		req := &fsservice.Request{
			Request: &fsservice.Request_CreateRegularFile{
				CreateRegularFile: &fsservice.CreateRegularFileRequest{
					Path: args[0],
				},
			},
		}
		fsclient.Session.DoRequest(req, fsclient.Session.CreateRegularFile, noop)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
