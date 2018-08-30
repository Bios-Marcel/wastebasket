// +build darwin

package wastebasket

import (
	"fmt"
	"os"
	"os/exec"
)

//Trash moves a files or folder including its content into the systems trashbin.
func Trash(path string) error {
	_, fileError := os.Stat(path)

	if os.IsNotExist(fileError) {
		return nil
	}

	if fileError != nil {
		return fileError
	}

	command := fmt.Sprintln("tell app \"Finder\" to delete \"%s\" as POSIX file", path)
	return exec.Command("osascript", "-e", command).Run()
}

//Empty clears the platforms trashbin. It uses the `Finder` app to empty the trashbin.
func Empty() error {
	return exec.Command("osascript", "-e", "tell app \"Finder\" to empty").Run()
}
