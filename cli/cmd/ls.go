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
	Use:   "ls",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
