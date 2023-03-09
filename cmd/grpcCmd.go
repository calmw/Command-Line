package cmd

import (
	"flag"
	"fmt"
	"io"
)

type grpcConfig struct {
	server string
	method string
	body   string
}

func HandleGrpc(w io.Writer, args []string) error {
	c := grpcConfig{}
	fs := flag.NewFlagSet("grpc", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&c.method, "method", "", "Method to call")
	fs.StringVar(&c.body, "body", "", "Body of request")
	fs.Usage = func() {
		var usageString = `
grpc: A grpc client.
grpc: <option> server`
		_, _ = fmt.Fprintf(w, usageString)
		_, _ = fmt.Fprintln(w) // 这里会输出空行
		_, _ = fmt.Fprintln(w, "Options:")
		fs.PrintDefaults()
	}
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	if fs.NArg() != 2 {
		return ErrNoServerSpecified
	}
	c.server = fs.Arg(0)
	_, _ = fmt.Fprintln(w, "Executing grpc command")
	return nil
}
