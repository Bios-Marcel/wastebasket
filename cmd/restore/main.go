package main

import (
	"os"

	"github.com/Bios-Marcel/wastebasket/cmd/impl"
)

func main() {
	if err := impl.RestoreCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
