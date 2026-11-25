package assembler

import "fmt"

const InstructionSize = 3

func Encode(instructions []Instruction) ([]byte, error) {
	machineCode := make([]byte, 0, len(instructions)*InstructionSize)
	for index, instruction := range instructions {
		word, err := encodeInstruction(instruction)
		if err != nil {
			return nil, fmt.Errorf("instruction %d: %w", index, err)
		}
		machineCode = append(
			machineCode,
			byte(word),
			byte(word>>8),
			byte(word>>16),
		)
	}
	return machineCode, nil
}

func encodeInstruction(instruction Instruction) (uint32, error) {
	switch instruction.Command {
	case "LOAD":
		if instruction.A != loadOpcode || instruction.B > loadMax {
			return 0, fmt.Errorf("invalid LOAD fields A=%d B=%d", instruction.A, instruction.B)
		}
	case "STORE":
		if instruction.A != storeOpcode || instruction.B > storeMax {
			return 0, fmt.Errorf("invalid STORE fields A=%d B=%d", instruction.A, instruction.B)
		}
	case "READ":
		if instruction.A != readOpcode || instruction.B != 0 {
			return 0, fmt.Errorf("invalid READ fields A=%d B=%d", instruction.A, instruction.B)
		}
	case "ABS":
		if instruction.A != absOpcode || instruction.B != 0 {
			return 0, fmt.Errorf("invalid ABS fields A=%d B=%d", instruction.A, instruction.B)
		}
	default:
		return 0, fmt.Errorf("unknown command %q", instruction.Command)
	}
	return uint32(instruction.A) | instruction.B<<operandShift, nil
}
