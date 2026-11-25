package assembler

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	loadMax      = 1<<17 - 1
	storeMax     = 1<<10 - 1
	loadOpcode   = 0
	storeOpcode  = 2
	readOpcode   = 5
	absOpcode    = 4
	operandShift = 3
)

type Instruction struct {
	Command string `json:"command"`
	A       uint8  `json:"A"`
	B       uint32 `json:"B"`
}

func Translate(source io.Reader) ([]Instruction, error) {
	reader := csv.NewReader(source)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read CSV header: %w", err)
	}
	if len(header) != 2 ||
		!strings.EqualFold(strings.TrimSpace(header[0]), "command") ||
		!strings.EqualFold(strings.TrimSpace(header[1]), "operand") {
		return nil, fmt.Errorf("CSV header must be command,operand")
	}

	var instructions []Instruction
	for row := 2; ; row++ {
		record, err := reader.Read()
		if err == io.EOF {
			return instructions, nil
		}
		if err != nil {
			return nil, fmt.Errorf("read CSV row %d: %w", row, err)
		}
		instruction, err := translateRecord(record)
		if err != nil {
			return nil, fmt.Errorf("CSV row %d: %w", row, err)
		}
		instructions = append(instructions, instruction)
	}
}

func translateRecord(record []string) (Instruction, error) {
	if len(record) < 1 || len(record) > 2 {
		return Instruction{}, fmt.Errorf("expected a command and optional operand")
	}
	command := strings.ToUpper(strings.TrimSpace(record[0]))
	operand := ""
	if len(record) == 2 {
		operand = strings.TrimSpace(record[1])
	}

	switch command {
	case "LOAD":
		value, err := parseOperand(operand, loadMax)
		if err != nil {
			return Instruction{}, fmt.Errorf("LOAD: %w", err)
		}
		return Instruction{Command: command, A: loadOpcode, B: value}, nil
	case "STORE":
		value, err := parseOperand(operand, storeMax)
		if err != nil {
			return Instruction{}, fmt.Errorf("STORE: %w", err)
		}
		return Instruction{Command: command, A: storeOpcode, B: value}, nil
	case "READ":
		if operand != "" {
			return Instruction{}, fmt.Errorf("READ does not take an operand")
		}
		return Instruction{Command: command, A: readOpcode}, nil
	case "ABS":
		if operand != "" {
			return Instruction{}, fmt.Errorf("ABS does not take an operand")
		}
		return Instruction{Command: command, A: absOpcode}, nil
	default:
		return Instruction{}, fmt.Errorf("unknown command %q", command)
	}
}

func parseOperand(operand string, max uint64) (uint32, error) {
	if operand == "" {
		return 0, fmt.Errorf("requires an unsigned integer operand")
	}
	value, err := strconv.ParseUint(operand, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid unsigned integer operand %q", operand)
	}
	if value > max {
		return 0, fmt.Errorf("operand %d exceeds maximum %d", value, max)
	}
	return uint32(value), nil
}
