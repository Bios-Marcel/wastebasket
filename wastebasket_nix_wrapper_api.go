//go:build !windows && !darwin && nix_wrapper

package wastebasket

import (
	"errors"
)

// Trash moves a file or folder including its content into the systems trashbin.
func Trash(paths ...string) error {
	//gio is the tool that replaces gvfs, therefore it is the first choice.
	if err := gioTrash(paths...); err != nil && err != errToolNotAvailable {
		return err
	} else if err == nil {
		return nil
	}

	if err := gvfsTrash(paths...); err != nil && err != errToolNotAvailable {
		return err
	} else if err == nil {
		return nil
	}

	if err := trashCli(paths...); err != nil && err != errToolNotAvailable {
		return err
	} else if err == nil {
		return nil
	}

	return errors.New("None of the commands `gio`, `gvfs-trash` or `trash` are available")
}

// Empty clears the platforms trashbin.
func Empty() error {
	//gio is the tool that replaces gvfs, therefore it is the first choice.
	if err := gioEmpty(); err != nil && err != errToolNotAvailable {
		return err
	} else if err == nil {
		return nil
	}

	if err := gvfsEmpty(); err != nil && err != errToolNotAvailable {
		return err
	} else if err == nil {
		return nil
	}

	if err := trashCliEmpty(); err != nil && err != errToolNotAvailable {
		return err
	} else if err == nil {
		return nil
	}

	return errors.New("None of the commands `gio`, `gvfs-trash` or `trash-empty` are available")
}
