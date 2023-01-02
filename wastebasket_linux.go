package wastebasket

import (
	"errors"
	"os"
	"os/exec"
)

var (
	availabilityCache   = make(map[string]bool)
	errToolNotAvailable = errors.New("tool not available")
)

func isCommandAvailable(name string) bool {
	avail, ok := availabilityCache[name]
	if avail && ok {
		return true
	}
	_, fileError := exec.LookPath(name)
	availabilityCache[name] = fileError == nil
	return fileError == nil
}

// Trash moves a file or folder including its content into the systems trashbin.
func Trash(path string) error {
	//gio is the tool that replaces gvfs, therefore it is the first choice.
	if err := gioTrash(path); err != nil && err != errToolNotAvailable {
		return err
	} else if err == nil {
		return nil
	}

	if err := gvfsTrash(path); err != nil && err != errToolNotAvailable {
		return err
	} else if err == nil {
		return nil
	}

	if err := trashCli(path); err != nil && err != errToolNotAvailable {
		return err
	} else if err == nil {
		return nil
	}

	return errors.New("None of the commands `gio`, `gvfs-trash` or `trash` are available")
}

func gioTrash(path string) error {
	if isCommandAvailable("gio") {
		// --force makes sure we don't get errors for non-existent files.
		return exec.Command("gio", "trash", "--force", path).Run()
	}

	return errToolNotAvailable
}

func gvfsTrash(path string) error {
	if isCommandAvailable("gvfs-trash") {
		// --force makes sure we don't get errors for non-existent files.
		return exec.Command("gvfs-trash", "--force", path).Run()
	}

	return errToolNotAvailable
}

func trashCli(path string) error {
	if isCommandAvailable("trash") {
		//trash-cli throws 74 in case the file doesn't exist
		_, fileError := os.Stat(path)

		if os.IsNotExist(fileError) {
			return nil
		}

		return exec.Command("trash", "--", path).Run()
	}

	return errToolNotAvailable
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

func gioEmpty() error {
	if isCommandAvailable("gio") {
		return exec.Command("gio", "trash", "--empty").Run()
	}

	return errToolNotAvailable
}

func gvfsEmpty() error {
	if isCommandAvailable("gvfs-trash") {
		return exec.Command("gvfs-trash", "--empty").Run()
	}

	return errToolNotAvailable
}

func trashCliEmpty() error {
	if isCommandAvailable("trash-empty") {
		return exec.Command("trash-empty").Run()
	}

	return errToolNotAvailable
}
