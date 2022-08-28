/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"bufio"
	"fmt"
	"material/filesystem/cli/cmd"
	"material/filesystem/cli/fsclient"
	"os"
	"strings"
)

func main() {
	conn, err := fsclient.Initialize()
	if err != nil {
		fmt.Printf("critical error: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Welcome to material filesystem client!")
	fmt.Println("Type help for help :)")

	for {
		fmt.Printf("material-filesystem %s $ ", fsclient.Session.WorkingDirName())
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()
		args := strings.Split(text, " ")
		cmd.Execute(args)
	}
}
