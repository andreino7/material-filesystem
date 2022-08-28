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

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find [PATH] [NAME]",
	Short: "Search for files in a directory hierarchy",
	Long: `Find searches the directory tree rooted at the given
path and prints list of files/directories matching the given name. 
Supports absolute and relative paths.

Examples:
find / some_name
find dir1/dir2 some_name`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("invalid argument")
		}
		req := &fsservice.Request{
			Request: &fsservice.Request_Find{
				Find: &fsservice.FindFilesRequest{
					Path: args[0],
					Name: args[1],
				},
			},
		}
		fsclient.Session.DoRequest(req, fsclient.Session.FindFiles, printFind)
		return nil
	},
}

func printFind(resp *fsservice.Response) {
	for _, path := range resp.GetFind().GetPaths() {
		fmt.Println(path)
	}
}

func init() {
	rootCmd.AddCommand(findCmd)
}
