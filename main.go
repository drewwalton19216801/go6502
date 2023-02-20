package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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
	fmt.Println("  --watch-addresses\tWatch the specified addresses (comma separated)")
	fmt.Println("  -f, --file\t\tLoad a program from a file")
	fmt.Println("Example: go6502 -c 1 -f program.bin --watch-addresses 0x6000,0x6002")
}

func loadProgramFromFile(fileName string) []uint8 {
	// Load the program from a file as a byte array
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()
	// Get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return nil
	}
	fileSize := fileInfo.Size()
	// Read the file
	program := make([]uint8, fileSize)
	_, err = file.Read(program)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}
	return program
}

func main() {
	debug := false
	speed := mhzToHz(1)
	loadFromFile := false
	watchAddresses := false
	var addressesToWatch []uint16
	var program []uint8
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
			case "--watch-addresses":
				if i+1 < len(os.Args) {
					watchAddresses = true
					i++
					addresses := os.Args[i]
					// Split the addresses by comma
					addressesArray := strings.Split(addresses, ",")
					// Convert the addresses to uint16
					for _, address := range addressesArray {
						addressInt, err := strconv.ParseUint(address, 0, 16)
						if err != nil {
							fmt.Println("Invalid address:", address)
							return
						}
						addressesToWatch = append(addressesToWatch, uint16(addressInt))
					}
				} else {
					fmt.Println("Missing addresses")
					return
				}
			case "-f", "--file":
				if i+1 < len(os.Args) {
					i++
					loadFromFile = true
					program = loadProgramFromFile(os.Args[i])
				} else {
					fmt.Println("Missing file name")
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
	mmu := &MMU{}
	cpu := CPU{clockSpeed: speed, MMU: mmu, debug: debug} // 0.00001 MHz (10 hz)
	// Load the demo program, if not loading from a file
	if !loadFromFile {
		program = demoProgram
	}
	mmu.loadProgram(program)
	// Reset the CPU to the demo program
	cpu.reset()
	// If we did not load from a file, we need to set the program counter to 0x8000
	if !loadFromFile {
		cpu.PC = 0x8000
	}
	// Run the CPU
	cpu.run(watchAddresses, addressesToWatch)
	fmt.Println("Emulation done in", cpu.cycles, "cycles", "at", hzToMHz(cpu.clockSpeed), "MHz")
}
