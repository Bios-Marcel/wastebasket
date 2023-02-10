package impl

import (
	"github.com/Bios-Marcel/wastebasket"
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
		result, err := wastebasket.Query(args...)
		if err != nil {
			cmd.PrintErrln(err)
		} else {
			for _, arg := range args {
				matches := result[arg]
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
