//go:build linux

package wastebasket

import (
	"os"
	"testing"
)

func Test_Empty_Gio(t *testing.T) {
	err := gioEmpty()
	if err == errToolNotAvailable {
		t.SkipNow()
	}

	if err != nil {
		t.Errorf("unexpected error clearing trash: %s", err)
	}
}

func Test_Empty_Gvfs(t *testing.T) {
	err := gvfsEmpty()
	if err == errToolNotAvailable {
		t.SkipNow()
	}

	if err != nil {
		t.Errorf("unexpected error clearing trash: %s", err)
	}
}

func Test_Empty_TrashCli(t *testing.T) {
	err := trashCliEmpty()
	if err == errToolNotAvailable {
		t.SkipNow()
	}

	if err != nil {
		t.Errorf("unexpected error clearing trash: %s", err)
	}
}

func Test_Trash_DifferentTools(t *testing.T) {
	type testCase struct {
		name     string
		path     string
		fnDelete func(string) error
		err      error
	}

	cases := []testCase{
		{
			name:     "basic gio test",
			path:     "test.txt",
			fnDelete: gioTrash,
			err:      nil,
		},
		{
			name:     "basic gvfs test",
			path:     "test.txt",
			fnDelete: gvfsTrash,
			err:      nil,
		},
		{
			name:     "basic trashCli test",
			path:     "test.txt",
			fnDelete: trashCli,
			err:      nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer writeTestData(t, c.path)()

			err := c.fnDelete(c.path)
			if err == errToolNotAvailable {
				t.SkipNow()
			}

			if err != c.err {
				t.Errorf("unexpected error: %v != %v", err, c.err)
			}

			if _, err := os.Stat(c.path); os.IsNotExist(err) {
				//Everything correct!
			} else {
				t.Errorf("File hasn't been deleted.")
			}
		})
	}
}
