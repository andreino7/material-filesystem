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

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls [DIRECTORY]",
	Short: "List directory contents",
	Long: `List the FILEs names (in the current directory by
default). Supports absolute and relative paths.

Examples:
ls
ls /dir/dir1/dir2
ls dir1/
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return fmt.Errorf("invalid argument")
		}

		path := ""
		if len(args) == 1 {
			path = args[0]
		}

		req := &fsservice.Request{
			Request: &fsservice.Request_List{
				List: &fsservice.ListFilesRequest{
					Path: path,
				},
			},
		}
		fsclient.Session.DoRequest(req, fsclient.Session.ListFiles, printFileNames)
		return nil
	},
}

func printFileNames(resp *fsservice.Response) {
	for _, name := range resp.GetList().GetNames() {
		fmt.Println(name)
	}
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
