package interpreter

import (
	"fmt"

	"uvmproject/src/assembler"
)

const (
	dataAddressBits = 10
	DataMemorySize  = 1 << dataAddressBits
)

func Run(machineCode []byte) ([]int32, error) {
	program, err := assembler.Decode(machineCode)
	if err != nil {
		return nil, fmt.Errorf("decode program: %w", err)
	}

	memory := make([]int32, DataMemorySize)
	stack := make([]int32, 0, len(program))
	for programCounter, instruction := range program {
		if err := executeInstruction(instruction, memory, &stack); err != nil {
			return nil, fmt.Errorf(
				"instruction %d (%s): %w",
				programCounter,
				instruction.Command,
				err,
			)
		}
	}
	return memory, nil
}

func executeInstruction(
	instruction assembler.Instruction,
	memory []int32,
	stack *[]int32,
) error {
	switch instruction.Command {
	case "LOAD":
		*stack = append(*stack, int32(instruction.B))
	case "READ":
		addressValue, err := pop(stack)
		if err != nil {
			return err
		}
		address, err := memoryAddress(addressValue, len(memory))
		if err != nil {
			return err
		}
		*stack = append(*stack, memory[address])
	case "STORE":
		value, err := pop(stack)
		if err != nil {
			return err
		}
		address, err := memoryAddress(int32(instruction.B), len(memory))
		if err != nil {
			return err
		}
		memory[address] = value
	case "ABS":
		return fmt.Errorf("ABS execution is not implemented yet")
	default:
		return fmt.Errorf("unsupported command %q", instruction.Command)
	}
	return nil
}

func pop(stack *[]int32) (int32, error) {
	if len(*stack) == 0 {
		return 0, fmt.Errorf("stack underflow")
	}
	last := len(*stack) - 1
	value := (*stack)[last]
	*stack = (*stack)[:last]
	return value, nil
}

func memoryAddress(value int32, memorySize int) (int, error) {
	if value < 0 || int64(value) >= int64(memorySize) {
		return 0, fmt.Errorf("memory address %d is out of range", value)
	}
	return int(value), nil
}
