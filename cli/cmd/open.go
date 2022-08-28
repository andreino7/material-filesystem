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

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open [PATH]",
	Short: "Open a file",
	Long: `The open call opens the file specified by pathname and
prints the file descriptor to stdout.
Supports absolute and relative paths.

Examples:
open /file1
open file1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("invalid argument")
		}
		req := &fsservice.Request{
			Request: &fsservice.Request_Open{
				Open: &fsservice.OpenRequest{
					Path: args[0],
				},
			},
		}
		fsclient.Session.DoRequest(req, fsclient.Session.Open, printFd)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}

func printFd(resp *fsservice.Response) {
	fmt.Println(resp.GetOpen().GetFileDescriptor())
}
