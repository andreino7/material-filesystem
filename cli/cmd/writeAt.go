/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"material/filesystem/cli/fsclient"
	"material/filesystem/pb/proto/fsservice"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// writeAtCmd represents the writeAt command
var writeAtCmd = &cobra.Command{
	Use:   "writeAt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return fmt.Errorf("invalid argument")
		}

		start, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid argument")
		}

		text := strings.Join(args[2:], " ")

		req := &fsservice.Request{
			Request: &fsservice.Request_WriteAt{
				WriteAt: &fsservice.WriteAtRequest{
					FileDescriptor: args[0],
					Pos:            int32(start),
					Content:        []byte(text),
				},
			},
		}
		fsclient.Session.DoRequest(req, fsclient.Session.WriteAt, noop)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(writeAtCmd)
}
