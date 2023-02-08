package main

type Instruction struct {
	mnemonic       string
	addressingMode int
	length         int
	cycles         int
	execute        func(*CPU, uint16)
}

var instructions = map[uint8]Instruction{
	0x00: {mnemonic: "BRK", addressingMode: 0, length: 1, cycles: 7, execute: func(cpu *CPU, operand uint16) {
		cpu.brk()
	}},
	0x01: {mnemonic: "ORA", addressingMode: 11, length: 2, cycles: 6, execute: func(cpu *CPU, operand uint16) {}},
	0x05: {mnemonic: "ORA", addressingMode: 3, length: 2, cycles: 3, execute: func(cpu *CPU, operand uint16) {}},
	0x06: {mnemonic: "ASL", addressingMode: 3, length: 2, cycles: 5, execute: func(cpu *CPU, operand uint16) {
		cpu.asl(operand)
	}},
	0x08: {mnemonic: "PHP", addressingMode: 0, length: 1, cycles: 3, execute: func(cpu *CPU, operand uint16) {}},
	0x09: {mnemonic: "ORA", addressingMode: 2, length: 2, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0x0A: {mnemonic: "ASL", addressingMode: 1, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {
		cpu.aslAcc()
	}},
	0x0D: {mnemonic: "ORA", addressingMode: 7, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0x0E: {mnemonic: "ASL", addressingMode: 7, length: 3, cycles: 6, execute: func(cpu *CPU, operand uint16) {
		cpu.asl(operand)
	}},
	0x10: {mnemonic: "BPL", addressingMode: 6, length: 2, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0x11: {mnemonic: "ORA", addressingMode: 9, length: 2, cycles: 5, execute: func(cpu *CPU, operand uint16) {}},
	0x15: {mnemonic: "ORA", addressingMode: 4, length: 2, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0x16: {mnemonic: "ASL", addressingMode: 4, length: 2, cycles: 6, execute: func(cpu *CPU, operand uint16) {
		cpu.asl(operand)
	}},
	0x18: {mnemonic: "CLC", addressingMode: 0, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {
		cpu.clc()
	}},
	0x19: {mnemonic: "ORA", addressingMode: 8, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0x1D: {mnemonic: "ORA", addressingMode: 8, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0x1E: {mnemonic: "ASL", addressingMode: 8, length: 3, cycles: 7, execute: func(cpu *CPU, operand uint16) {
		cpu.asl(operand)
	}},
	0x20: {mnemonic: "JSR", addressingMode: 7, length: 3, cycles: 6, execute: func(cpu *CPU, operand uint16) {}},
	0x21: {mnemonic: "AND", addressingMode: 11, length: 2, cycles: 6, execute: func(cpu *CPU, operand uint16) {}},
	0x24: {mnemonic: "BIT", addressingMode: 3, length: 2, cycles: 3, execute: func(cpu *CPU, operand uint16) {}},
	0x25: {mnemonic: "AND", addressingMode: 3, length: 2, cycles: 3, execute: func(cpu *CPU, operand uint16) {}},
	0x26: {mnemonic: "ROL", addressingMode: 3, length: 2, cycles: 5, execute: func(cpu *CPU, operand uint16) {}},
	0x28: {mnemonic: "PLP", addressingMode: 0, length: 1, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0x29: {mnemonic: "AND", addressingMode: 2, length: 2, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0x2A: {mnemonic: "ROL", addressingMode: 1, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0x2C: {mnemonic: "BIT", addressingMode: 7, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0x2D: {mnemonic: "AND", addressingMode: 7, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0x2E: {mnemonic: "ROL", addressingMode: 7, length: 3, cycles: 6, execute: func(cpu *CPU, operand uint16) {}},
	0x30: {mnemonic: "BMI", addressingMode: 6, length: 2, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0x31: {mnemonic: "AND", addressingMode: 9, length: 2, cycles: 5, execute: func(cpu *CPU, operand uint16) {}},
	0x35: {mnemonic: "AND", addressingMode: 4, length: 2, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0x36: {mnemonic: "ROL", addressingMode: 4, length: 2, cycles: 6, execute: func(cpu *CPU, operand uint16) {}},
	0x38: {mnemonic: "SEC", addressingMode: 0, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {
		cpu.sec()
	}},
	0x39: {mnemonic: "AND", addressingMode: 8, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0x3D: {mnemonic: "AND", addressingMode: 8, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0x3E: {mnemonic: "ROL", addressingMode: 8, length: 3, cycles: 7, execute: func(cpu *CPU, operand uint16) {}},
	0x40: {mnemonic: "RTI", addressingMode: 0, length: 1, cycles: 6, execute: func(cpu *CPU, operand uint16) {}},
	0x41: {mnemonic: "EOR", addressingMode: 11, length: 2, cycles: 6, execute: func(cpu *CPU, operand uint16) {}},
	0x45: {mnemonic: "EOR", addressingMode: 3, length: 2, cycles: 3, execute: func(cpu *CPU, operand uint16) {}},
	0x46: {mnemonic: "LSR", addressingMode: 3, length: 2, cycles: 5, execute: func(cpu *CPU, operand uint16) {}},
	0x48: {mnemonic: "PHA", addressingMode: 0, length: 1, cycles: 3, execute: func(cpu *CPU, operand uint16) {}},
	0x49: {mnemonic: "EOR", addressingMode: 2, length: 2, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0x4A: {mnemonic: "LSR", addressingMode: 1, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0x4C: {mnemonic: "JMP", addressingMode: 7, length: 3, cycles: 3, execute: func(cpu *CPU, operand uint16) {
		cpu.jmp(operand)
	}},
	0x4D: {mnemonic: "EOR", addressingMode: 7, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0x4E: {mnemonic: "LSR", addressingMode: 7, length: 3, cycles: 6, execute: func(cpu *CPU, operand uint16) {}},
	0x50: {mnemonic: "BVC", addressingMode: 6, length: 2, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0x51: {mnemonic: "EOR", addressingMode: 9, length: 2, cycles: 5, execute: func(cpu *CPU, operand uint16) {}},
	0x55: {mnemonic: "EOR", addressingMode: 4, length: 2, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0x56: {mnemonic: "LSR", addressingMode: 4, length: 2, cycles: 6, execute: func(cpu *CPU, operand uint16) {}},
	0x58: {mnemonic: "CLI", addressingMode: 0, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {
		cpu.cli()
	}},
	0x59: {mnemonic: "EOR", addressingMode: 8, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0x5D: {mnemonic: "EOR", addressingMode: 8, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0x5E: {mnemonic: "LSR", addressingMode: 8, length: 3, cycles: 7, execute: func(cpu *CPU, operand uint16) {}},
	0x60: {mnemonic: "RTS", addressingMode: 0, length: 1, cycles: 6, execute: func(cpu *CPU, operand uint16) {}},
	0x61: {mnemonic: "ADC", addressingMode: 11, length: 2, cycles: 6, execute: func(cpu *CPU, operand uint16) {
		cpu.adc(operand)
	}},
	0x65: {mnemonic: "ADC", addressingMode: 3, length: 2, cycles: 3, execute: func(cpu *CPU, operand uint16) {
		cpu.adc(operand)
	}},
	0x66: {mnemonic: "ROR", addressingMode: 3, length: 2, cycles: 5, execute: func(cpu *CPU, operand uint16) {
		cpu.ror(operand)
	}},
	0x68: {mnemonic: "PLA", addressingMode: 0, length: 1, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0x69: {mnemonic: "ADC", addressingMode: 2, length: 2, cycles: 2, execute: func(cpu *CPU, operand uint16) {
		cpu.adc(operand)
	}},
	0x6A: {mnemonic: "ROR", addressingMode: 1, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {
		cpu.ror(operand)
	}},
	0x6C: {mnemonic: "JMP", addressingMode: 10, length: 3, cycles: 5, execute: func(cpu *CPU, operand uint16) {
		cpu.jmp(operand)
	}},
	0x6D: {mnemonic: "ADC", addressingMode: 7, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {
		cpu.adc(operand)
	}},
	0x6E: {mnemonic: "ROR", addressingMode: 7, length: 3, cycles: 6, execute: func(cpu *CPU, operand uint16) {
		cpu.ror(operand)
	}},
	0x70: {mnemonic: "BVS", addressingMode: 6, length: 2, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0x71: {mnemonic: "ADC", addressingMode: 9, length: 2, cycles: 5, execute: func(cpu *CPU, operand uint16) {
		cpu.adc(operand)
	}},
	0x75: {mnemonic: "ADC", addressingMode: 4, length: 2, cycles: 4, execute: func(cpu *CPU, operand uint16) {
		cpu.adc(operand)
	}},
	0x76: {mnemonic: "ROR", addressingMode: 4, length: 2, cycles: 6, execute: func(cpu *CPU, operand uint16) {
		cpu.ror(operand)
	}},
	0x78: {mnemonic: "SEI", addressingMode: 0, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {
		cpu.sei()
	}},
	0x79: {mnemonic: "ADC", addressingMode: 8, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {
		cpu.adc(operand)
	}},
	0x7D: {mnemonic: "ADC", addressingMode: 8, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {
		cpu.adc(operand)
	}},
	0x7E: {mnemonic: "ROR", addressingMode: 8, length: 3, cycles: 7, execute: func(cpu *CPU, operand uint16) {
		cpu.ror(operand)
	}},
	0x81: {mnemonic: "STA", addressingMode: 11, length: 2, cycles: 6, execute: func(cpu *CPU, operand uint16) {
		cpu.sta(operand)
	}},
	0x84: {mnemonic: "STY", addressingMode: 3, length: 2, cycles: 3, execute: func(cpu *CPU, operand uint16) {}},
	0x85: {mnemonic: "STA", addressingMode: 3, length: 2, cycles: 3, execute: func(cpu *CPU, operand uint16) {
		cpu.sta(operand)
	}},
	0x86: {mnemonic: "STX", addressingMode: 3, length: 2, cycles: 3, execute: func(cpu *CPU, operand uint16) {
		cpu.stx(operand)
	}},
	0x88: {mnemonic: "DEY", addressingMode: 0, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0x8A: {mnemonic: "TXA", addressingMode: 0, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0x8C: {mnemonic: "STY", addressingMode: 7, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0x8D: {mnemonic: "STA", addressingMode: 7, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {
		cpu.sta(operand)
	}},
	0x8E: {mnemonic: "STX", addressingMode: 7, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {
		cpu.stx(operand)
	}},
	0x90: {mnemonic: "BCC", addressingMode: 6, length: 2, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0x91: {mnemonic: "STA", addressingMode: 9, length: 2, cycles: 6, execute: func(cpu *CPU, operand uint16) {
		cpu.sta(operand)
	}},
	0x94: {mnemonic: "STY", addressingMode: 4, length: 2, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0x95: {mnemonic: "STA", addressingMode: 4, length: 2, cycles: 4, execute: func(cpu *CPU, operand uint16) {
		cpu.sta(operand)
	}},
	0x96: {mnemonic: "STX", addressingMode: 5, length: 2, cycles: 4, execute: func(cpu *CPU, operand uint16) {
		cpu.stx(operand)
	}},
	0x98: {mnemonic: "TYA", addressingMode: 0, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0x99: {mnemonic: "STA", addressingMode: 8, length: 3, cycles: 5, execute: func(cpu *CPU, operand uint16) {
		cpu.sta(operand)
	}},
	0x9A: {mnemonic: "TXS", addressingMode: 0, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0x9D: {mnemonic: "STA", addressingMode: 8, length: 3, cycles: 5, execute: func(cpu *CPU, operand uint16) {
		cpu.sta(operand)
	}},
	0xA0: {mnemonic: "LDY", addressingMode: 2, length: 2, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0xA1: {mnemonic: "LDA", addressingMode: 11, length: 2, cycles: 6, execute: func(cpu *CPU, operand uint16) {
		cpu.lda(operand)
	}},
	0xA2: {mnemonic: "LDX", addressingMode: 2, length: 2, cycles: 2, execute: func(cpu *CPU, operand uint16) {
		cpu.ldx(operand)
	}},
	0xA4: {mnemonic: "LDY", addressingMode: 3, length: 2, cycles: 3, execute: func(cpu *CPU, operand uint16) {}},
	0xA5: {mnemonic: "LDA", addressingMode: 3, length: 2, cycles: 3, execute: func(cpu *CPU, operand uint16) {
		cpu.lda(operand)
	}},
	0xA6: {mnemonic: "LDX", addressingMode: 3, length: 2, cycles: 3, execute: func(cpu *CPU, operand uint16) {
		cpu.ldx(operand)
	}},
	0xA8: {mnemonic: "TAY", addressingMode: 0, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0xA9: {mnemonic: "LDA", addressingMode: 2, length: 2, cycles: 2, execute: func(cpu *CPU, operand uint16) {
		cpu.lda(operand)
	}},
	0xAA: {mnemonic: "TAX", addressingMode: 0, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0xAC: {mnemonic: "LDY", addressingMode: 7, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0xAD: {mnemonic: "LDA", addressingMode: 7, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {
		cpu.lda(operand)
	}},
	0xAE: {mnemonic: "LDX", addressingMode: 7, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {
		cpu.ldx(operand)
	}},
	0xB0: {mnemonic: "BCS", addressingMode: 6, length: 2, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0xB1: {mnemonic: "LDA", addressingMode: 9, length: 2, cycles: 5, execute: func(cpu *CPU, operand uint16) {
		cpu.lda(operand)
	}},
	0xB4: {mnemonic: "LDY", addressingMode: 4, length: 2, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0xB5: {mnemonic: "LDA", addressingMode: 4, length: 2, cycles: 4, execute: func(cpu *CPU, operand uint16) {
		cpu.lda(operand)
	}},
	0xB6: {mnemonic: "LDX", addressingMode: 5, length: 2, cycles: 4, execute: func(cpu *CPU, operand uint16) {
		cpu.ldx(operand)
	}},
	0xB8: {mnemonic: "CLV", addressingMode: 0, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {
		cpu.clv()
	}},
	0xB9: {mnemonic: "LDA", addressingMode: 8, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {
		cpu.lda(operand)
	}},
	0xBA: {mnemonic: "TSX", addressingMode: 0, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0xBC: {mnemonic: "LDY", addressingMode: 8, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0xBD: {mnemonic: "LDA", addressingMode: 8, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {
		cpu.lda(operand)
	}},
	0xBE: {mnemonic: "LDX", addressingMode: 8, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {
		cpu.ldx(operand)
	}},
	0xC0: {mnemonic: "CPY", addressingMode: 2, length: 2, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0xC1: {mnemonic: "CMP", addressingMode: 11, length: 2, cycles: 6, execute: func(cpu *CPU, operand uint16) {}},
	0xC4: {mnemonic: "CPY", addressingMode: 3, length: 2, cycles: 3, execute: func(cpu *CPU, operand uint16) {}},
	0xC5: {mnemonic: "CMP", addressingMode: 3, length: 2, cycles: 3, execute: func(cpu *CPU, operand uint16) {}},
	0xC6: {mnemonic: "DEC", addressingMode: 3, length: 2, cycles: 5, execute: func(cpu *CPU, operand uint16) {}},
	0xC8: {mnemonic: "INY", addressingMode: 0, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0xC9: {mnemonic: "CMP", addressingMode: 2, length: 2, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0xCA: {mnemonic: "DEX", addressingMode: 0, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0xCC: {mnemonic: "CPY", addressingMode: 7, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0xCD: {mnemonic: "CMP", addressingMode: 7, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0xCE: {mnemonic: "DEC", addressingMode: 7, length: 3, cycles: 6, execute: func(cpu *CPU, operand uint16) {}},
	0xD0: {mnemonic: "BNE", addressingMode: 6, length: 2, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0xD1: {mnemonic: "CMP", addressingMode: 9, length: 2, cycles: 5, execute: func(cpu *CPU, operand uint16) {}},
	0xD5: {mnemonic: "CMP", addressingMode: 4, length: 2, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0xD6: {mnemonic: "DEC", addressingMode: 4, length: 2, cycles: 6, execute: func(cpu *CPU, operand uint16) {}},
	0xD8: {mnemonic: "CLD", addressingMode: 0, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {
		cpu.cld()
	}},
	0xD9: {mnemonic: "CMP", addressingMode: 8, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0xDD: {mnemonic: "CMP", addressingMode: 8, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0xDE: {mnemonic: "DEC", addressingMode: 8, length: 3, cycles: 7, execute: func(cpu *CPU, operand uint16) {}},
	0xE0: {mnemonic: "CPX", addressingMode: 2, length: 2, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0xE1: {mnemonic: "SBC", addressingMode: 11, length: 2, cycles: 6, execute: func(cpu *CPU, operand uint16) {}},
	0xE4: {mnemonic: "CPX", addressingMode: 3, length: 2, cycles: 3, execute: func(cpu *CPU, operand uint16) {}},
	0xE5: {mnemonic: "SBC", addressingMode: 3, length: 2, cycles: 3, execute: func(cpu *CPU, operand uint16) {}},
	0xE6: {mnemonic: "INC", addressingMode: 3, length: 2, cycles: 5, execute: func(cpu *CPU, operand uint16) {
		cpu.inc(operand)
	}},
	0xE8: {mnemonic: "INX", addressingMode: 0, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0xE9: {mnemonic: "SBC", addressingMode: 2, length: 2, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0xEA: {mnemonic: "NOP", addressingMode: 0, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {
		cpu.nop()
	}},
	0xEC: {mnemonic: "CPX", addressingMode: 7, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0xED: {mnemonic: "SBC", addressingMode: 7, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0xEE: {mnemonic: "INC", addressingMode: 7, length: 3, cycles: 6, execute: func(cpu *CPU, operand uint16) {
		cpu.inc(operand)
	}},
	0xF0: {mnemonic: "BEQ", addressingMode: 6, length: 2, cycles: 2, execute: func(cpu *CPU, operand uint16) {}},
	0xF1: {mnemonic: "SBC", addressingMode: 9, length: 2, cycles: 5, execute: func(cpu *CPU, operand uint16) {}},
	0xF5: {mnemonic: "SBC", addressingMode: 4, length: 2, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0xF6: {mnemonic: "INC", addressingMode: 4, length: 2, cycles: 6, execute: func(cpu *CPU, operand uint16) {
		cpu.inc(operand)
	}},
	0xF8: {mnemonic: "SED", addressingMode: 0, length: 1, cycles: 2, execute: func(cpu *CPU, operand uint16) {
		cpu.sed()
	}},
	0xF9: {mnemonic: "SBC", addressingMode: 8, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0xFD: {mnemonic: "SBC", addressingMode: 8, length: 3, cycles: 4, execute: func(cpu *CPU, operand uint16) {}},
	0xFE: {mnemonic: "INC", addressingMode: 8, length: 3, cycles: 7, execute: func(cpu *CPU, operand uint16) {
		cpu.inc(operand)
	}},
}

func boolToInt(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func (cpu *CPU) adc(address uint16) {
	// Fetch the data from the address
	data := cpu.MMU.readByte(address)
	// Add the data to the accumulator
	result := uint16(cpu.A) + uint16(data)
	// Add one to the result if the carry flag is set
	if cpu.getFlag(Carry) {
		result++
	}

	// Set the zero flag to opposite of (result & 0xFF)
	cpu.setFlag(Zero, result&0xFF == 0)

	// Check if decimal mode is enabled
	if cpu.getFlag(Decimal) {
		// Do BCD addition
		// C# version: if (((registers.A & 0xF) + (value & 0xF) + (GetFlag(StatusFlags.Carry) ? 1 : 0)) > 9)
		if ((cpu.A & 0xF) + (data & 0xF) + boolToInt(cpu.getFlag(Carry))) > 9 {
			result += 6
		}

		// Set the negative flag if the result is negative
		if result&0x80 == 0x80 {
			cpu.setFlag(Negative, true)
		}
		// Set the overflow flag if needed
		cpu.setFlag(Overflow, ^uint16(cpu.A)^uint16(data)&0x80 != 0 && ^uint16(cpu.A)^result&0x80 != 0)

		// If the result is greater than 0x99, we need to add 96 to it
		if result > 0x99 {
			result += 96
		}

		// Set the carry flag if the result is greater than 0x99
		cpu.setFlag(Carry, result > 0x99)
	} else {
		// Set the negative flag to opposite of (result & 0x80)
		cpu.setFlag(Negative, result&0x80 == 0x80)
		// Set overflow flag to opposite of ((registers.A ^ result) & (value ^ result) & 0x80)
		cpu.setFlag(Overflow, ^uint16(cpu.A)^result&^uint16(data)^result&0x80 != 0)
		// Set carry flag to opposite of (result > 0xFF)
		cpu.setFlag(Carry, result > 0xFF)
	}
	// Set the accumulator to the result
	cpu.A = uint8(result & 0xFF)
}

func (cpu *CPU) aslAcc() {
	// Shift the accumulator left
	result := uint16(cpu.A) << 1
	// Set the carry flag to opposite of (result > 0xFF)
	cpu.setFlag(Carry, result > 0xFF)
	// Set the zero flag to opposite of (result & 0xFF)
	cpu.setFlag(Zero, result&0xFF == 0)
	// Set the negative flag to opposite of (result & 0x80)
	cpu.setFlag(Negative, result&0x80 == 0x80)
	// Set the accumulator to the result
	cpu.A = uint8(result & 0xFF)
}

func (cpu *CPU) asl(address uint16) {
	// Fetch the data from the address
	data := cpu.MMU.readByte(address)
	// Shift the data left
	result := uint16(data) << 1
	// Set the carry flag to opposite of (result > 0xFF)
	cpu.setFlag(Carry, result > 0xFF)
	// Set the zero flag to opposite of (result & 0xFF)
	cpu.setFlag(Zero, result&0xFF == 0)
	// Set the negative flag to opposite of (result & 0x80)
	cpu.setFlag(Negative, result&0x80 == 0x80)
	// Write the result back to the address
	cpu.MMU.writeByte(address, uint8(result&0xFF))
}

func (cpu *CPU) brk() {
	// Increment the program counter
	cpu.PC++
	// Push the program counter to the stack
	cpu.pushWord(cpu.PC)
	// Push the status register to the stack
	cpu.pushByte(cpu.getStatus())
	// Set the break flag
	cpu.setFlag(Break, true)
	// Set the Interrupt flag
	cpu.setFlag(Interrupt, true)
	// Set the program counter to the address stored at 0xFFFE
	cpu.PC = cpu.MMU.readWord(0xFFFE)
}

func (cpu *CPU) clc() {
	cpu.setFlag(Carry, false)
}

func (cpu *CPU) cld() {
	cpu.setFlag(Decimal, false)
}

func (cpu *CPU) cli() {
	cpu.setFlag(Interrupt, false)
}

func (cpu *CPU) clv() {
	cpu.setFlag(Overflow, false)
}

func (cpu *CPU) inc(address uint16) {
	// Fetch the data from the address
	data := cpu.MMU.readByte(address)
	// Increment the data
	result := data + 1
	// Set the zero flag to opposite of (result & 0xFF)
	cpu.setFlag(Zero, result&0xFF == 0)
	// Set the negative flag to opposite of (result & 0x80)
	cpu.setFlag(Negative, result&0x80 == 0x80)
	// Write the result back to the address
	cpu.MMU.writeByte(address, result)
}

func (cpu *CPU) nop() {}

func (cpu *CPU) jmp(address uint16) {
	cpu.PC = address
}

func (cpu *CPU) lda(address uint16) {
	data := cpu.MMU.readByte(address)
	// Set the accumulator to the data
	cpu.A = data
	// Set the zero flag if the data is zero
	if data == 0 {
		cpu.setFlag(Zero, true)
	} else {
		cpu.setFlag(Zero, false)
	}
	// Set the negative flag if the data is negative
	if data&0x80 == 0x80 {
		cpu.setFlag(Negative, true)
	} else {
		cpu.setFlag(Negative, false)
	}
}

func (cpu *CPU) ldx(address uint16) {
	data := cpu.MMU.readByte(address)
	// Set the X register to the data
	cpu.X = data
	// Set the zero flag if the data is zero
	if data == 0 {
		cpu.setFlag(Zero, true)
	} else {
		cpu.setFlag(Zero, false)
	}
	// Set the negative flag if the data is negative
	if data&0x80 == 0x80 {
		cpu.setFlag(Negative, true)
	} else {
		cpu.setFlag(Negative, false)
	}
}

func ror_wrapped(cpu *CPU, value uint8) uint8 {
	// Rotate the value right
	result := (value >> 1) | (value << 7)
	// Set the carry flag to the value of the least significant bit
	cpu.setFlag(Carry, value&0x01 == 0x01)
	// Set the zero flag if the result is zero
	cpu.setFlag(Zero, result == 0)
	// Set the negative flag if the result is negative
	cpu.setFlag(Negative, result&0x80 == 0x80)
	// Return the result
	return result
}

func (cpu *CPU) ror(address uint16) {
	if address == 0 {
		cpu.A = ror_wrapped(cpu, cpu.A)
	} else {
		data := cpu.MMU.readByte(address)
		result := ror_wrapped(cpu, data)
		cpu.MMU.writeByte(address, result)
	}
}

func (cpu *CPU) sec() {
	cpu.setFlag(Carry, true)
}

func (cpu *CPU) sed() {
	cpu.setFlag(Decimal, true)
}

func (cpu *CPU) sei() {
	cpu.setFlag(Interrupt, true)
}

func (cpu *CPU) sta(address uint16) {
	cpu.MMU.writeByte(address, cpu.A)
}

func (cpu *CPU) stx(address uint16) {
	cpu.MMU.writeByte(address, cpu.X)
}
