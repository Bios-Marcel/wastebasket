//go:build !windows
// +build !windows

package wastebasket

import (
	"os"
	"testing"
)

// TestTrashWithExistentFileWithDoubleQuotes tests trashing a single file with a double quote in its name
func TestTrashWithExistentFileWithDoubleQuotes(t *testing.T) {
	defer cleanup()

	writeTestFile(testFilePathWithDoubleQuotes, t)
	if errTrash := Trash(testFilePathWithDoubleQuotes); errTrash != nil {
		t.Errorf("Error trashing file. (%s)", errTrash.Error())
	}

	if _, errStat := os.Stat(testFilePathWithDoubleQuotes); os.IsNotExist(errStat) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted.")
	}
}
