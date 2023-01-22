//go:build !windows && !darwin

package wastebasket

import (
	"fmt"
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
	cachedInformation *cache
	cacheMutex        = &sync.Mutex{}
)

type cache struct {
	// device (or partition) is required in order to decide which trashbin to
	// use, as our implementation also supports non-home trashbins.
	device uint64
	path   string
}

func getCache() (*cache, error) {
	// Prevent locking, since we never null it again.
	if cachedInformation != nil {
		return cachedInformation, nil
	}

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Make sure it hasn't been initialised in the meantime.
	if cachedInformation != nil {
		return cachedInformation, nil
	}

	var homeTrashDir string
	// On some big distros, such as Ubuntu for example, this variable isn't
	// set. Instead, we will fallback to what Ubuntu does for now.
	if dataHome := os.Getenv("XDG_DATA_HOME"); dataHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
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
				return nil, err
			}
		} else if err == nil {
			break
		}

		closestExistingParent = filepath.Dir(closestExistingParent)
	}

	cachedInformation = &cache{
		device: stat.Dev,
		path:   homeTrashDir,
	}
	return cachedInformation, nil
}

func fileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}

		// File doesn't exist
		return false, nil
	}

	return true, nil
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
	// FIXME Check for sticky bit
	// FIXME Check for symbolic links
	// FIXME Check for permissions and set the correctly
	// FIXME Copy metadata when moving across file systems
	// FIXME Check which types of files we aren't allowed to delete.

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

	for _, path := range paths {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}

		if err := unix.Lstat(path, &stat); err != nil {
			// File doesn't exist, there is no work to actually do, as there
			// is nothing to delete. Additionally, the following code won't
			// treat file absence without errors.
			if err == unix.ENOENT {
				continue
			}

			return err
		}

		var trashDir, filesDir, infoDir string
		// Deleting accross partitions / mounts
		if cache.device != stat.Dev {
			topDir, err := getTopDir(&stat, path)
			if err != nil {
				return err
			}

			var uid string
			if currentUser, err := user.Current(); err != nil {
				return err
			} else {
				uid = currentUser.Uid
			}

			trashDir = filepath.Join(topDir, ".Trash")
			if exists, err := fileExists(trashDir); err != nil {
				return err
			} else if !exists {
				// If .Trash doesn't exist, we need to check for .Trash-$uid
				// and create it if it doesn't exist. The spec however
				// doesn't indicate that we should do the same with .Trash.
				trashDir = filepath.Join(topDir, ".Trash-"+uid)
				filesDir = filepath.Join(trashDir, "files")
				infoDir = filepath.Join(trashDir, "info")
			} else {
				filesDir = filepath.Join(trashDir, uid, "files")
				infoDir = filepath.Join(trashDir, uid, "info")
			}
		}

		// Fallback to home trash.
		if trashDir == "" {
			trashDir = cache.path
			filesDir = filepath.Join(trashDir, "files")
			infoDir = filepath.Join(trashDir, "info")
		}

		if err := os.MkdirAll(filesDir, 0700); err != nil && !os.IsExist(err) {
			return fmt.Errorf("error creating directory '%s': %w", filesDir, err)
		}
		if err := os.MkdirAll(infoDir, 0700); err != nil && !os.IsExist(err) {
			return fmt.Errorf("error creating directory '%s': %w", infoDir, err)
		}

		baseName := filepath.Base(path)
		trashedFilePath := filepath.Join(filesDir, baseName)
		trashedFileInfoPath := filepath.Join(infoDir, baseName) + ".trashinfo"

		// We need to check whether the trash already contains a file with this
		// name, since deleted files from different directories often have the
		// same name. An example would be .gitignore files, they always have
		// the same basename and therefore always the same trash path.
		// We simply count up in this case. Since we've got the info file, we
		// can map back to the original name later on.

		var trashedFilePathExists, trashedFileInfoPathExists bool
		if exists, err := fileExists(trashedFilePath); err != nil {
			return err
		} else {
			trashedFilePathExists = exists
		}
		if !trashedFilePathExists {
			if exists, err := fileExists(trashedFileInfoPath); err != nil {
				return err
			} else {
				trashedFileInfoPathExists = exists
			}
		}

		if trashedFilePathExists || trashedFileInfoPathExists {
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
				trashedFileInfoPath = filepath.Join(infoDir, newBaseName+".trashinfo")
				if exists, err := fileExists(trashedFileInfoPath); err != nil || exists {
					continue
				}

				// We found a valid name, where neither the file itself, nor
				// the trashinfo file exist.
				break
			}
		}

		if err := os.Rename(path, trashedFilePath); err != nil {
			// Moving across different filesystems causes os.Rename to fail.
			// Therefore we need to do a costly copy followed by a remove.
			if linkErr, ok := err.(*os.LinkError); ok && linkErr.Err.Error() == "invalid cross-device link" {
				if err := unix.Statfs(trashDir, &statfs); err == nil {
					trashDirFsType := statfs.Type
					if err := unix.Statfs(path, &statfs); err == nil {
						if trashDirFsType != statfs.Type {
							if err := copyLib.Copy(path, trashedFilePath); err != nil {
								return fmt.Errorf("error copying files into trash directory: %w", err)
							}

							if err := os.RemoveAll(path); err != nil {
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
		// FIXME Support relative paths. Absolute paths may only be used
		// in the home trash.
		if err := os.WriteFile(trashedFileInfoPath, []byte(fmt.Sprintf("[Trash Info]\nPath=%s\nDeletionDate=%s\n", escapeUrl(path), deletionDate)), 0600); err != nil {
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
