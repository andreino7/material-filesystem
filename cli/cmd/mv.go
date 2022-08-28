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

// mvCmd represents the mv command
var mvCmd = &cobra.Command{
	Use:   "mv [SOURCE] [DEST]",
	Short: "Move (rename) files",
	Long: `Rename SOURCE to DEST, or move SOURCE(s) to DIRECTORY
and resolves any name conflict by merging directories and renaming files. 
Creates all parent directories of DEST.
Supports relative and absolute paths.

Examples:
mv dir1 dir2
mv /dir1/file1 dir1/file2`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("invalid argument")
		}

		req := &fsservice.Request{
			Request: &fsservice.Request_Move{
				Move: &fsservice.MoveRequest{
					SrcPath:  args[0],
					DestPath: args[1],
				},
			},
		}
		fsclient.Session.DoRequest(req, fsclient.Session.Move, noop)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(mvCmd)
}
