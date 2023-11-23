package wastebasket_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Bios-Marcel/wastebasket/v2"
	"github.com/Bios-Marcel/wastebasket/v2/wastebasket_nix"
	"github.com/Bios-Marcel/wastebasket/v2/wastebasket_windows"
	"github.com/stretchr/testify/assert"
)

var (
	_ wastebasket.TrashedFileInfo = &wastebasket_nix.TrashedFileInfo{}
	_ wastebasket.TrashedFileInfo = &wastebasket_windows.TrashedFileInfo{}
)

func generateManyFileNames(count int) []string {
	fileNames := make([]string, 0, count)
	for i := 1; i <= count; i++ {
		fileNames = append(fileNames, fmt.Sprintf("%d.txt", i))
	}
	return fileNames
}

func Test_Trash(t *testing.T) {
	type trashExpectation struct {
		trasher     func() error
		expectedErr error
	}
	type testCase struct {
		name              string
		testDataToCreate  []string
		trashExpectations []trashExpectation
		expectedFiles     []string
	}

	cases := []testCase{
		{
			name:             "existent file, no edge cases",
			testDataToCreate: []string{"test.txt"},
			trashExpectations: []trashExpectation{
				{trash("test.txt"), nil},
			},
			expectedFiles: nil,
		},
		{
			name:             "existent file, spaces in name",
			testDataToCreate: []string{"te st.txt"},
			trashExpectations: []trashExpectation{
				{trash("te st.txt"), nil},
			},
			expectedFiles: nil,
		},
		{
			name:             "existent file, spaces in name and ./ when deleting",
			testDataToCreate: []string{"te st.txt"},
			trashExpectations: []trashExpectation{
				{trash("./te st.txt"), nil},
			},
			expectedFiles: nil,
		},
		{
			name:             "non-existent file",
			testDataToCreate: nil,
			trashExpectations: []trashExpectation{
				{trash("doesntexist.txt"), nil},
			},
			expectedFiles: nil,
		},
		{
			name:             "existent empty directory",
			testDataToCreate: []string{"folder/"},
			trashExpectations: []trashExpectation{
				{trash("folder"), nil},
			},
			expectedFiles: nil,
		},
		{
			name:             "existent non empty directory",
			testDataToCreate: []string{"folder/", "folder/file.txt"},
			trashExpectations: []trashExpectation{
				{trash("folder"), nil},
			},
			expectedFiles: nil,
		},
		{
			name:             "existent directory with spaces in name",
			testDataToCreate: []string{"fol der/"},
			trashExpectations: []trashExpectation{
				{trash("fol der"), nil},
			},
			expectedFiles: nil,
		},
		{
			name:             "delete two files in one call",
			testDataToCreate: []string{"a.txt", "b.txt"},
			trashExpectations: []trashExpectation{
				{trash("a.txt", "b.txt"), nil},
			},
			expectedFiles: nil,
		},
		{
			name:             "mix of folders and files",
			testDataToCreate: []string{"a/", "b.txt", "c/", "c/file-in-c.txt", "c/another-file-in-c.txt"},
			trashExpectations: []trashExpectation{
				{trash("a", "b.txt", "c/file-in-c.txt", "c"), nil},
			},
			expectedFiles: nil,
		},
		{
			name:             "large amount of files",
			testDataToCreate: generateManyFileNames(257),
			trashExpectations: []trashExpectation{
				{trash(generateManyFileNames(257)...), nil},
			},
			expectedFiles: nil,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			defer writeTestData(t, testCase.testDataToCreate...)()

			for _, expectation := range testCase.trashExpectations {
				err := expectation.trasher()
				if !errors.Is(err, expectation.expectedErr) {
					t.Errorf("unexpected error: %v != %v", err, expectation.expectedErr)
				}
			}

		OUTER_LOOP:
			for _, file := range testCase.testDataToCreate {
				_, err := os.Stat(file)
				for _, expectedFile := range testCase.expectedFiles {
					if file == expectedFile {
						if os.IsNotExist(err) {
							t.Errorf("File %s doesn't exist, but was expected to", file)
						}
						continue OUTER_LOOP
					}
				}

				if err == nil {
					t.Errorf("file %s shouldn't exist, but does", file)
				}
			}
		})
	}
}

func Test_Empty(t *testing.T) {
	if err := wastebasket.Empty(); err != nil {
		t.Errorf("Error emptying trashbin. (%s)", err.Error())
	}

	// Can I find a way to see if this actually worked?
}

func Test_Query_Restore_Homedir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
		return
	}

	path := filepath.Join(home, "path.txt")
	defer writeTestData(t, path)()
	assertExists(t, path)

	if err := wastebasket.Trash(path); err != nil {
		t.Errorf("Error trashing file. (%s)", err.Error())
	}
	assertNotExists(t, path)

	result, err := wastebasket.Query(wastebasket.QueryOptions{Paths: []string{path}})
	if assert.NoError(t, err) {
		if assert.Len(t, result.Matches[path], 1) {
			if assert.NoError(t, result.Matches[path][0].Restore()) {
				assertExists(t, path)
			}
		}
	}
}

func assertExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); err != nil {
		t.Errorf("path '%s' doesn't exist", path)
	}
}

func assertNotExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("path '%s' exists, but shouldn't", path)
	}
}

func trash(paths ...string) func() error {
	return func() error {
		return wastebasket.Trash(paths...)
	}
}

func writeTestData(t *testing.T, paths ...string) func() {
	t.Helper()

	for _, path := range paths {
		var err error
		if strings.HasSuffix(path, "/") {
			err = os.Mkdir(path, os.ModePerm)
		} else {
			err = os.WriteFile(path, []byte("test"), os.ModePerm)
		}
		if err != nil {
			t.Errorf("Error writing test data. (%s)", err.Error())
		}
	}

	return func() {
		t.Log("Cleaning up test files")
		for _, path := range paths {
			os.RemoveAll(path)
		}

		for _, path := range paths {
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				continue
			}

			if err == nil {
				t.Errorf("error, file hasn't been deleted.")
			} else {
				t.Errorf("unexpected error cleaning up: %s", err)
			}
		}
		t.Log("Done cleaning up test files")
	}
}
