package main

var addressingModes = map[int]func(*CPU) uint16{
	// Implied
	0: func(cpu *CPU) uint16 { return 0 }, // no operand
	// Accumulator
	1: func(cpu *CPU) uint16 { return 0 }, // no operand
	// Immediate
	2: func(cpu *CPU) uint16 {
		// Read the address from the next byte
		address := cpu.PC
		cpu.PC++
		return address
	},
	// Zero Page
	3: func(cpu *CPU) uint16 {
		address := cpu.MMU.readByte(cpu.PC)
		cpu.PC++
		return uint16(address)
	},
	// Zero Page, X
	4: func(cpu *CPU) uint16 {
		address := cpu.MMU.readByte(cpu.PC)
		cpu.PC++
		return uint16(address + cpu.X)
	},
	// Zero Page, Y
	5: func(cpu *CPU) uint16 {
		address := cpu.MMU.readByte(cpu.PC)
		cpu.PC++
		return uint16(address + cpu.Y)
	},
	// Relative
	6: func(cpu *CPU) uint16 {
		offset := cpu.MMU.readByte(cpu.PC)
		cpu.PC++
		return uint16(int16(cpu.PC) + int16(int8(offset)))
	},
	// Absolute
	7: func(cpu *CPU) uint16 {
		address := cpu.MMU.readWord(cpu.PC)
		cpu.PC += 2
		return address
	},
	// Absolute, X
	8: func(cpu *CPU) uint16 {
		baseAddress := cpu.MMU.readWord(cpu.PC)
		cpu.PC += 2
		effectiveAddress := baseAddress + uint16(cpu.X)

		// Check for page boundary crossing
		if effectiveAddress>>8 != baseAddress>>8 {
			cpu.cycles++
		}
		// Return the effective address
		return effectiveAddress
	},
	// Absolute, Y
	9: func(cpu *CPU) uint16 {
		baseAddress := cpu.MMU.readWord(cpu.PC)
		cpu.PC += 2
		effectiveAddress := baseAddress + uint16(cpu.Y)

		// Check for page boundary crossing
		if effectiveAddress>>8 != baseAddress>>8 {
			cpu.cycles++
		}
		// Return the effective address
		return effectiveAddress
	},
	// Indirect
	10: func(cpu *CPU) uint16 {
		address := cpu.MMU.readWord(cpu.PC)
		cpu.PC += 2

		// Check for page boundary crossing
		if address&0xFF == 0xFF {
			return uint16(cpu.MMU.readByte(address)) | uint16(cpu.MMU.readByte(address&0xFF00))<<8
		}
		return cpu.MMU.readWord(address)
	},
	// (Indirect, X)
	11: func(cpu *CPU) uint16 {
		// Read the address from the next byte and add the X register
		address := (uint16(cpu.MMU.readByte(cpu.PC) + cpu.X))
		cpu.PC++
		// Read the address from the same page
		address = cpu.MMU.readWord(address)

		// Check if the address is in the page boundary
		if (address+uint16(cpu.X))&0xFF < address&0xFF {
			// Address is in the page boundary, so read the address from the next page
			low := cpu.MMU.readByte(address & 0xFF00)
			high := cpu.MMU.readByte((address & 0xFF00) | ((address + uint16(cpu.X)) & 0xFF))
			return uint16(high)<<8 | uint16(low)
		} else {
			// Address is not in the page boundary, so read the address from the same page
			return cpu.MMU.readWord(address + uint16(cpu.X))
		}
	},
	// (Indirect), Y
	12: func(cpu *CPU) uint16 {
		// Read the address from the next byte
		address := uint16(cpu.MMU.readByte(cpu.PC))
		cpu.PC++
		// Read the address from the same page
		address = cpu.MMU.readWord(address)

		// Check if the address is in the page boundary
		if (address+uint16(cpu.Y))&0xFF < address&0xFF {
			// Address is in the page boundary, so read the address from the next page
			low := cpu.MMU.readByte(address & 0xFF00)
			high := cpu.MMU.readByte((address & 0xFF00) | ((address + uint16(cpu.Y)) & 0xFF))
			return uint16(high)<<8 | uint16(low)
		} else {
			// Address is not in the page boundary, so read the address from the same page
			return cpu.MMU.readWord(address + uint16(cpu.Y))
		}
	},
}

var addressingModeNames = map[int]string{
	0:  "Implied",
	1:  "Accumulator",
	2:  "Immediate",
	3:  "ZeroPage",
	4:  "ZeroPageX",
	5:  "ZeroPageY",
	6:  "Relative",
	7:  "Absolute",
	8:  "AbsoluteX",
	9:  "AbsoluteY",
	10: "Indirect",
	11: "IndirectX",
	12: "IndirectY",
}
