// +build windows

package wastebasket

import "errors"

//Trash moves a files or folder including its content into the systems trashbin.
func Trash(path string) error {
	return errors.New("Not supported yet")
}

//Empty clears the platforms trashbin.
func Empty() error {
	return errors.New("Not supported yet")
}
