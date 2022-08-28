/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"material/filesystem/cli/fsclient"
	"material/filesystem/pb/proto/fsservice"
	"strings"

	"github.com/spf13/cobra"
)

// appendCmd represents the append command
var appendCmd = &cobra.Command{
	Use:   "append [FILE] [CONTENT]",
	Short: "Append all content to a file",
	Long: `Append the contents passed as argument to the end of the file.
Supports absolute and relative paths. If the file does not exist,
a new one is automatically created.

Examples:
append file1 some text to append
append /dir1/file1 some text to append
appen dir1/file1 some text to append
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("invalid argument")
		}

		text := strings.Join(args[1:], " ")

		req := &fsservice.Request{
			Request: &fsservice.Request_AppendAll{
				AppendAll: &fsservice.AppendAllRequest{
					Path:    args[0],
					Content: []byte(text),
				},
			},
		}
		fsclient.Session.DoRequest(req, fsclient.Session.AppendAll, noop)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(appendCmd)
}
