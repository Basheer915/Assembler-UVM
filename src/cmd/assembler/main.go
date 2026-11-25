package main

import (
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
	outputPath := flags.String("output", "", "path to binary machine-code output")
	testMode := flags.Bool("test", false, "print instruction fields and encoded bytes")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *inputPath == "" || *outputPath == "" {
		return fmt.Errorf("both -input and -output are required")
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments: %v", flags.Args())
	}

	instructions, machineCode, err := assembleFile(*inputPath)
	if err != nil {
		return err
	}
	if err := writeOutput(*outputPath, machineCode); err != nil {
		return err
	}
	if *testMode {
		if err := printTestOutput(output, instructions, machineCode); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(output, "Размер двоичного файла: %d байт\n", len(machineCode)); err != nil {
		return fmt.Errorf("print output size: %w", err)
	}
	return nil
}

func assembleFile(inputPath string) ([]assembler.Instruction, []byte, error) {
	source, err := os.Open(inputPath)
	if err != nil {
		return nil, nil, fmt.Errorf("open input: %w", err)
	}
	defer source.Close()

	instructions, err := assembler.Translate(source)
	if err != nil {
		return nil, nil, err
	}
	machineCode, err := assembler.Encode(instructions)
	if err != nil {
		return nil, nil, fmt.Errorf("encode machine code: %w", err)
	}
	return instructions, machineCode, nil
}

func writeOutput(outputPath string, machineCode []byte) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	if err := os.WriteFile(outputPath, machineCode, 0644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}

func printTestOutput(output io.Writer, instructions []assembler.Instruction, machineCode []byte) error {
	for index, instruction := range instructions {
		offset := index * assembler.InstructionSize
		if _, err := fmt.Fprintf(
			output,
			"%d: %s A=%d B=%d -> 0x%02X, 0x%02X, 0x%02X\n",
			index,
			instruction.Command,
			instruction.A,
			instruction.B,
			machineCode[offset],
			machineCode[offset+1],
			machineCode[offset+2],
		); err != nil {
			return fmt.Errorf("print encoded instruction: %w", err)
		}
	}
	return nil
}
