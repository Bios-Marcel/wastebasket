//go:build !windows
// +build !windows

package wastebasket

import (
	"os"
	"testing"
)

// TestTrashWithExistentFileWithDoubleQuotes tests trashing a single file with a double quote in its name
func Test_Trash_ExistentFileWithDoubleQuotes(t *testing.T) {

	path := "foo\"bar\".txt"
	defer writeTestData(t, path)

	if errTrash := Trash(path); errTrash != nil {
		t.Errorf("Error trashing file. (%s)", errTrash.Error())
	}

	if _, errStat := os.Stat(path); os.IsNotExist(errStat) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted.")
	}
}
