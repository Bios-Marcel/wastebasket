package wastebasket

import (
	"errors"
	"os"
	"os/exec"
)

var (
	availabilityCache   = make(map[string]bool)
	errToolNotAvailable = errors.New("tool not available")
	errNoToolsAvailable = errors.New("none of the commands `gio`, `gvfs-trash` or `trash` are available")
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

	return errNoToolsAvailable
}

func gioTrash(paths ...string) error {
	if isCommandAvailable("gio") {
		// --force makes sure we don't get errors for non-existent files.
		parameters := append([]string{"trash", "--force"}, paths...)
		return exec.Command("gio", parameters...).Run()
	}

	return errToolNotAvailable
}

func gvfsTrash(paths ...string) error {
	if isCommandAvailable("gvfs-trash") {
		// --force makes sure we don't get errors for non-existent files.
		parameters := append([]string{"--force"}, paths...)
		return exec.Command("gvfs-trash", parameters...).Run()
	}

	return errToolNotAvailable
}

func trashCli(paths ...string) error {
	if isCommandAvailable("trash") {
		//trash-cli throws 74 in case the file doesn't exist
		existingFiles := make([]string, 0, len(paths))
		for _, path := range paths {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				continue
			}

			existingFiles = append(existingFiles, path)
		}

		parameters := append([]string{"--"}, existingFiles...)
		return exec.Command("trash", parameters...).Run()
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

	return errNoToolsAvailable
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
