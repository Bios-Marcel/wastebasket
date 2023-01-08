//go:build !windows && !darwin && !nix_wrapper

package wastebasket

import (
	"os"
)

func customImplTrash(paths ...string) error {
	for _, path := range paths {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}

	return nil
}

func customImplEmpty() error {
	return nil
}
