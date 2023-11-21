package wastebasket

import (
	"errors"
	"time"
)

// TrashedFileInfo represents a file that has been deleted and now resides
// in the trashbin.
type TrashedFileInfo interface {
	// OriginalPath is the files path before it was deleted.
	OriginalPath() string
	// DeletionDate is the deletion date in the computers local timezone.
	DeletionDate() time.Time
	// Restore will attempt restoring the file to its previous location.
	Restore() error
	// Delete will permanently deleting the underlying file. Note that we do not
	// zero the respective bytes on the disk.
	Delete() error
}

type QueryResult struct {
	// Matches are the query results, mapped from the query input (globa / path)
	// to the respective trashed files. Note that the same file can be trashed
	// multiple times and a glob can match multiple files. Therefore, expect
	// multiple entries in both scenarios.
	Matches map[string][]TrashedFileInfo
	// Failures are non-fatal errors that occured during the query. This can
	// for example be a failure of reading a certain file from the trashbin,
	// which however won't prevent us from reading other files. This field is
	// only populated if [QueryOptions.FailFast] is set to `false`.
	Failures []error
}

// QueryOptions allows to configure the Query-Call. Options Globs and Paths
// aren't allowed to both be set.
type QueryOptions struct {
	Globs []string
	// Paths can be relative or absolute.
	Paths []string

	// FailFast will insantly cause the Query call to return if an error occurs.
	// This goes for both errors that occur working on a specifc trashed file
	// and errors that occur for preparing the querying.
	FailFast bool
}

var ErrOnlyOneOfGlobsOrPaths = errors.New("only one of the options .Globs or .Paths must be set.")

func (options QueryOptions) validate() error {
	if len(options.Globs) > 0 && len(options.Paths) > 0 {
		return ErrOnlyOneOfGlobsOrPaths
	}
	return nil
}
