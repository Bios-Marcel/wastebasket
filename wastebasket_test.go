package wastebasket

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	testFilePath                 = "test.txt"
	testDirPath                  = "test-delete-me"
	testFilePathWithSpaces       = "te st.txt"
	testDirPathWithSpaces        = "test-del ete-me"
	testFilePathWithDoubleQuotes = "foo\"bar\".txt"
)

//TestTrash tests trashing a single file which is created beforehand
func TestTrashWithExistentFile(t *testing.T) {
	writeTestFile(testFilePath, t)
	error := Trash(testFilePath)
	if error != nil {
		t.Errorf("Error trashing file. (%s)", error.Error())
	}

	_, error = os.Stat(testFilePath)
	if os.IsNotExist(error) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted.")
	}

	cleanup()
}

//TestTrash tests trashing a single file which is created beforehand
func TestTrashWithExistentFileWithSpaces(t *testing.T) {
	writeTestFile(testFilePathWithSpaces, t)
	error := Trash(testFilePathWithSpaces)
	if error != nil {
		t.Errorf("Error trashing file. (%s)", error.Error())
	}

	_, error = os.Stat(testFilePathWithSpaces)
	if os.IsNotExist(error) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted.")
	}

	cleanup()
}

//TestTrash tests trashing a single file which is created beforehand. the path is of format `./filename`.
func TestTrashWithExistentFileWithSpacesAndDotSlashAppended(t *testing.T) {
	writeTestFile(testFilePathWithSpaces, t)
	error := Trash("./" + testFilePathWithSpaces)
	if error != nil {
		t.Errorf("Error trashing file. (%s)", error.Error())
	}

	_, error = os.Stat(testFilePathWithSpaces)
	if os.IsNotExist(error) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted.")
	}

	cleanup()
}

//TestTrash tests trashing a single file which doesn't exist
func TestTrashWithNonexistentFile(t *testing.T) {
	_, error := os.Stat(testFilePath)
	if !os.IsNotExist(error) {
		t.Errorf("File shouldn't exist at start of this test.")
	}

	error = Trash(testFilePath)
	if error != nil {
		t.Errorf("Error trashing file. (%s)", error)
	}

	_, error = os.Stat(testFilePath)
	if os.IsNotExist(error) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted.")
	}

	cleanup()
}

//TestTrash tests trashing a single file which is created beforehand
func TestTrashWithExistentFolder(t *testing.T) {
	writeTestDirectory(testDirPath, t)
	error := Trash(testDirPath)
	if error != nil {
		t.Errorf("Error trashing file. (%s)", error.Error())
	}

	_, error = os.Stat(testDirPath)
	if os.IsNotExist(error) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted.")
	}

	cleanup()
}

//TestTrash tests trashing a single file which is created beforehand
func TestTrashWithExistentNonEmptyFolder(t *testing.T) {
	writeTestDirectory(testDirPath, t)
	writeTestFile(filepath.Join(testDirPath, testFilePath), t)
	error := Trash(testDirPath)
	if error != nil {
		t.Errorf("Error trashing file. (%s)", error.Error())
	}

	_, error = os.Stat(testDirPath)
	if os.IsNotExist(error) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted.")
	}

	cleanup()
}

//TestTrash tests trashing a single file which is created beforehand
func TestTrashWithExistentFolderWithSpaces(t *testing.T) {
	writeTestDirectory(testDirPathWithSpaces, t)
	error := Trash(testDirPathWithSpaces)
	if error != nil {
		t.Errorf("Error trashing file. (%s)", error.Error())
	}

	_, error = os.Stat(testDirPathWithSpaces)
	if os.IsNotExist(error) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted.")
	}
	cleanup()
}

//TestTrash tests trashing a single file which doesn't exist
func TestTrashWithNonexistentFolder(t *testing.T) {
	_, error := os.Stat(testDirPath)
	if !os.IsNotExist(error) {
		t.Errorf("File shouldn't exist at start of this test.")
	}

	error = Trash(testDirPath)
	if error != nil {
		t.Errorf("Error trashing file. (%s)", error.Error())
	}

	_, error = os.Stat(testDirPath)
	if os.IsNotExist(error) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted.")
	}

	cleanup()
}

// TestTrashWithExistentFileWithDoubleQuotes tests trashing a single file with a double quote in its name
func TestTrashWithExistentFileWithDoubleQuotes(t *testing.T) {
	writeTestFile(testFilePathWithDoubleQuotes, t)
	error := Trash(testFilePathWithDoubleQuotes)
	if error != nil {
		t.Errorf("Error trashing file. (%s)", error.Error())
	}

	_, error = os.Stat(testFilePathWithDoubleQuotes)
	if os.IsNotExist(error) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted.")
	}

	cleanup()
}

//TestEmpty tests emptying the systems trashbin
func TestEmpty(t *testing.T) {
	error := Empty()
	if error != nil {
		t.Errorf("Error emptying trashbin. (%s)", error.Error())
	}

	//Can I found a way to see if this actually worked?
}

func writeTestFile(path string, t *testing.T) {
	writeError := ioutil.WriteFile(path, []byte("Test"), os.ModePerm)
	if writeError != nil {
		t.Errorf("Error writing test file. (%s)", writeError.Error())
	}
}

func writeTestDirectory(path string, t *testing.T) {
	mkdirError := os.Mkdir(testDirPath, os.ModePerm)
	if mkdirError != nil {
		t.Errorf("Error creating test directory. (%s)", mkdirError.Error())
	}
}

func cleanup() {
	os.Remove(testDirPath)
	os.Remove(testDirPathWithSpaces)
	os.Remove(testFilePath)
	os.Remove(testFilePathWithSpaces)
	os.Remove(testFilePathWithDoubleQuotes)
}
