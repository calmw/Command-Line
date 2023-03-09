package cmd

import (
	"flag"
	"fmt"
	"io"
)

type httpConfig struct {
	url  string
	verb string
}

func HandleHttp(w io.Writer, args []string) error {
	var v string
	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&v, "verb", "GET", "HTTP method") // verb变量名
	fs.Usage = func() {
		var usageString = `
http: A HTTP client.
http: <option> server`
		_, _ = fmt.Fprintf(w, usageString)
		_, _ = fmt.Fprintln(w)
		_, _ = fmt.Fprintln(w)
		_, _ = fmt.Fprintln(w, "Options:")
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return ErrNoServerSpecified
	}

	c := httpConfig{verb: v}
	c.url = fs.Arg(0)
	_, _ = fmt.Fprintln(w, "Executing http command")
	return nil
}
