package wastebasket

import (
	"errors"
	"os"
	"os/exec"
)

func isCommandAvailable(name string) bool {
	_, fileError := exec.LookPath(name)
	return fileError == nil
}

//Trash moves a file or folder including its content into the systems trashbin.
func Trash(path string) error {
	//gio us the tool that replaces gvfs, therefore it is the first choice.
	if isCommandAvailable("gio") {
		return exec.Command("gio", "trash", "--force", path).Run()
	} else if isCommandAvailable("gvfs-trash") {
		return exec.Command("gvfs-trash", "--force", path).Run()
	} else if isCommandAvailable("trash") {
		//trash-cli throws 74 in case the file doesn't exist
		_, fileError := os.Stat(path)

		if os.IsNotExist(fileError) {
			return nil
		}

		return exec.Command("trash", "--", path).Run()
	}

	return errors.New("None of the commands `gio`, `gvfs-trash` or `trash` are available")
}

//Empty clears the platforms trashbin.
func Empty() error {
	//gio us the tool that replaces gvfs, therefore it is the first choice.
	if isCommandAvailable("gio") {
		return exec.Command("gio", "trash", "--empty").Run()
	} else if isCommandAvailable("gvfs-trash") {
		return exec.Command("gvfs-trash", "--empty").Run()
	} else if isCommandAvailable("trash-empty") {
		return exec.Command("trash-empty").Run()
	}

	return errors.New("None of the commands `gio`, `gvfs-trash` or `trash-empty` are available")
}
