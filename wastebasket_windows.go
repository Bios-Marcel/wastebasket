// +build windows

package wastebasket

import (
	"fmt"
	"os"
	"os/exec"
)

//Trash moves a files or folder including its content into the systems trashbin.
func Trash(path string) error {
	file, fileError := os.Stat(path)
	if fileError != nil {
		return fileError
	}

	if os.IsNotExist(fileError) {
		return nil
	}

	psCommand := ""
	if file.IsDir() {
		psCommand = fmt.Sprintf("\"Add-Type -AssemblyName Microsoft.VisualBasic;[Microsoft.VisualBasic.FileIO.FileSystem]::DeleteFolder(%s, 'OnlyErrorDialogs','SendToRecycleBin')\"", path)
	} else {
		psCommand = fmt.Sprintf("\"Add-Type -AssemblyName Microsoft.VisualBasic;[Microsoft.VisualBasic.FileIO.FileSystem]::DeleteFile(%s, 'OnlyErrorDialogs','SendToRecycleBin')\"", path)
	}

	return exec.Command("powershell", "-Command", psCommand).Run()
}

//Empty clears the platforms trashbin.
func Empty() error {
	return exec.Command("powershell", "-Command", "\"Clear-RecycleBin\"").Run()
}
