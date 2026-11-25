package tests

import (
	"bytes"
	"strings"
	"testing"

	"uvmproject/src/assembler"
)

func TestTranslateSpecificationCommands(t *testing.T) {
	source := strings.NewReader(
		"command,operand\n" +
			"LOAD,319\n" +
			"READ,\n" +
			"STORE,675\n" +
			"ABS,\n",
	)

	got, err := assembler.Translate(source)
	if err != nil {
		t.Fatalf("Translate() error = %v", err)
	}

	want := []assembler.Instruction{
		{Command: "LOAD", A: 0, B: 319},
		{Command: "READ", A: 5, B: 0},
		{Command: "STORE", A: 2, B: 675},
		{Command: "ABS", A: 4, B: 0},
	}
	if len(got) != len(want) {
		t.Fatalf("Translate() returned %d instructions, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("instruction %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestTranslateRejectsInvalidCommandsAndOperands(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{name: "unknown command", source: "command,operand\nJUMP,\n"},
		{name: "load out of range", source: "command,operand\nLOAD,131072\n"},
		{name: "store out of range", source: "command,operand\nSTORE,1024\n"},
		{name: "missing operand", source: "command,operand\nLOAD,\n"},
		{name: "unexpected operand", source: "command,operand\nABS,1\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := assembler.Translate(strings.NewReader(test.source)); err == nil {
				t.Fatal("Translate() error = nil, want an error")
			}
		})
	}
}

func TestEncodeSpecificationCommands(t *testing.T) {
	instructions := []assembler.Instruction{
		{Command: "LOAD", A: 0, B: 319},
		{Command: "READ", A: 5, B: 0},
		{Command: "STORE", A: 2, B: 675},
		{Command: "ABS", A: 4, B: 0},
	}
	want := []byte{
		0xF8, 0x09, 0x00,
		0x05, 0x00, 0x00,
		0x1A, 0x15, 0x00,
		0x04, 0x00, 0x00,
	}

	got, err := assembler.Encode(instructions)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("Encode() = % X, want % X", got, want)
	}
}

func TestEncodeRejectsInvalidInstructionFields(t *testing.T) {
	instructions := []assembler.Instruction{
		{Command: "LOAD", A: 1, B: 319},
	}
	if _, err := assembler.Encode(instructions); err == nil {
		t.Fatal("Encode() error = nil, want invalid field error")
	}
}

func TestDecodeSpecificationCommands(t *testing.T) {
	machineCode := []byte{
		0xF8, 0x09, 0x00,
		0x05, 0x00, 0x00,
		0x1A, 0x15, 0x00,
		0x04, 0x00, 0x00,
	}
	want := []assembler.Instruction{
		{Command: "LOAD", A: 0, B: 319},
		{Command: "READ", A: 5, B: 0},
		{Command: "STORE", A: 2, B: 675},
		{Command: "ABS", A: 4, B: 0},
	}

	got, err := assembler.Decode(machineCode)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("Decode() returned %d instructions, want %d", len(got), len(want))
	}
	for index := range want {
		if got[index] != want[index] {
			t.Errorf("instruction %d = %+v, want %+v", index, got[index], want[index])
		}
	}
}

func TestDecodeRejectsInvalidMachineCode(t *testing.T) {
	tests := []struct {
		name string
		code []byte
	}{
		{name: "incomplete instruction", code: []byte{0x01}},
		{name: "unknown opcode", code: []byte{0x06, 0x00, 0x00}},
		{name: "operand exceeds field", code: []byte{0x00, 0x00, 0x10}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := assembler.Decode(test.code); err == nil {
				t.Fatal("Decode() error = nil, want invalid machine-code error")
			}
		})
	}
}
