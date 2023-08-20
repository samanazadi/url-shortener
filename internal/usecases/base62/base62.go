package base62

import (
	"time"
)

// charSet is the mapping used in base62 encoding
const charSet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

const baseUnixMilli int64 = 1577836800000 // time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

var unixMilliMask int64
var lastUnixMilli int64
var lastSeq = uint8(0)

func init() {
	for i := 0; i < 39; i++ {
		unixMilliMask <<= 1
		unixMilliMask++
	}
}

func int64ToString(i int64) string {
	const base = int64(len(charSet))
	encoded := ""
	for i > 0 {
		encoded = string(charSet[i%base]) + encoded
		i /= base
	}
	return encoded
}

// schema: 0[39 bits milliseconds from base / 10][16 bits machine ID][8 bits seq]
func makeInt64(unixMilli int64, machineID uint16, seq uint8) int64 {
	id := int64(seq) // start with seq on the right
	id += int64(machineID) << 8
	diff10Milli := (unixMilli - baseUnixMilli) / 10
	id += (diff10Milli & unixMilliMask) << 24
	return id
}

// GenerateID creates new ID
func GenerateID(machineID uint16) string {
	um := time.Now().UnixMilli()
	var seq uint8
	if lastUnixMilli/10 != um/10 {
		lastSeq = 0 // restart seq
		seq = 0
	} else {
		lastSeq++
		seq = lastSeq
	}
	lastUnixMilli = um

	return int64ToString(makeInt64(um, machineID, seq))
}
