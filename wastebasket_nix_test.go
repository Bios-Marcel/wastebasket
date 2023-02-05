//go:build !windows && !darwin

package wastebasket

import (
	"testing"
)

// TestTrashWithExistentFileWithDoubleQuotes tests trashing a single file with a double quote in its name
func Test_Trash_ExistentFileWithDoubleQuotes(t *testing.T) {
	path := "foo\"bar\".txt"
	defer writeTestData(t, path)()
	assertExists(t, path)

	if err := Trash(path); err != nil {
		t.Errorf("Error trashing file. (%s)", err.Error())
	}

	assertNotExists(t, path)
}
