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
	Use:   "readAt [FILE_DESCRIPTOR] [START] [END]",
	Short: "Print files",
	Long: `Print the contents of a file to the standard output.

Examples:
readAt fd 10 100
`,
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
