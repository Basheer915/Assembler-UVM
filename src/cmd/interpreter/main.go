package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"uvmproject/src/interpreter"
)

type options struct {
	inputPath string
	dumpPath  string
	start     int
	end       int
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "interpreter:", err)
		os.Exit(1)
	}
}

func run(args []string, output, errorOutput io.Writer) error {
	configuration, err := parseOptions(args, errorOutput)
	if err != nil {
		return err
	}
	machineCode, err := os.ReadFile(configuration.inputPath)
	if err != nil {
		return fmt.Errorf("read machine-code file: %w", err)
	}
	memory, err := interpreter.Run(machineCode)
	if err != nil {
		return err
	}
	if err := interpreter.WriteXMLDump(
		configuration.dumpPath,
		memory,
		configuration.start,
		configuration.end,
	); err != nil {
		return err
	}
	_, err = fmt.Fprintf(output, "XML-дамп сохранен: %s\n", configuration.dumpPath)
	return err
}

func parseOptions(args []string, errorOutput io.Writer) (options, error) {
	flags := flag.NewFlagSet("interpreter", flag.ContinueOnError)
	flags.SetOutput(errorOutput)
	inputPath := flags.String("input", "", "path to binary machine-code file")
	dumpPath := flags.String("dump", "", "path to XML memory dump")
	start := flags.Int("start", -1, "first memory address to include")
	end := flags.Int("end", -1, "last memory address to include")
	if err := flags.Parse(args); err != nil {
		return options{}, err
	}
	if flags.NArg() != 0 {
		return options{}, fmt.Errorf("unexpected positional arguments: %v", flags.Args())
	}
	if *inputPath == "" || *dumpPath == "" {
		return options{}, fmt.Errorf("both -input and -dump are required")
	}
	if *start < 0 || *end < 0 || *start > *end {
		return options{}, fmt.Errorf("provide a valid -start and -end address range")
	}
	return options{
		inputPath: *inputPath,
		dumpPath:  *dumpPath,
		start:     *start,
		end:       *end,
	}, nil
}
