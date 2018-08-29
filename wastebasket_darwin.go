// +build darwin

package wastebasket

import (
	"errors"
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

	return exec.Command("trash", path).Run()
}

//Empty clears the platforms trashbin.
func Empty() error {
	return errors.New("Not supported yet")
}
