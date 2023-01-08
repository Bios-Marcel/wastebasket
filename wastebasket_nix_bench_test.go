//go:build !windows && !darwin && nix_wrapper_impl && !nix_wrapper

package wastebasket

import (
	"fmt"
	"os"
	"testing"
)

func create_files(count int) []string {
	files := make([]string, 0, count)
	for i := 1; i <= count; i++ {
		path := fmt.Sprintf("./%d.txt", i)
		files = append(files, path)
		err := os.WriteFile(path, []byte("test"), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	return files
}

func Benchmark_gio_trash_manyFiles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		files := create_files(100)
		b.StartTimer()

		if err := gioTrash(files...); err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_customImpl_trash_manyFiles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		files := create_files(100)
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
