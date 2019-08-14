package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/k0kubun/go-ansi"
	"github.com/mattn/go-isatty"
	"os"
)

var inputFileName string = "-"

func main() {
	flag.StringVar(&inputFileName, "f", "-", "file to process or '-' to process stdin")
	flag.Parse()

	in := os.Stdin
	out := ansi.NewAnsiStdout()
	if inputFileName != "-" {
		f, err := os.Open(inputFileName)
		exitOnUsageError(err)
		in = f
		defer in.Close()
	}

	if isatty.IsTerminal(in.Fd()) || isatty.IsCygwinTerminal(in.Fd()) {
		exitOnUsageError(fmt.Errorf("the command is intended to work with pipes or files"))
	}

	reader := bufio.NewReader(in)
	zap := zapper{reader, out}
	err := zap.pipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(2)
	}
}

func exitOnUsageError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, `Usage: 

  zapper [options]

Options:
  -f string
        file to process or '-' to process stdin (default "-")

Examples:
  cat yourlogfile.log | zapper
  zapper < yourlogfile.log
  zapper -f yourlogfile.log

`)
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
