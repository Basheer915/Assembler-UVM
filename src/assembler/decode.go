package assembler

import "fmt"

func Decode(machineCode []byte) ([]Instruction, error) {
	if len(machineCode)%InstructionSize != 0 {
		return nil, fmt.Errorf(
			"machine code length %d is not a multiple of %d",
			len(machineCode),
			InstructionSize,
		)
	}

	instructions := make([]Instruction, 0, len(machineCode)/InstructionSize)
	for offset := 0; offset < len(machineCode); offset += InstructionSize {
		word := uint32(machineCode[offset]) |
			uint32(machineCode[offset+1])<<8 |
			uint32(machineCode[offset+2])<<16
		instruction, err := decodeInstruction(word)
		if err != nil {
			return nil, fmt.Errorf("instruction %d: %w", offset/InstructionSize, err)
		}
		instructions = append(instructions, instruction)
	}
	return instructions, nil
}

func decodeInstruction(word uint32) (Instruction, error) {
	opcode := uint8(word & ((1 << operandShift) - 1))
	operand := word >> operandShift

	var instruction Instruction
	switch opcode {
	case loadOpcode:
		instruction = Instruction{Command: "LOAD", A: opcode, B: operand}
	case storeOpcode:
		instruction = Instruction{Command: "STORE", A: opcode, B: operand}
	case readOpcode:
		instruction = Instruction{Command: "READ", A: opcode, B: operand}
	case absOpcode:
		instruction = Instruction{Command: "ABS", A: opcode, B: operand}
	default:
		return Instruction{}, fmt.Errorf("unknown opcode %d", opcode)
	}
	if _, err := encodeInstruction(instruction); err != nil {
		return Instruction{}, err
	}
	return instruction, nil
}
