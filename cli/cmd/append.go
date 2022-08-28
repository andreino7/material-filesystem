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

var text *string

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
		fmt.Println("qui")
		fmt.Println(args)
		if len(args) != 1 {
			return fmt.Errorf("invalid argument")
		}

		req := &fsservice.Request{
			Request: &fsservice.Request_AppendAll{
				AppendAll: &fsservice.AppendAllRequest{
					Path:    args[0],
					Content: []byte(*text),
				},
			},
		}
		fsclient.Session.DoRequest(req, fsclient.Session.AppendAll, noop)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(appendCmd)
	appendCmd.PostRun = appendPostRun
	appendPostRun(nil, nil)
}

func appendPostRun(cmd *cobra.Command, args []string) {
	appendCmd.ResetFlags()
	text = appendCmd.Flags().String("string", "", "usage")
}
