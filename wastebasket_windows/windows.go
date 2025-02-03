package wastebasket_windows

import (
	"hash/fnv"
	"time"
)

type TrashedFileInfo struct {
	fileSize     uint64
	infoPath     string
	originalPath string
	deletionDate time.Time
	restoreFunc  func(force bool) error
	deleteFunc   func() error
}

func NewTrashedFileInfo(
	fileSize uint64,
	infoPath, originalPath string, deletionDate time.Time,
	restore func(force bool) error,
	deleteFunc func() error,
) *TrashedFileInfo {
	return &TrashedFileInfo{
		fileSize:     fileSize,
		infoPath:     infoPath,
		originalPath: originalPath,
		deletionDate: deletionDate,
		restoreFunc:  restore,
		deleteFunc:   deleteFunc,
	}
}

func (t TrashedFileInfo) FileSize() uint64 {
	return t.fileSize
}

func (t TrashedFileInfo) UniqueIdentifier() string {
	hash := fnv.New64()
	return string(hash.Sum([]byte(t.infoPath)))
}

func (t TrashedFileInfo) OriginalPath() string {
	return t.originalPath
}

func (t TrashedFileInfo) DeletionDate() time.Time {
	return t.deletionDate
}

func (t TrashedFileInfo) Restore(force bool) error {
	return t.restoreFunc(force)
}

func (t TrashedFileInfo) Delete() error {
	return t.deleteFunc()
}
