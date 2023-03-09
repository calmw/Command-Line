package main

import (
	"Command-Line/cmd"
	"errors"
	"fmt"
	"io"
	"os"
)

var errInvalidSubCommand = errors.New("invalid sub-command specified")

func main() {
	err := handleCommand(os.Stdout, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
}

func printUsage(w io.Writer) {
	_, _ = fmt.Fprintf(w, "Usage: mync [http|grpc] -h\n")
	_ = cmd.HandleHttp(w, []string{"-h"})
	_ = cmd.HandleGrpc(w, []string{"-h"})
}

func handleCommand(w io.Writer, args []string) error {
	var err error
	if len(args) < 1 {
		err = errInvalidSubCommand
	} else {
		switch args[0] {
		case "http":
			err = cmd.HandleHttp(w, args[1:])
		case "grpc":
			err = cmd.HandleGrpc(w, args[1:])
		case "-h":
			printUsage(w)
		case "--help":
			printUsage(w)
		default:
			err = errInvalidSubCommand
		}
	}

	if errors.Is(err, cmd.ErrNoServerSpecified) || errors.Is(err, errInvalidSubCommand) {
		_, _ = fmt.Fprintln(w, err)
		printUsage(w)
	}
	return err
}
