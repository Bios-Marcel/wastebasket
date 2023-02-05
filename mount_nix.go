//go:build !windowss && !darwin

package wastebasket

import (
	"bufio"
	"bytes"
	"math"
	"os"
	"strings"
	"sync/atomic"
)

// The mounts could change anytime, but probably won't all too often. So we
// read them each time Mount() is called, but remember the count, so we can
// at least prepare a buffer of that size. This must only be accessed atomicly.
var lastMountCount int32

// Mounts reads /proc/mounts and returns all mounts found, excluding
// everything on the devices sysfs, rootfs, cgroup and /dev/ paths.
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
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			if firstSpace := bytes.IndexByte(data, ' '); firstSpace >= 0 {
				// Some device won't contain a trash either way or might be
				// dangerous to interact with.
				device := string(data[:firstSpace])
				switch device {
				case "rootfs", "sysfs", "cgroup":
					// Skip line
					return i + 1, nil, nil
				default:
					// Devices don't usually contain files
					if strings.HasPrefix(device, "/dev/") {
						return i + 1, nil, nil
					}
				}

				if nextSpace := bytes.IndexByte(data[firstSpace+1:], ' '); nextSpace >= 0 {
					return i + 1, data[firstSpace+1 : firstSpace+1+nextSpace], nil
				}
			}

			return i + 1, nil, nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	})
	scanner.Buffer(buffer, math.MaxInt)
	mounts := make([]string, 0, atomic.LoadInt32(&lastMountCount))

	for scanner.Scan() {
		mounts = append(mounts, scanner.Text())
	}

	atomic.StoreInt32(&lastMountCount, int32(len(mounts)))

	return mounts, nil
}
