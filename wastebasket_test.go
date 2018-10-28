package wastebasket

import (
	"io/ioutil"
	"os"
	"testing"
)

var testFilePath = "test.txt"
var testDirPath = "test-delete-me"
var testFilePathWithSpaces = "te st.txt"
var testDirPathWithSpaces = "test-del ete-me"

//TestTrash tests trashing a single file which is created beforehand
func TestTrashWithExistentFile(t *testing.T) {
	writeTestFile(testFilePath, t)
	error := Trash(testFilePath)
	if error != nil {
		t.Errorf("Error trashing file. (%s)", error)
	}

	_, error = os.Stat(testFilePath)
	if (os.IsNotExist(error)) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted. (%s)", error)
	}

	cleanup()
}

//TestTrash tests trashing a single file which is created beforehand
func TestTrashWithExistentFileWithSpaces(t *testing.T) {
	writeTestFile(testFilePathWithSpaces, t)
	error := Trash(testFilePathWithSpaces)
	if error != nil {
		t.Errorf("Error trashing file. (%s)", error)
	}

	_, error = os.Stat(testFilePathWithSpaces)
	if (os.IsNotExist(error)) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted. (%s)", error)
	}

	cleanup()
}

//TestTrash tests trashing a single file which is created beforehand. the path is of format `./filename`.
func TestTrashWithExistentFileWithSpacesAndDotSlashAppended(t *testing.T) {
	writeTestFile(testFilePathWithSpaces, t)
	error := Trash("./" + testFilePathWithSpaces)
	if error != nil {
		t.Errorf("Error trashing file. (%s)", error)
	}

	_, error = os.Stat(testFilePathWithSpaces)
	if (os.IsNotExist(error)) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted. (%s)", error)
	}

	cleanup()
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
	if (os.IsNotExist(error)) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted. (%s)", error)
	}

	cleanup()
}

//TestTrash tests trashing a single file which is created beforehand
func TestTrashWithExistentFolder(t *testing.T) {
	writeTestDirectory(testDirPath, t)
	error := Trash(testDirPath)
	if error != nil {
		t.Errorf("Error trashing file. (%s)", error)
	}

	_, error = os.Stat(testDirPath)
	if (os.IsNotExist(error)) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted. (%s)", error)
	}

	cleanup()
}

//TestTrash tests trashing a single file which is created beforehand
func TestTrashWithExistentFolderWithSpaces(t *testing.T) {
	writeTestDirectory(testDirPathWithSpaces, t)
	error := Trash(testDirPathWithSpaces)
	if error != nil {
		t.Errorf("Error trashing file. (%s)", error)
	}

	_, error = os.Stat(testDirPathWithSpaces)
	if (os.IsNotExist(error)) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted. (%s)", error)
	}
	cleanup()
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
	if (os.IsNotExist(error)) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted. (%s)", error)
	}

	cleanup()
}

//TestEmpty tests emptying the systems trashbin
func TestEmpty(t *testing.T) {
	error := Empty()
	if error != nil {
		t.Errorf("Error emptying trashbin. (%s)", error)
	}
}

func writeTestFile(path string, t *testing.T) {
	writeError := ioutil.WriteFile(path, []byte("Test"), os.ModePerm)
	if writeError != nil {
		t.Errorf("Error writing testfile. (%s)", writeError)
	}
}

func writeTestDirectory(path string, t *testing.T) {
	mkdirError := os.Mkdir(testDirPath, os.ModePerm)
	if mkdirError != nil {
		t.Errorf("Error creating test directory. (%s)", mkdirError)
	}
}

func cleanup() {
	os.Remove(testDirPath)
	os.Remove(testFilePath)
}
