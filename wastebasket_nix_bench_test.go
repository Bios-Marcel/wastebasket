//go:build !windows && !darwin && nix_wrapper_impl

package wastebasket

import (
	"fmt"
	"os"
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

		if err := customImplTrash(files...); err != nil {
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

		if err := customImplTrash(files...); err != nil {
			b.Error(err)
		}
	}
}
