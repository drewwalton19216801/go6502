package main

// Flag bits
const (
	None      byte = 0         // No flags set
	Carry     byte = 1 << iota // Carry
	Zero                       // Zero
	Interrupt                  // Interrupt
	Decimal                    // Decimal
	Break                      // Break
	Overflow                   // Overflow
	Negative                   // Negative
)

func (cpu *CPU) setFlag(flag byte, value bool) {
	if value {
		cpu.P |= flag
	} else {
		cpu.P &^= flag
	}
}

func (cpu *CPU) getFlag(flag byte) bool {
	return cpu.P&flag != 0
}
