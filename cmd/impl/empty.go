package impl

import (
	"github.com/Bios-Marcel/wastebasket/v2"
	"github.com/spf13/cobra"
)

var EmptyCmd = &cobra.Command{
	Use:   "empty",
	Short: "empty clears all trashbins that can be found",
	// If used as root cmd, these will be ignored.
	SuggestFor: []string{"clear"},
	// Currently none, as empty just clears every trashbin it can find.
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := wastebasket.Empty(); err != nil {
			cmd.PrintErrln(err)
		}
	},
}
