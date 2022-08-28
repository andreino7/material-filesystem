/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"material/filesystem/cli/cmd"
	"material/filesystem/cli/fsclient"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const defaultPort = "2222"

func main() {

	port := os.Getenv("FS_DAEMON_PORT")
	if port == "" {
		port = defaultPort
	}
	// Start grpc
	if err := fsclient.Initialize(port); err != nil {
		log.Fatalf("critical error: %v", err)
		os.Exit(1)
	}

	// Graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		<-sigs
		fsclient.Close()
		done <- true
	}()

	fmt.Println("Welcome to material filesystem client!")
	fmt.Println("Type help for help :)")

	go startCli()

	// Wait for shutdown signal
	<-done
}

func startCli() {
	for {
		fmt.Printf("material-filesystem %s $ ", fsclient.Session.WorkingDirName())
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()
		args := strings.Split(text, " ")
		cmd.Execute(args)
	}
}
