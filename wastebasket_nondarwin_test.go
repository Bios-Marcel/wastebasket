//go:build !darwin

package wastebasket_test

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Bios-Marcel/wastebasket/v2"
)

func Test_Query_Restore_Homedir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("single_file", func(t *testing.T) {
		path := filepath.Join(home, "path.txt")
		t.Cleanup(writeTestData(t, path))
		assertExists(t, path)

		if err := wastebasket.Trash(path); err != nil {
			t.Errorf("Error trashing file. (%s)", err.Error())
		}
		assertNotExists(t, path)

		result, err := wastebasket.Query(wastebasket.QueryOptions{Search: []string{path}})
		if assert.NoError(t, err) {
			if assert.Len(t, result.Matches[path], 1) {
				if assert.NoError(t, result.Matches[path][0].Restore(false)) {
					assertExists(t, path)
				}
			}
		}
	})

	// Makes sure restore requires force for overwrite.
	t.Run("same_file_twice", func(t *testing.T) {
		path := filepath.Join(home, "path.txt")

		for i := range 2 {
			t.Cleanup(writeTestDataWith(t, strconv.FormatInt(int64(i), 10), path))
			assertExists(t, path)
			if err := wastebasket.Trash(path); err != nil {
				t.Errorf("Error trashing file. (%s)", err.Error())
			}
			assertNotExists(t, path)
		}

		result, err := wastebasket.Query(wastebasket.QueryOptions{Search: []string{path}})
		require.NoError(t, err)
		require.Len(t, result.Matches[path], 2)
		matchOne := result.Matches[path][0]
		require.NoError(t, matchOne.Restore(false))
		f, err := os.ReadFile(matchOne.OriginalPath())
		require.NoError(t, err)
		var wasFirst bool
		if string(f) != "0" && string(f) != "1" {
			t.Fatalf("File content unexpected: %s", string(f))
		}
		if string(f) == "0" {
			wasFirst = true
		}

		assertExists(t, path)
		matchTwo := result.Matches[path][1]
		require.Error(t, matchTwo.Restore(false), wastebasket.ErrAlreadyExists)
		f, err = os.ReadFile(matchTwo.OriginalPath())
		require.NoError(t, err)
		if wasFirst {
			require.Equal(t, "0", string(f))
		} else {
			require.Equal(t, "1", string(f))
		}

		assertExists(t, path)
		require.NoError(t, matchTwo.Restore(true))
		f, err = os.ReadFile(matchTwo.OriginalPath())
		require.NoError(t, err)
		if wasFirst {
			require.Equal(t, "1", string(f))
		} else {
			require.Equal(t, "0", string(f))
		}
		assertExists(t, path)
	})
}
