package impl

import (
	"github.com/Bios-Marcel/wastebasket/v2"
	"github.com/spf13/cobra"
)

var TrashCmd = &cobra.Command{
	Use:   "trash files...",
	Short: "trash moves the specified files into the trashbin",
	Long:  "TODO",
	// If used as root cmd, these will be ignored.
	SuggestFor: []string{"delete", "remove", "recycle"},
	Args:       cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := wastebasket.Trash(args...); err != nil {
			cmd.PrintErrln(err)
		}
	},
}
