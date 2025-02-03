package impl

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/Bios-Marcel/wastebasket/v2"
	"github.com/spf13/cobra"
)

var RestoreCmd = &cobra.Command{
	Use:   "restore files...",
	Short: "Restores a specific file (or many files via a glob) to the original location",
	Example: `
  wastebasket restore /home/user/document/file.pdf
  # ID for the case that there are multiple versions of the given file.
  wastebasket restore /home/user/document/file.pdf@5033e67d2c9d
`,
	// If used as root cmd, these will be ignored.
	SuggestFor: []string{"recover"},
	Aliases:    []string{"recover"},
	// Currently none, as empty just clears every trashbin it can find.
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		options := wastebasket.QueryOptions{}

		var err error
		options.Glob, err = cmd.Flags().GetBool("glob")
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		var id string
		arg := args[0]
		indexOfAt := strings.LastIndexByte(arg, '@')
		// 17 = identifiersize - 1.
		if indexOfAt != -1 && indexOfAt == len(arg)-17 {
			id = string(arg[indexOfAt+1:])
			options.Search = []string{arg[:indexOfAt]}
		} else {
			options.Search = args
		}
		arg = options.Search[0]

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		result, err := wastebasket.Query(options)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		matches := result.Matches[arg]
		if id != "" {
			matches = slices.DeleteFunc(matches,
				func(item wastebasket.TrashedFileInfo) bool {
					return item.UniqueIdentifier() != id
				})
		}

		if len(matches) == 0 {
			cmd.PrintErrf("No matching file found for '%s'\n", arg)
			return
		}

		if len(matches) == 1 {
			cmd.Printf("Restoring '%s' ...\n", matches[0].OriginalPath())
			if err := matches[0].Restore(force); err != nil {
				cmd.PrintErrf("error restoring '%s':\n\t%s\n", arg, err)
				os.Exit(1)
			}
			return
		}

		dedupe := make(map[string][]wastebasket.TrashedFileInfo)
		for _, match := range matches {
			dedupe[match.OriginalPath()] = append(dedupe[match.OriginalPath()], match)
		}

		for _, arr := range dedupe {
			if len(arr) == 1 {
				match := arr[0]
				fmt.Printf("Restoring '%s' from '%s'\n",
					match.OriginalPath(), match.DeletionDate())
				if err := match.Restore(force); err != nil {
					cmd.PrintErrf("Error restoring '%s': %s\n", match.OriginalPath(), err)
					os.Exit(1)
				}
			}
		}
		for _, arr := range dedupe {
			if len(arr) > 1 {
				fmt.Printf("Not restoring '%s'; multiple matches:\n",
					arr[0].OriginalPath())
				for _, match := range arr {
					fmt.Printf("\t'%s' %s (ID %s)\n", match.OriginalPath(), match.DeletionDate(), match.UniqueIdentifier())
				}
			}
		}
	},
}

func init() {
	RestoreCmd.Flags().Bool("glob", false, "If set, the given paths will be treated as globs instead of normal paths.")
	RestoreCmd.Flags().Bool("force", false, "If set, restore will overwrite existing files.")
}
