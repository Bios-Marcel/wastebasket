//go:build !windows && !darwin

package wastebasket

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

func create_files(count int) []string {
	// Prevent slowdowns of consecutive runs, since we have to check for file
	// existence more often if we create the file multiple times.
	timeNow := time.Now().Format(time.RFC3339)
	files := make([]string, 0, count)
	for i := 1; i <= count; i++ {
		path := fmt.Sprintf("./%s-%d.txt", timeNow, i)
		files = append(files, path)
		err := os.WriteFile(path, []byte("test"), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	return files
}

const manyFilesCount = 20

func Benchmark_gio_trash_manyFiles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		files := create_files(manyFilesCount)
		b.StartTimer()

		if err := gioTrash(files...); err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_customImpl_trash_manyFiles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		files := create_files(manyFilesCount)
		b.StartTimer()

		if err := Trash(files...); err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_gio_trash_singleFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		files := create_files(1)
		b.StartTimer()

		if err := gioTrash(files...); err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_customImpl_trash_singleFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		files := create_files(1)
		b.StartTimer()

		if err := Trash(files...); err != nil {
			b.Error(err)
		}
	}
}

var (
	availabilityCache   = make(map[string]bool)
	errToolNotAvailable = errors.New("tool not available")
)

func isCommandAvailable(name string) bool {
	avail, ok := availabilityCache[name]
	if avail && ok {
		return true
	}
	_, fileError := exec.LookPath(name)
	availabilityCache[name] = fileError == nil
	return fileError == nil
}

func gioTrash(paths ...string) error {
	if isCommandAvailable("gio") {
		// --force makes sure we don't get errors for non-existent files.
		parameters := append([]string{"trash", "--force"}, paths...)
		return exec.Command("gio", parameters...).Run()
	}

	return errToolNotAvailable
}

func trashCli(paths ...string) error {
	if isCommandAvailable("trash") {
		//trash-cli throws 74 in case the file doesn't exist
		existingFiles := make([]string, 0, len(paths))
		for _, path := range paths {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				continue
			}

			existingFiles = append(existingFiles, path)
		}

		parameters := append([]string{"--"}, existingFiles...)
		return exec.Command("trash", parameters...).Run()
	}

	return errToolNotAvailable
}

func gioEmpty() error {
	if isCommandAvailable("gio") {
		return exec.Command("gio", "trash", "--empty").Run()
	}

	return errToolNotAvailable
}

func trashCliEmpty() error {
	if isCommandAvailable("trash-empty") {
		return exec.Command("trash-empty").Run()
	}

	return errToolNotAvailable
}
