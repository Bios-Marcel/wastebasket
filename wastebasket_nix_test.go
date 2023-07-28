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

// FIXME Write tests for:
// * Restore on topdir of mount
// * Restore of file with multiple versions
// * Restore of file with multiple versions in different trashbins
//   (technically not possible if only storing with wastebasket, but can happen technically on a system)
// * Restore of nonexistent files
