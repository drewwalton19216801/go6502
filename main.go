package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

var demoProgram = []uint8{
	// NOP
	0xEA, // 2 cycles
	// SED
	0xF8, // 2 cycles
	// LDA #$01
	0xA9, 0x01, // 2 cycles
	// CLC
	0x18, // 2 cycles
	// ADC #$03
	0x69, 0x03, // 2 cycles
	// STA $0200
	0x8D, 0x00, 0x02, // 4 cycles
	// LDX $0200 (load the value of $0200 into X)
	0xAE, 0x00, 0x02, // 4 cycles
	// STX $0300 (store the value of X into $0300)
	0x8E, 0x00, 0x03, // 4 cycles
	// CLD
	0xD8, // 2 cycles
	// BRK
	0x00, // 7 cycles
	// Expected cycles: 2 + 2 + 2 + 2 + 2 + 4 + 4 + 4 + 2 + 7 = 31
}

func printUsage() {
	fmt.Println("Usage: go6502 [options]")
	fmt.Println("Options:")
	fmt.Println("  -h, --help\t\tPrint this help message")
	fmt.Println("  -d, --debug\t\tEnable debug mode")
	fmt.Println("  -c, --clock-speed\tSet the clock speed in MHz")
}

func main() {
	debug := false
	speed := mhzToHz(1)
	// Parse the command line arguments
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "-h", "--help":
				printUsage()
				return
			case "-d", "--debug":
				debug = true
			case "-c", "--clock-speed":
				if i+1 < len(os.Args) {
					i++
					clockSpeed, err := strconv.ParseFloat(os.Args[i], 64)
					if err != nil {
						fmt.Println("Invalid clock speed:", os.Args[i])
						return
					}
					speed = mhzToHz(clockSpeed)
				} else {
					fmt.Println("Missing clock speed")
					return
				}
			default:
				fmt.Println("Invalid option:", os.Args[i])
				return
			}
		}
	} else {
		printUsage()
		return
	}
	startTime := time.Now()
	mmu := &MMU{}
	cpu := CPU{clockSpeed: speed, MMU: mmu, debug: debug} // 0.00001 MHz (10 hz)
	// Reset the CPU to the demo program
	cpu.reset(0x8000)
	// Load the demo program
	mmu.loadProgram(demoProgram)
	// Run the CPU
	cpu.run()
	fmt.Println("Emulation done in", cpu.cycles, "cycles", "at", hzToMHz(cpu.clockSpeed), "MHz")
	// Print the CPU registers on one line
	fmt.Printf("A: %02X X: %02X Y: %02X P: %02X SP: %02X PC: %04X\n", cpu.A, cpu.X, cpu.Y, cpu.P, cpu.SP, cpu.PC)
	// Print the memory at $0200
	fmt.Printf("Memory at $0200: %02X\n", mmu.readByte(0x0200))
	// Print the memory at $0300
	fmt.Printf("Memory at $0300: %02X\n", mmu.readByte(0x0300))
	// Print the time elapsed
	finishTime := time.Now()
	fmt.Println("Time elapsed:", finishTime.Sub(startTime))
}
