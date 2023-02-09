package main

const (
	RAMSize = 0x10000 // 64KB
)

type MMU struct {
	RAM [RAMSize]uint8
	// TODO: add PPU, APU, etc.
}

// Read a byte from the memory
func (mmu *MMU) readByte(address uint16) uint8 {
	// return the value at the address
	return mmu.RAM[address]
}

// Read a word from the memory
func (mmu *MMU) readWord(address uint16) uint16 {
	// return a 16-bit word from the memory
	return uint16(mmu.RAM[address]) | uint16(mmu.RAM[address+1])<<8
}

func (mmu *MMU) writeByte(address uint16, value uint8) {
	// TODO: add PPU, APU, etc.
	mmu.RAM[address] = value
}

// Write a word to the memory
func (mmu *MMU) writeWord(address uint16, value uint16) {
	mmu.writeByte(address, uint8(value))
	mmu.writeByte(address+1, uint8(value>>8))
}

// Load a program into the memory
func (mmu *MMU) loadProgram(program []uint8) {
	for i, instruction := range program {
		mmu.RAM[0x8000+uint16(i)] = instruction
	}
}
