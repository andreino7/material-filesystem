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
	Use:   "find",
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
