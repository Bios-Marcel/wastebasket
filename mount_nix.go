//go:build !windowss && !darwin

package wastebasket

import (
	"bufio"
	"math"
	"os"
	"sync/atomic"
	"unicode/utf8"
)

// The mounts could change anytime, but probably won't all too often. So we
// read them each time Mount() is called, but remember the count, so we can
// at least prepare a buffer of that size. This must only be accessed atomicly.
var lastMountCount int32

func Mounts() ([]string, error) {
	handle, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, err
	}

	defer handle.Close()

	scanner := bufio.NewScanner(handle)
	// I just assume that entries won't be much longer than 200, since my
	// own /proc/mounts was max 157 characters. Worst case, we'll probably
	// expand to 400. Either way, this should handle most cases efficiently.
	buffer := make([]byte, 0, 200)
	scanner.Split(bufio.ScanLines)
	scanner.Buffer(buffer, math.MaxInt)
	mounts := make([]string, 0, atomic.LoadInt32(&lastMountCount))

	for scanner.Scan() {
		bytes := scanner.Bytes()
		// Skip first character of device, it doesn't matter anyway.
		from := 1
		var startOfPath int
		for {
			r, size := utf8.DecodeRune(bytes[from:])
			from += size
			if r == ' ' {
				if startOfPath > 0 {
					mounts = append(mounts, string(bytes[startOfPath:from-1]))
					break
				}

				startOfPath = from
			}
		}
	}

	atomic.StoreInt32(&lastMountCount, int32(len(mounts)))

	return mounts, nil
}
