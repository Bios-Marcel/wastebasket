package impl

import (
	"fmt"

	"github.com/Bios-Marcel/wastebasket"
	"github.com/spf13/cobra"
)

var QueryCmd = &cobra.Command{
	Use:   "query files...",
	Short: "query attempts to find the given files in the trash and print information",
	Long:  "query attempts to find the given files in all available trashbins, this works for both files relative to the working directory and absolute paths. Upon match, information is printed line by line.",
	// If used as root cmd, these will be ignored.
	SuggestFor: []string{"lookup"},
	Aliases:    []string{"lookup"},
	// Currently none, as empty just clears every trashbin it can find.
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := wastebasket.Query(args...)
		if err != nil {
			cmd.PrintErrln(err)
		} else {
			for key, value := range result {
				fmt.Println(key, len(value))
			}
		}
	},
}
