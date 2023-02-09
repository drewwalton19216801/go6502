# go6502
A simple 6502 emulator written in Go

This may not ever be complete or useful, but it's nonetheless cool.

## Features
- [x] Decimal mode
- [X] Controllable clock speed
- [ ] Interrupts
- [ ] Non-maskable interrupts
- [ ] 100% legal instruction coverage
- [X] 100% legal addressing mode coverage
- [ ] 100% illegal instruction coverage
- [ ] 6502 variant support
- [ ] Loading ROMs from files

## Building
From the go6502 directory: `go build .`

## Running
`./go6502 [options]` (run `./go6502 --help` to see options)

Alternatively, you can run `go run . [options]`

## Options
`--clock-speed (-c)` - Clock speed in MHz (default 1, can go down to 0.00001)

`--debug (-d)` - Enable debug mode
