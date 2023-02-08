package main

import (
	"fmt"
	"time"
)

var demoProgram = []uint8{
	// NOP
	0xEA,
	// SED
	0xF8,
	// LDA #$01
	0xA9, 0x01,
	// CLC
	0x18,
	// ADC #$01
	0x69, 0x01,
	// STA $0200
	0x8D, 0x00, 0x02,
	// CLD
	0xD8,
	// BRK
	0x00,
}

func main() {
	startTime := time.Now()
	mmu := &MMU{}
	cpu := CPU{clockSpeed: mhzToHz(0.00001), MMU: mmu, debug: true} // 0.00001 MHz (10 hz)
	// Reset the CPU to the demo program
	cpu.reset(0x8000)
	// Load the demo program
	mmu.loadProgram(demoProgram)
	cpu.run()
	fmt.Println("Emulation done in", cpu.cycles, "cycles", "at", hzToMHz(cpu.clockSpeed), "MHz")
	// Print the CPU registers on one line
	fmt.Printf("A: %02X X: %02X Y: %02X P: %02X SP: %02X PC: %04X\n", cpu.A, cpu.X, cpu.Y, cpu.P, cpu.SP, cpu.PC)
	// Print the memory at $0200
	fmt.Printf("Memory at $0200: %02X\n", mmu.readByte(0x0200))
	finishTime := time.Now()
	fmt.Println("Time elapsed:", finishTime.Sub(startTime))
}
