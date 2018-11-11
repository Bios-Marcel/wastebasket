package wastebasket

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

//Trash moves a file or folder including its content into the systems trashbin.
func Trash(path string) error {
	_, fileError := os.Stat(path)

	if os.IsNotExist(fileError) {
		return nil
	}

	if fileError != nil {
		return fileError
	}

	//Passing a relative path will lead to the Finder not being able to find the file at all.
	path, pathToAbsPathError := filepath.Abs(path)
	if pathToAbsPathError != nil {
		return pathToAbsPathError
	}

	osascriptCommand := fmt.Sprintf("tell app \"Finder\" to delete POSIX file \"%s\"", path)
	return exec.Command("osascript", "-e", osascriptCommand).Run()
}

//Empty clears the platforms trashbin. It uses the `Finder` app to empty the trashbin.
func Empty() error {
	return exec.Command("osascript", "-e", "tell app \"Finder\" to empty").Run()
}
