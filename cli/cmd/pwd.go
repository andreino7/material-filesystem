/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"material/filesystem/cli/fsclient"

	"github.com/spf13/cobra"
)

// pwdCmd represents the pwd command
var pwdCmd = &cobra.Command{
	Use:   "pwd",
	Short: "Print the path of the current/working directory",
	Long: `Print the full path of the current working directory.

Examples:
pwd
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fsclient.Session.WorkingDirPath())
	},
}

func init() {
	rootCmd.AddCommand(pwdCmd)
}
