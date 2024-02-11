//go:build android || ios || js

package wastebasket

import (
	"os"
)

func Trash(paths ...string) error {
	for _, path := range paths {
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func Empty() error {
	return nil
}
