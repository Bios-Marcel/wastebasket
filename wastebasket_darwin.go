// +build darwin

package wastebasket

import (
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

//Empty isn't supported on the darwin platform.
func Empty() error {
	return nil
}
