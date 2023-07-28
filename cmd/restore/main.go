package main

import (
	"os"

	"github.com/Bios-Marcel/wastebasket/v2/cmd/impl"
)

func main() {
	if err := impl.RestoreCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
