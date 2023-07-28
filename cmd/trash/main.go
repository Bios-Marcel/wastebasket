package main

import (
	"os"

	"github.com/Bios-Marcel/wastebasket/v2/cmd/impl"
)

func main() {
	if err := impl.TrashCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
