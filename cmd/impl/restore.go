package impl

import (
	"github.com/Bios-Marcel/wastebasket/v2"
	"github.com/spf13/cobra"
)

var RestoreCmd = &cobra.Command{
	Use:   "restore files...",
	Short: "TODO",
	Long:  "TODO",
	// If used as root cmd, these will be ignored.
	SuggestFor: []string{"recover"},
	Aliases:    []string{"recover"},
	// Currently none, as empty just clears every trashbin it can find.
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		options := wastebasket.QueryOptions{}

		glob, err := cmd.Flags().GetBool("glob")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		if glob {
			options.Globs = args
		} else {
			options.Paths = args
		}

		failfast, err := cmd.Flags().GetBool("failfast")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		options.FailFast = failfast

		result, err := wastebasket.Query(options)
		if err != nil {
			cmd.PrintErrln(err)
		} else {
			for _, arg := range args {
				matches := result.Matches[arg]
				if len(matches) == 0 {
					cmd.PrintErrf("no matching file found for '%s'\n", arg)
				} else if len(matches) > 1 {
					cmd.PrintErrf("multiple matching files found for '%s'\n", arg)
					for _, match := range matches {
						cmd.PrintErrf("\t'%s'\n", match)
					}
				} else {
					if err := matches[0].Restore(); err != nil {
						cmd.PrintErrf("error restoring '%s':\n\t%s\n", arg, err)
					}
				}
			}
		}
	},
}

func init() {
	RestoreCmd.Flags().Bool("glob", false, "If set, the given paths will be treated as globs instead of normal paths.")
	RestoreCmd.Flags().Bool("failfast", false, "If set, the query will stop upon the first error.")
}
