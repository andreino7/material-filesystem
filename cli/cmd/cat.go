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

// readCmd represents the read command
var catCmd = &cobra.Command{
	Use:   "cat [FILE]",
	Short: "Print files",
	Long: `Print the contents of a file to the standard output.
Supports absolute and relative paths.

Examples:
cat file1
cat /dir1/file1
cat dir1/file1
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("invalid argument")
		}

		req := &fsservice.Request{
			Request: &fsservice.Request_ReadAll{
				ReadAll: &fsservice.ReadAllRequest{
					Path: args[0],
				},
			},
		}
		fsclient.Session.DoRequest(req, fsclient.Session.ReadAll, printReadAll)
		return nil
	},
}

func printReadAll(resp *fsservice.Response) {
	fmt.Println(string(resp.GetReadAll().GetContent()))
}

func init() {
	rootCmd.AddCommand(catCmd)
}
