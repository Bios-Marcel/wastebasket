package wastebasket

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	testFilePath           = "test.txt"
	testDirPath            = "test-delete-me"
	testFilePathWithSpaces = "te st.txt"
	testDirPathWithSpaces  = "test-del ete-me"
)

//TestTrash tests trashing a single file which is created beforehand
func TestTrashWithExistentFile(t *testing.T) {
	defer cleanup()

	writeTestFile(testFilePath, t)
	if errTrash := Trash(testFilePath); errTrash != nil {
		t.Errorf("Error trashing file. (%s)", errTrash.Error())
	}

	if _, errStat := os.Stat(testFilePath); os.IsNotExist(errStat) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted.")
	}
}

//TestTrash tests trashing a single file which is created beforehand
func TestTrashWithExistentFileWithSpaces(t *testing.T) {
	defer cleanup()

	writeTestFile(testFilePathWithSpaces, t)
	if errTrash := Trash(testFilePathWithSpaces); errTrash != nil {
		t.Errorf("Error trashing file. (%s)", errTrash.Error())
	}

	if _, errStat := os.Stat(testFilePathWithSpaces); os.IsNotExist(errStat) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted.")
	}
}

//TestTrash tests trashing a single file which is created beforehand. the path is of format `./filename`.
func TestTrashWithExistentFileWithSpacesAndDotSlashAppended(t *testing.T) {
	defer cleanup()

	writeTestFile(testFilePathWithSpaces, t)
	if errTrash := Trash("./" + testFilePathWithSpaces); errTrash != nil {
		t.Errorf("Error trashing file. (%s)", errTrash.Error())
	}

	if _, errStat := os.Stat(testFilePathWithSpaces); os.IsNotExist(errStat) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted.")
	}
}

//TestTrash tests trashing a single file which doesn't exist
func TestTrashWithNonexistentFile(t *testing.T) {
	defer cleanup()

	if _, errStat := os.Stat(testFilePath); !os.IsNotExist(errStat) {
		t.Errorf("File shouldn't exist at start of this test.")
	}

	if errTrash := Trash(testFilePath); errTrash != nil {
		t.Errorf("Error trashing file. (%s)", errTrash)
	}

	if _, errStat := os.Stat(testFilePath); os.IsNotExist(errStat) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted.")
	}
}

//TestTrash tests trashing a single file which is created beforehand
func TestTrashWithExistentFolder(t *testing.T) {
	defer cleanup()

	writeTestDirectory(testDirPath, t)
	if errTrash := Trash(testDirPath); errTrash != nil {
		t.Errorf("Error trashing file. (%s)", errTrash.Error())
	}

	if _, errStat := os.Stat(testDirPath); os.IsNotExist(errStat) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted.")
	}
}

//TestTrash tests trashing a single file which is created beforehand
func TestTrashWithExistentNonEmptyFolder(t *testing.T) {
	defer cleanup()

	writeTestDirectory(testDirPath, t)
	writeTestFile(filepath.Join(testDirPath, testFilePath), t)

	if errTrash := Trash(testDirPath); errTrash != nil {
		t.Errorf("Error trashing file. (%s)", errTrash.Error())
	}

	if _, errStat := os.Stat(testDirPath); os.IsNotExist(errStat) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted.")
	}
}

//TestTrash tests trashing a single file which is created beforehand
func TestTrashWithExistentFolderWithSpaces(t *testing.T) {
	defer cleanup()

	writeTestDirectory(testDirPathWithSpaces, t)
	if errTrash := Trash(testDirPathWithSpaces); errTrash != nil {
		t.Errorf("Error trashing file. (%s)", errTrash.Error())
	}

	if _, errStat := os.Stat(testDirPathWithSpaces); os.IsNotExist(errStat) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted.")
	}
}

//TestTrash tests trashing a single file which doesn't exist
func TestTrashWithNonexistentFolder(t *testing.T) {
	defer cleanup()

	if _, errStat := os.Stat(testDirPath); !os.IsNotExist(errStat) {
		t.Errorf("File shouldn't exist at start of this test.")
	}

	if errTrash := Trash(testDirPath); errTrash != nil {
		t.Errorf("Error trashing file. (%s)", errTrash.Error())
	}

	if _, errStat := os.Stat(testDirPath); os.IsNotExist(errStat) {
		//Everything correct!
	} else {
		t.Errorf("File hasn't been deleted.")
	}
}

//TestEmpty tests emptying the systems trashbin
func TestEmpty(t *testing.T) {
	error := Empty()
	if error != nil {
		t.Errorf("Error emptying trashbin. (%s)", error.Error())
	}

	//Can I find a way to see if this actually worked?
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
