package wastebasket

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

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
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer writeTestData(t, c.testDataToCreate...)()

			for _, expectation := range c.trashExpectations {
				err := expectation.trasher()
				if err != expectation.expectedErr {
					t.Errorf("unexpected error: %v != %v", err, expectation.expectedErr)
				}
			}

		OUTER_LOOP:
			for _, file := range c.testDataToCreate {
				_, err := os.Stat(file)
				for _, expectedFile := range c.expectedFiles {
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

// TestEmpty tests emptying the systems trashbin
func TestEmpty(t *testing.T) {
	error := Empty()
	if error != nil {
		t.Errorf("Error emptying trashbin. (%s)", error.Error())
	}

	//Can I find a way to see if this actually worked?
}

func trash(paths ...string) func() error {
	return func() error {
		return Trash(paths...)
	}
}

func writeTestData(t *testing.T, paths ...string) func() {
	for _, path := range paths {
		var err error
		if strings.HasSuffix(path, "/") {
			err = os.Mkdir(path, os.ModePerm)
		} else {
			err = ioutil.WriteFile(path, []byte("test"), os.ModePerm)
		}
		if err != nil {
			t.Errorf("Error writing test data. (%s)", err.Error())
		}
	}

	return func() {
		for _, path := range paths {
			os.RemoveAll(path)
		}

		for _, path := range paths {
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				return
			}

			if err == nil {
				t.Errorf("error, file hasn't been deleted.")
			} else {
				t.Errorf("unexpected error cleaning up: %s", err)
			}
		}
	}
}
