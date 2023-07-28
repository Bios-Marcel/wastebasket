package wastebasket_windows

import "time"

type TrashedFileInfo struct {
	fileSize     uint64
	originalPath string
	deletionDate time.Time
	restore      func() error
}

func NewTrashedFileInfo(
	fileSize uint64,
	originalPath string, deletionDate time.Time,
	// FIXME
	writeProtected bool,
	// FIXME
	infoPath string,
	restore func() error) *TrashedFileInfo {
	return &TrashedFileInfo{
		fileSize:     uint64(fileSize),
		originalPath: originalPath,
		deletionDate: deletionDate,
		restore:      restore,
	}
}

func (t TrashedFileInfo) FileSize() uint64 {
	return t.fileSize
}

func (t TrashedFileInfo) OriginalPath() string {
	return t.originalPath
}

func (t TrashedFileInfo) DeletionDate() time.Time {
	return t.deletionDate
}

func (t TrashedFileInfo) Restore() error {
	return t.restore()
}
