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

// closeCmd represents the close command
var closeCmd = &cobra.Command{
	Use:   "close [FILE_DESCRIPTOR]",
	Short: "Close a file",
	Long: `Close a file using the file descriptor.

Examples:
	close file-descriptor`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("invalid argument")
		}
		req := &fsservice.Request{
			Request: &fsservice.Request_Close{
				Close: &fsservice.CloseRequest{
					FileDescriptor: args[0],
				},
			},
		}
		fsclient.Session.DoRequest(req, fsclient.Session.Close, noop)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(closeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// closeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// closeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
