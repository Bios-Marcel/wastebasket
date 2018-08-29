// +build darwin

package wastebasket

import (
	"errors"
	"os/exec"
)

//Trash moves a files or folder including its content into the systems trashbin.
func Trash(path string) error {
	return exec.Command("trash", path).Run()
}

//Empty clears the platforms trashbin.
func Empty() error {
	return errors.New("Not supported yet")
}
