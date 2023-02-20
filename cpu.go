package main

import (
	"fmt"
	"time"
)

var resetVector = 0xFFFC

type CPU struct {
	A, X, Y, P uint8
	PC         uint16
	SP         uint16
	clockSpeed int64 // in Hz
	MMU        *MMU  // memory management unit
	running    bool  // is the CPU running?
	cycles     int   // number of cycles executed
	debug      bool  // is the CPU in debug mode?
}

func (cpu *CPU) reset() {
	// Reset the CPU
	cpu.PC = cpu.MMU.readWord(uint16(resetVector))
	cpu.SP = 0x0100
	cpu.A = 0x00
	cpu.X = 0x00
	cpu.Y = 0x00
	// Set None flags
	cpu.setFlag(None, true)
	cpu.running = true
}

func (cpu *CPU) writeByte(value uint8) {
	// Write the value to the memory
	cpu.MMU.writeByte(cpu.PC, value)
	// Increment the program counter
	cpu.PC++
}

func (cpu *CPU) writeWord(value uint16) {
	// Write the low byte
	cpu.writeByte(uint8(value))
	// Write the high byte
	cpu.writeByte(uint8(value >> 8))
}

func (cpu *CPU) fetchByte() uint8 {
	// Read the byte
	value := cpu.MMU.readByte(cpu.PC)
	// Increment the program counter
	cpu.PC++
	// Return the value
	return value
}

func (cpu *CPU) fetchSignedByte() int8 {
	// Read the byte
	value := cpu.fetchByte()
	// Convert it to a signed byte
	return int8(value)
}

func (cpu *CPU) fetchWord() uint16 {
	// Read the low byte
	low := cpu.fetchByte()
	// Read the high byte
	high := cpu.fetchByte()
	// Return the word
	return uint16(high)<<8 | uint16(low)
}

func (cpu *CPU) spToAddress() uint16 {
	// Convert the stack pointer to an address
	return 0x0100 | uint16(cpu.SP)
}

func (cpu *CPU) setZNFlags() {
	// Set the Zero flag
	cpu.setFlag(Zero, cpu.A == 0)
	// Set the Negative flag
	cpu.setFlag(Negative, cpu.A&0x80 != 0)
}

func (cpu *CPU) pushByte(value uint8) {
	// Decrement the stack pointer
	cpu.SP--
	// Write the value to the stack
	cpu.MMU.writeByte(cpu.spToAddress(), value)
}

func (cpu *CPU) pushWord(value uint16) {
	// Push the high byte
	cpu.pushByte(uint8(value >> 8))
	// Push the low byte
	cpu.pushByte(uint8(value))
}

func (cpu *CPU) popByte() uint8 {
	// Increment the stack pointer
	cpu.SP++
	// Read the value from the stack
	return cpu.MMU.readByte(cpu.spToAddress())
}

func (cpu *CPU) popWord() uint16 {
	// Pop the low byte
	low := cpu.popByte()
	// Pop the high byte
	high := cpu.popByte()
	// Return the word
	return uint16(high)<<8 | uint16(low)
}

func (cpu *CPU) pushPC() {
	// Push the program counter to the stack
	cpu.pushWord(cpu.PC)
}

func (cpu *CPU) pushPCPlusOne() {
	// Push the program counter + 1 to the stack
	cpu.pushWord(cpu.PC + 1)
}

func (cpu *CPU) pushPCMinusOne() {
	// Push the program counter - 1 to the stack
	cpu.pushWord(cpu.PC - 1)
}

func (cpu *CPU) fetchOperand(instruction uint8) uint16 {
	// Get the addressing mode
	inst := instructions[instruction]
	// Get the addressing mode function
	mode, ok := addressingModes[inst.addressingMode]
	if !ok {
		return 0x0000
	}
	// Return the operand
	return mode(cpu)
}

func (cpu *CPU) disassemble(instruction uint8) string {
	// Get the instruction
	inst := instructions[instruction]
	// Get the addressing mode
	mode := addressingModeNames[inst.addressingMode]

	operandString := ""

	// Get the operand string
	switch mode {
	case "Implied":
		operandString = ""
	case "Accumulator":
		operandString = "A"
	case "Immediate":
		operandString = fmt.Sprintf("#$%02X", cpu.MMU.readByte(cpu.PC))
	case "ZeroPage":
		operandString = fmt.Sprintf("$%02X", cpu.MMU.readByte(cpu.PC))
	case "ZeroPageX":
		operandString = fmt.Sprintf("$%02X,X", cpu.MMU.readByte(cpu.PC))
	case "ZeroPageY":
		operandString = fmt.Sprintf("$%02X,Y", cpu.MMU.readByte(cpu.PC))
	case "Relative":
		operandString = fmt.Sprintf("$%02X", cpu.MMU.readByte(cpu.PC))
	case "Absolute":
		operandString = fmt.Sprintf("$%04X", cpu.MMU.readWord(cpu.PC))
	case "AbsoluteX":
		operandString = fmt.Sprintf("$%04X,X", cpu.MMU.readWord(cpu.PC))
	case "AbsoluteY":
		operandString = fmt.Sprintf("$%04X,Y", cpu.MMU.readWord(cpu.PC))
	case "Indirect":
		operandString = fmt.Sprintf("($%04X)", cpu.MMU.readWord(cpu.PC))
	case "IndirectX":
		operandString = fmt.Sprintf("($%02X,X)", cpu.MMU.readByte(cpu.PC))
	case "IndirectY":
		operandString = fmt.Sprintf("($%02X),Y", cpu.MMU.readByte(cpu.PC))

	}
	// Return the disassembly
	return (inst.mnemonic + " " + operandString)
}

func (cpu *CPU) run(watch bool, watchAddresses []uint16) {
	// Set the running flag
	cpu.running = true
	// Set the cycle duration based on the clock speed
	cycleDuration := time.Duration(1000000000 / cpu.clockSpeed)
	for {
		// Fetch the instruction
		instruction := cpu.fetchByte()
		if cpu.debug {
			// Disassemble the instruction if debugging is enabled
			fmt.Println("CPU run: disassembly =", cpu.disassemble(instruction))
			// Print the CPU registers in hex
			fmt.Printf("CPU run: A = %02X, X = %02X, Y = %02X, SP = %02X, PC = %04X, P = %02X (%08b), Cycles = %d\n", cpu.A, cpu.X, cpu.Y, cpu.SP, cpu.PC, cpu.P, cpu.P, cpu.cycles)
		}
		// Fetch the operand
		operand := cpu.fetchOperand(instruction)
		// Get the instruction from the instruction map
		inst := instructions[instruction]
		// Execute the instruction
		inst.execute(cpu, operand)
		// If the watch flag is set, print the memory addresses
		if watch {
			for _, address := range watchAddresses {
				fmt.Printf("CPU run: $%04X = %02X\n", address, cpu.MMU.readByte(address))
			}
		}
		for i := 0; i < inst.cycles; i++ {
			start := time.Now()
			// TODO: perform any necessary operations
			elapsed := time.Since(start)
			// Check if the elapsed time is less than the cycle duration
			if elapsed < cycleDuration {
				// Sleep for the remaining time
				time.Sleep(cycleDuration - elapsed)
			}
		}
		// Increment the cycle count by the number of cycles the instruction takes
		cpu.cycles += inst.cycles
		// Check if the CPU is running
		if !cpu.running {
			break
		}
		// Check if the break flag is set
		if cpu.getFlag(Break) {
			break
		}
	}
}
