package main

import "strconv"

func mhzToHz(mhz float64) int64 {
	return int64(mhz * 1000000)
}

func hzToMHz(hz int64) string {
	return strconv.FormatFloat(float64(hz)/1000000, 'f', 5, 64)
}
