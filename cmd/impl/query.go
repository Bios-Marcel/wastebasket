package impl

import (
	"fmt"

	"github.com/Bios-Marcel/wastebasket/v2"
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
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		options := wastebasket.QueryOptions{}

		glob, err := cmd.Flags().GetBool("glob")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		options.Glob = glob
		options.Search = args

		result, err := wastebasket.Query(options)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		for _, value := range result.Matches {
			for _, value := range value {
				fmt.Printf("%s %s\n", value.OriginalPath(), value.DeletionDate())
			}
		}
	},
}

func init() {
	QueryCmd.Flags().Bool("glob", false, "If set, the given paths will be treated as globs instead of normal paths.")
}
