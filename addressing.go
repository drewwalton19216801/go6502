package main

// addressingModes is a map of addressing mode functions
var addressingModes = map[int]func(*CPU) uint16{
	// Implied
	0: func(cpu *CPU) uint16 { return 0 }, // no operand
	// Accumulator
	1: func(cpu *CPU) uint16 { return 0 }, // no operand
	// Immediate
	2: func(cpu *CPU) uint16 {
		// Read the address from the next byte
		address := cpu.PC
		// Increment the program counter
		cpu.PC++
		// Return the address
		return address
	},
	// Zero Page
	3: func(cpu *CPU) uint16 {
		// Read the address from the next byte
		address := cpu.MMU.readByte(cpu.PC)
		// Increment the program counter
		cpu.PC++
		// Return the address
		return uint16(address)
	},
	// Zero Page, X
	4: func(cpu *CPU) uint16 {
		// Read the address from the next byte
		address := cpu.MMU.readByte(cpu.PC)
		// Increment the program counter
		cpu.PC++
		// Return the address with the X register
		return uint16(address + cpu.X)
	},
	// Zero Page, Y
	5: func(cpu *CPU) uint16 {
		// Read the address from the next byte
		address := cpu.MMU.readByte(cpu.PC)
		// Increment the program counter
		cpu.PC++
		// Return the address with the Y register
		return uint16(address + cpu.Y)
	},
	// Relative
	6: func(cpu *CPU) uint16 {
		// Read the offset from the next byte
		offset := cpu.MMU.readByte(cpu.PC)
		// Increment the program counter
		cpu.PC++
		// Return the program counter with the offset
		return uint16(int16(cpu.PC) + int16(int8(offset)))
	},
	// Absolute
	7: func(cpu *CPU) uint16 {
		// Read the address from the next two bytes
		address := cpu.MMU.readWord(cpu.PC)
		// Increment the program counter by two
		cpu.PC += 2
		// Return the address
		return address
	},
	// Absolute, X
	8: func(cpu *CPU) uint16 {
		// Read the address from the next two bytes
		baseAddress := cpu.MMU.readWord(cpu.PC)
		// Increment the program counter by two
		cpu.PC += 2
		// Add the X register to the address
		effectiveAddress := baseAddress + uint16(cpu.X)

		// Check for page boundary crossing
		if effectiveAddress>>8 != baseAddress>>8 {
			// Page boundary crossed, so increment the cycle count
			cpu.cycles++
		}
		// Return the effective address
		return effectiveAddress
	},
	// Absolute, Y
	9: func(cpu *CPU) uint16 {
		// Read the address from the next two bytes
		baseAddress := cpu.MMU.readWord(cpu.PC)
		// Increment the program counter by two
		cpu.PC += 2
		// Add the Y register to the address
		effectiveAddress := baseAddress + uint16(cpu.Y)

		// Check for page boundary crossing
		if effectiveAddress>>8 != baseAddress>>8 {
			// Page boundary crossed, so increment the cycle count
			cpu.cycles++
		}
		// Return the effective address
		return effectiveAddress
	},
	// Indirect
	10: func(cpu *CPU) uint16 {
		// Read the address from the next two bytes
		address := cpu.MMU.readWord(cpu.PC)
		// Increment the program counter by two
		cpu.PC += 2

		// Check for page boundary crossing
		if address&0xFF == 0xFF {
			// Page boundary crossed, so read the address from the next page (fixes JMP ($xxFF) bug)
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
			// Address is in the page boundary, so read the address from the next page (fixes JMP ($xxFF) bug)
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
			// Address is in the page boundary, so read the address from the next page (fixes JMP ($xxFF) bug)
			low := cpu.MMU.readByte(address & 0xFF00)
			high := cpu.MMU.readByte((address & 0xFF00) | ((address + uint16(cpu.Y)) & 0xFF))
			return uint16(high)<<8 | uint16(low)
		} else {
			// Address is not in the page boundary, so read the address from the same page
			return cpu.MMU.readWord(address + uint16(cpu.Y))
		}
	},
}

// addressingModeNames is a map of addressing mode names
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
