package main

import (
	"fmt"
	"strconv"
	"time"
)

type CPU struct {
	A, X, Y, P uint8
	PC         uint16
	SP         uint16
	clockSpeed int64 // in Hz
	MMU        *MMU  // memory management unit
	running    bool  // is the CPU running?
	cycles     int   // number of cycles executed
}

func (cpu *CPU) reset(resetVector uint16) {
	cpu.PC = resetVector
	cpu.SP = 0x0100
	cpu.A = 0x00
	cpu.X = 0x00
	cpu.Y = 0x00
	// Set None flags
	cpu.setFlag(None, true)
	// fill the memory with FFs
	for i := 0; i < RAMSize; i++ {
		cpu.MMU.RAM[i] = 0xFF
	}
	cpu.running = true
}

func (cpu *CPU) fetchByte() uint8 {
	value := cpu.MMU.readByte(cpu.PC)
	// Increment the program counter
	cpu.PC++
	// Increment the cycles
	cpu.cycles++
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
	inst := instructions[instruction]
	mode, ok := addressingModes[inst.addressingMode]
	if !ok {
		return 0x0000
	}
	// Print the address mode
	fmt.Println("Addressing mode:", addressingModeNames[inst.addressingMode])
	return mode(cpu)
}

func (cpu *CPU) disassemble(instruction uint8) string {
	inst, ok := instructions[instruction]
	if !ok {
		return "$" + strconv.FormatUint(uint64(instruction), 16)
	}
	return inst.mnemonic
}

func (cpu *CPU) run() {
	// load the demo program
	cycleDuration := time.Duration(1000000000 / cpu.clockSpeed)
	for {
		instruction := cpu.fetchByte()
		operand := cpu.fetchOperand(instruction)
		fmt.Printf("%04X: %s %04X\n", cpu.PC, cpu.disassemble(instruction), operand)
		inst := instructions[instruction]
		inst.execute(cpu, operand)
		for i := 0; i < inst.cycles; i++ {
			start := time.Now()
			// TODO: perform any necessary operations
			elapsed := time.Since(start)
			if elapsed < cycleDuration {
				time.Sleep(cycleDuration - elapsed)
			}
			if !cpu.running {
				break
			}
			cpu.cycles += inst.cycles
		}
		if !cpu.running {
			break
		}
	}
}
