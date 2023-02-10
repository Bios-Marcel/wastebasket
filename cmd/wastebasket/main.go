package main

import (
	"os"

	"github.com/Bios-Marcel/wastebasket/cmd/impl"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "wastebasket",
		Short:   "wastebasket allows interaction with the system trashbin",
		Long:    "TODO",
		Example: `wastebasket trash file_1.txt file_2.txt`,
	}
	rootCmd.AddCommand(impl.TrashCmd)
	rootCmd.AddCommand(impl.EmptyCmd)
	rootCmd.AddCommand(impl.QueryCmd)
	rootCmd.AddCommand(impl.RestoreCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
