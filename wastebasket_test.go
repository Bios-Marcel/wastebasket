package wastebasket

import (
	"io/ioutil"
	"os"
	"testing"
)

var testFilePath = "test.txt"
var testDirPath = "test-delete-me"

//TestTrash tests trashing a single file which is created beforehand
func TestTrashWithExistentFile(t *testing.T) {
	writeTestFile(t)
	error := Trash(testFilePath)
	if error != nil {
		t.Errorf("Error trashing file. (%s)", error)
	}

	_, error = os.Stat(testFilePath)
	if os.IsExist(error) {
		t.Errorf("File hasn't been deleted. (%s)", error)
	}
}

//TestTrash tests trashing a single file which doesn't exist
func TestTrashWithNonexistentFile(t *testing.T) {
	_, error := os.Stat(testFilePath)
	if os.IsExist(error) {
		t.Errorf("File shouldn'T exist at start of this test. (%s)", error)
	}

	error = Trash(testFilePath)
	if error != nil {
		t.Errorf("Error trashing file. (%s)", error)
	}

	_, error = os.Stat(testFilePath)
	if os.IsExist(error) {
		t.Errorf("File hasn't been deleted. (%s)", error)
	}
}

//TestTrash tests trashing a single file which is created beforehand
func TestTrashWithExistentFolder(t *testing.T) {
	writeTestDirectory(t)
	error := Trash(testDirPath)
	if error != nil {
		t.Errorf("Error trashing file. (%s)", error)
	}

	_, error = os.Stat(testDirPath)
	if os.IsExist(error) {
		t.Errorf("File hasn't been deleted. (%s)", error)
	}
}

//TestTrash tests trashing a single file which doesn't exist
func TestTrashWithNonexistentFolder(t *testing.T) {
	_, error := os.Stat(testDirPath)
	if os.IsExist(error) {
		t.Errorf("File shouldn'T exist at start of this test. (%s)", error)
	}

	error = Trash(testDirPath)
	if error != nil {
		t.Errorf("Error trashing file. (%s)", error)
	}

	_, error = os.Stat(testDirPath)
	if os.IsExist(error) {
		t.Errorf("File hasn't been deleted. (%s)", error)
	}
}

//TestEmpty tests emptying the systems trashbin
func TestEmpty(t *testing.T) {
	error := Empty()
	if error != nil {
		t.Errorf("Error emptying trashbin. (%s)", error)
	}
}

func writeTestFile(t *testing.T) {
	writeError := ioutil.WriteFile(testFilePath, []byte("Test"), os.ModePerm)
	if writeError != nil {
		t.Errorf("Error writing testfile. (%s)", writeError)
	}
}

func writeTestDirectory(t *testing.T) {
	mkdirError := os.Mkdir(testDirPath, os.ModePerm)
	if mkdirError != nil {
		t.Errorf("Error creating test directory. (%s)", mkdirError)
	}
}
