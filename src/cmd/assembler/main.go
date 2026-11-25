package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"uvmproject/src/assembler"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "assembler:", err)
		os.Exit(1)
	}
}

func run(args []string, output, errorOutput io.Writer) error {
	flags := flag.NewFlagSet("assembler", flag.ContinueOnError)
	flags.SetOutput(errorOutput)
	inputPath := flags.String("input", "", "path to CSV source")
	outputPath := flags.String("output", "", "path to JSON intermediate output")
	testMode := flags.Bool("test", false, "print intermediate instruction fields")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *inputPath == "" || *outputPath == "" {
		return fmt.Errorf("both -input and -output are required")
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments: %v", flags.Args())
	}

	source, err := os.Open(*inputPath)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}
	defer source.Close()

	instructions, err := assembler.Translate(source)
	if err != nil {
		return err
	}
	encoded, err := json.MarshalIndent(instructions, "", "  ")
	if err != nil {
		return fmt.Errorf("encode intermediate representation: %w", err)
	}
	encoded = append(encoded, '\n')
	if err := os.MkdirAll(filepath.Dir(*outputPath), 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	if err := os.WriteFile(*outputPath, encoded, 0644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	if *testMode {
		for index, instruction := range instructions {
			if _, err := fmt.Fprintf(
				output,
				"%d: %s A=%d B=%d\n",
				index,
				instruction.Command,
				instruction.A,
				instruction.B,
			); err != nil {
				return fmt.Errorf("print intermediate representation: %w", err)
			}
		}
	}
	return nil
}
