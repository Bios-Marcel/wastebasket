package wastebasket

import "time"

// TrashedFileInfo represents a file that has been deleted and now resides
// in the trashbin.
type TrashedFileInfo interface {
	// OriginalPath is the files path before it was deleted.
	OriginalPath() string
	// DeletionDate is the deletion date in the computers local timezone.
	DeletionDate() time.Time
	// Restore will attempt restoring the file to its previous location.
	Restore() error
}
