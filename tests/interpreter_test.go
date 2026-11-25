package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"uvmproject/src/assembler"
	"uvmproject/src/interpreter"
)

func TestRunCopiesArrayAndWritesXMLDump(t *testing.T) {
	source, err := os.Open("../examples/copy_array.csv")
	if err != nil {
		t.Fatalf("open copy program: %v", err)
	}
	defer source.Close()

	instructions, err := assembler.Translate(source)
	if err != nil {
		t.Fatalf("Translate() error = %v", err)
	}
	machineCode, err := assembler.Encode(instructions)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	memory, err := interpreter.Run(machineCode)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	for index, want := range []int32{7, 12, 24} {
		if got := memory[200+index]; got != want {
			t.Errorf("memory[%d] = %d, want %d", 200+index, got, want)
		}
	}
	assertXMLDump(t, memory)
}

func assertXMLDump(t *testing.T, memory []int32) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "memory.xml")
	if err := interpreter.WriteXMLDump(path, memory, 200, 202); err != nil {
		t.Fatalf("WriteXMLDump() error = %v", err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read memory dump: %v", err)
	}
	want := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
		"<memory>\n" +
		"  <cell address=\"200\">7</cell>\n" +
		"  <cell address=\"201\">12</cell>\n" +
		"  <cell address=\"202\">24</cell>\n" +
		"</memory>\n"
	if strings.TrimSpace(string(content)) != strings.TrimSpace(want) {
		t.Errorf("XML dump =\n%s\nwant:\n%s", content, want)
	}
}

func TestRunRejectsStackUnderflow(t *testing.T) {
	machineCode, err := assembler.Encode([]assembler.Instruction{
		{Command: "READ", A: 5},
	})
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	if _, err := interpreter.Run(machineCode); err == nil {
		t.Fatal("Run() error = nil, want stack-underflow error")
	}
}

func TestRunRejectsOutOfRangeReadAddress(t *testing.T) {
	machineCode, err := assembler.Encode([]assembler.Instruction{
		{Command: "LOAD", A: 0, B: interpreter.DataMemorySize},
		{Command: "READ", A: 5},
	})
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	if _, err := interpreter.Run(machineCode); err == nil {
		t.Fatal("Run() error = nil, want out-of-range address error")
	}
}

func TestWriteXMLDumpRejectsInvalidRange(t *testing.T) {
	path := filepath.Join(t.TempDir(), "memory.xml")
	if err := interpreter.WriteXMLDump(path, []int32{1, 2}, 0, 2); err == nil {
		t.Fatal("WriteXMLDump() error = nil, want invalid range error")
	}
}
