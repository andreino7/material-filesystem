/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"material/filesystem/cli/fsclient"
	"material/filesystem/pb/proto/fsservice"
	"strconv"

	"github.com/spf13/cobra"
)

// readAtCmd represents the readAt command
var readAtCmd = &cobra.Command{
	Use:   "readAt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 3 {
			return fmt.Errorf("invalid argument")
		}

		start, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid argument")
		}

		end, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("invalid argument")
		}

		req := &fsservice.Request{
			Request: &fsservice.Request_ReadAt{
				ReadAt: &fsservice.ReadAtRequest{
					FileDescriptor: args[0],
					StartPos:       int32(start),
					EndPos:         int32(end),
				},
			},
		}
		fsclient.Session.DoRequest(req, fsclient.Session.ReadAt, printReadAt)
		return nil
	},
}

func printReadAt(resp *fsservice.Response) {
	fmt.Println(string(resp.GetReadAt().GetContent()))
}

func init() {
	rootCmd.AddCommand(readAtCmd)
}
