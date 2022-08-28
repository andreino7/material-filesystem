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
	Use:   "append",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
