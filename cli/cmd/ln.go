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

var isSymbolicLink *bool

// lnCmd represents the ln command
var lnCmd = &cobra.Command{
	Use:   "ln [LINK_TARGET] [LINK_PATH]",
	Short: "Make links between files",
	Long: `Create a link to LINK_TARGET at the given LINK_PATH.
Create hard links by default, symbolic links with --symbolic.
Hard links to directories are not supported.
Supports absolute and relative paths.

Examples:
ln file1 file1-link
ln /dir/file1 /dir2/file1-link

ln -s file1 file1-link
ln -s /dir/file1 /dir-link
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("invalid argument")
		}

		req := &fsservice.Request{}
		if *isSymbolicLink {
			req.Request = &fsservice.Request_SymLink{
				SymLink: &fsservice.CreateSymbolicLinkRequest{
					SrcPath:  args[0],
					DestPath: args[1],
				},
			}
			fsclient.Session.DoRequest(req, fsclient.Session.CreateSymbolicLink, noop)
		} else {
			req.Request = &fsservice.Request_HardLink{
				HardLink: &fsservice.CreateHardLinkRequest{
					SrcPath:  args[0],
					DestPath: args[1],
				},
			}
			fsclient.Session.DoRequest(req, fsclient.Session.CreateHardLink, noop)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(lnCmd)
	lnCmd.PostRun = lnPostRun
	lnPostRun(nil, nil)
}

func lnPostRun(cmd *cobra.Command, args []string) {
	lnCmd.ResetFlags()
	isSymbolicLink = lnCmd.Flags().BoolP("symbolic", "s", false, "make symbolic links instead of hard links")
}
