//go:build !windows && !darwin

package wastebasket

import (
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"

	copyLib "github.com/otiai10/copy"
	"golang.org/x/sys/unix"
)

var (
	// cachedInformation makes sure we don't constantly check for the
	// directory and which drive it is on again.
	cachedInformation = &cache{
		init: &sync.Once{},
	}
)

type cache struct {
	// device (or partition) is required in order to decide which trashbin to
	// use, as our implementation also supports non-home trashbins.
	device uint64
	path   string

	init *sync.Once
	err  error
}

func getCache() (*cache, error) {
	cachedInformation.init.Do(func() {
		var homeTrashDir string
		// On some big distros, such as Ubuntu for example, this variable isn't
		// set. Instead, we will fallback to what Ubuntu does for now.
		if dataHome := os.Getenv("XDG_DATA_HOME"); dataHome == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				cachedInformation.err = err
				return
			}

			homeTrashDir = filepath.Join(homeDir, ".local", "share", "Trash")
		} else {
			homeTrashDir = filepath.Join(dataHome, "Trash")
		}

		// Since the actual trash directory might not yet exist, we try to go up
		// the hierarchy in order to find out which device / partition the trash
		// will be on.
		stat := unix.Stat_t{}
		for closestExistingParent := filepath.Dir(homeTrashDir); ; {
			if err := unix.Lstat(closestExistingParent, &stat); err != nil {
				if errno, ok := err.(unix.Errno); !ok || errno != unix.ENOENT {
					cachedInformation.err = err
					return
				}
			} else if err == nil {
				break
			}

			closestExistingParent = filepath.Dir(closestExistingParent)
		}

		cachedInformation.device = stat.Dev
		cachedInformation.path = homeTrashDir
		cachedInformation.err = nil
	})

	return cachedInformation, cachedInformation.err
}

// fileExists omits the parts to make this usable cross-platform and
// therefore saves a minimal amount of CPU cycles and some allocations.
func fileExists(path string) (bool, error) {
	var (
		stat unix.Stat_t
		err  error
	)

RETRY:
	for {
		err = unix.Stat(path, &stat)
		switch err {
		case nil:
			// No issue, exists
			return true, nil
		case unix.EINTR:
			// Doesn't exist
			continue RETRY
		case unix.ENOENT:
			// A non-error basically which tells you to try again.
			return false, nil
		default:
			// Unexpected error
			return false, err
		}
	}
}

// escapeUrl escapes the path according to the FreeDesktop Trash specification.
// Which basically just refers to "RFC 2396, section 2".
func escapeUrl(path string) string {
	u := &url.URL{Path: path}
	return u.EscapedPath()

}

// getTopDir returns the high folder til it reaches the top or finds the
// mountpoint of a parent directory.
func getTopDir(stat *unix.Stat_t, path string) (string, error) {
	parentDirStat := unix.Stat_t{}
	// If path isn't a directory, then parentDir should be the mounted
	// directory or other directory, but not the mounted directories parent
	// and we shouldn't return at that point, therefore we assume this to be
	// a safe way of doing things.
	lastParentDir := path
	parentDir := filepath.Dir(path)
	for {
		if parentDir == path || parentDir == "/" {
			return parentDir, nil
		}

		if err := unix.Lstat(parentDir, &parentDirStat); err != nil {
			return "", err
		}

		if parentDirStat.Dev != stat.Dev {
			return lastParentDir, nil
		}

		lastParentDir = parentDir
		parentDir = filepath.Dir(parentDir)
	}
}

func customImplTrash(paths ...string) error {
	// FIXME Allow absence of sticky b it via option, if not supported) by FS
	// FIXME Check for permissions and set the correctly
	// FIXME Copy metadata when moving across file systems
	// FIXME Check which types of files we aren't allowed to delete.
	// FIXME Query all mounts instead, as this will require only one file read
	// instead of multiple Lstat calls, resulting in less time consumed.
	// FIXME Decide whether we early exit on errors or try to delete all paths.
	// Later, this can be a setting. The decision should be documented.

	// RFC3339 defined in the time package contains the timezone offset, which
	// isn't defined by the spec and causes issues in some trash tools, such
	// as trash-cli.
	deletionDate := time.Now().Format("2006-01-02T15:04:05")
	cache, err := getCache()
	if err != nil {
		return fmt.Errorf("error determining user trash directory: %w", err)
	}

	var (
		stat   unix.Stat_t
		statfs unix.Statfs_t
	)

	for _, absPath := range paths {
		var err error
		absPath, err = filepath.Abs(absPath)
		if err != nil {
			return err
		}

		if err := unix.Lstat(absPath, &stat); err != nil {
			// File doesn't exist, there is no work to actually do, as there
			// is nothing to delete. Additionally, the following code won't
			// treat file absence without errors.
			if err == unix.ENOENT {
				continue
			}

			return err
		}

		// We only support absolute filenames in the home trash. For
		// topdirs, we use relative paths. This allows us to move a
		// mount, while still keeping trash files recoverable.
		var pathForTrashInfo string

		var trashDir, filesDir, infoDir string
		// Deleting accross partitions / mounts
		if cache.device != stat.Dev {
			topDir, err := getTopDir(&stat, absPath)
			if err != nil {
				return err
			}

			// While getTopDir won't return an empty string with its current
			// impl, this can change in the future, so beteter be safe than
			// sorry.
			if topDir != "" {
				var uid string
				if currentUser, err := user.Current(); err != nil {
					return err
				} else {
					uid = currentUser.Uid
				}

				trashDir = filepath.Join(topDir, ".Trash")

				var useFallbackTopdirTrash bool
				if trashDirStat, err := os.Stat(trashDir); err != nil {
					if !os.IsNotExist(err) {
						return err
					}
					useFallbackTopdirTrash = true
				} else {
					if trashDirStat.Mode()&fs.ModeSticky != 0 {
						// If the topdir trash directory contains a trash for all
						// users, it needs to have the sticky bit set. This is only
						// required for .Trash though, not for .Trash-$uid.
						useFallbackTopdirTrash = true
					} else if trashDirStat.Mode()&os.ModeSymlink != 0 {
						// Symlinks must not be used as per spec.
						useFallbackTopdirTrash = true
					}
				}

				pathForTrashInfo, err = filepath.Rel(topDir, absPath)
				if err != nil {
					return err
				}

				if !useFallbackTopdirTrash {
					filesDir = filepath.Join(trashDir, uid, "files")
					infoDir = filepath.Join(trashDir, uid, "info")
				} else {
					// If .Trash doesn't exist, we need to check for .Trash-$uid
					// and create it if it doesn't exist. The spec however
					// doesn't indicate that we should do the same with .Trash.
					trashDir = filepath.Join(topDir, ".Trash-"+uid)
					filesDir = filepath.Join(trashDir, "files")
					infoDir = filepath.Join(trashDir, "info")
				}

			}
		}

		if trashDir == "" {
			// Fallback to home trash.
			trashDir = cache.path
			filesDir = filepath.Join(trashDir, "files")
			infoDir = filepath.Join(trashDir, "info")
			// Home trash only supports absolute paths.
			pathForTrashInfo = absPath
		}

		if err := os.MkdirAll(filesDir, 0700); err != nil && !os.IsExist(err) {
			return fmt.Errorf("error creating directory '%s': %w", filesDir, err)
		}
		if err := os.MkdirAll(infoDir, 0700); err != nil && !os.IsExist(err) {
			return fmt.Errorf("error creating directory '%s': %w", infoDir, err)
		}

		baseName := filepath.Base(absPath)
		trashedFilePath := filepath.Join(filesDir, baseName)
		trashedFileInfoPath := filepath.Join(infoDir, baseName) + ".trashinfo"

		// We need to check whether the trash already contains a file with this
		// name, since deleted files from different directories often have the
		// same name. An example would be .gitignore files, they always have
		// the same basename and therefore always the same trash path.
		// We simply count up in this case. Since we've got the info file, we
		// can map back to the original name later on.

		var infoFileHandle *os.File
		if exists, err := fileExists(trashedFilePath); err != nil {
			return err
		} else if !exists {
			// We save ourselves the fileExists check, as we can combine it
			// with the opening of the file handle. This is a performance
			// optimisation.
			infoFileHandle, err = os.OpenFile(trashedFileInfoPath, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0600)
			if err != nil {
				if !os.IsExist(err) {
					return err
				}
			}

			// While we close manually later, we want to prevent a leak.
			defer func() {
				infoFileHandle.Close()
			}()
		}

		// If there isn't a valid info file handle yet, it means that one
		// of the two file names were already in use, requiring us to find
		// two unique filenames eiter way.
		if infoFileHandle == nil {
			extension := filepath.Ext(baseName)
			baseNameNoExtension := strings.TrimSuffix(baseName, extension)
			for i := uint64(1); i != 0; i = i + 1 {
				newBaseName := fmt.Sprintf("%s.%d%s", baseNameNoExtension, i, extension)

				// The names of both files must always be the same, putting
				// aside the .trashinfo extension.
				trashedFilePath = filepath.Join(filesDir, newBaseName)
				if exists, err := fileExists(trashedFilePath); err != nil || exists {
					continue
				}
				infoFileHandle, err = os.OpenFile(filepath.Join(infoDir, newBaseName+".trashinfo"), os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0600)
				if err != nil {
					if os.IsExist(err) {
						continue
					}
					return err
				}

				// We found a valid name, where neither the file itself, nor
				// the trashinfo file exist.
				break
			}
		}

		if err := os.Rename(absPath, trashedFilePath); err != nil {
			// Moving across different filesystems causes os.Rename to fail.
			// Therefore we need to do a costly copy followed by a remove.
			if linkErr, ok := err.(*os.LinkError); ok && linkErr.Err.Error() == "invalid cross-device link" {
				if err := unix.Statfs(trashDir, &statfs); err == nil {
					trashDirFsType := statfs.Type
					if err := unix.Statfs(absPath, &statfs); err == nil {
						if trashDirFsType != statfs.Type {
							if err := copyLib.Copy(absPath, trashedFilePath); err != nil {
								return fmt.Errorf("error copying files into trash directory: %w", err)
							}

							if err := os.RemoveAll(absPath); err != nil {
								return fmt.Errorf("error removing files (a copy into the trash has been made successfully): %w", err)
							}

							// Success of special treatment, proceed normally
							goto WRITE_FILE_INFO
						}
					}
				}
			}

			// All special treatment failed, return original os.Rename error
			return err
		}

	WRITE_FILE_INFO:
		if _, err = infoFileHandle.WriteString(fmt.Sprintf("[Trash Info]\nPath=%s\nDeletionDate=%s\n", escapeUrl(pathForTrashInfo), deletionDate)); err != nil {
			return err
		}
	}

	return nil
}

func customImplEmpty() error {
	// FIXME Figure out, whether this should only empty whatever the spec would
	// also demanding deleting into or all reachable trashbins. An alternative
	// would be to clear the topdir trash, if available and the user trash.
	// Considering the nature of wastebasket, it would probably be best to clear
	// the user trash.
	//
	// In the future, we could optionally allow passing a path here, so the
	// user can define a custom path or clearing options.
	//
	// This could have a format where you can define different options for
	// different platforms, such as:
	//   wastebasket.Empty(
	//	      wastebasket.Pattern("*.txt"),
	//        nix.ClearUserTrashbin(),
	//        darwin.ClearAllAvailableTrashbins(),
	//   )
	return nil
}
