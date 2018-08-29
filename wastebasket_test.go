package wastebasket

import (
	"io/ioutil"
	"os"
	"testing"
)

var testFilePath = "test.txt"

//TestTrash tests trashing a single file which is created beforehand
func TestTrash(t *testing.T) {
	error := Trash(testFilePath)
	if error != nil {
		t.Errorf("Error trashing file. (%s)", error)
	}

	_, error = os.Stat(testFilePath)
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

func writeTestFile() {
	ioutil.WriteFile(testFilePath, []byte("Test"), os.ModePerm)
}
