package main

import "strconv"

// mhzToHz converts a frequency in MHz to Hz
func mhzToHz(mhz float64) int64 {
	return int64(mhz * 1000000)
}

// hzToMHz converts a frequency in Hz to MHz
func hzToMHz(hz int64) string {
	return strconv.FormatFloat(float64(hz)/1000000, 'f', 5, 64)
}

func boolToInt(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
