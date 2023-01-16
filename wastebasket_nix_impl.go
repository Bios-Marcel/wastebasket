//go:build !windows && !darwin && !nix_wrapper

package wastebasket

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	copyLib "github.com/otiai10/copy"
	"golang.org/x/sys/unix"
)

func determineHomeTrashDir() (string, error) {
	// On some big distros, such as Ubuntu for example, this variable isn't
	// set. Instead, we will fallback to what Ubuntu does for now.
	if dataHome := os.Getenv("XDG_DATA_HOME"); dataHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		return filepath.Join(homeDir, ".local", "share", "Trash"), nil
	} else {
		return filepath.Join(dataHome, "Trash"), nil
	}
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

func customImplTrash(paths ...string) error {
	// FIXME Move logic that is uselessly repeated into init() function
	// or cache where it makes sense.

	trashDir, err := determineHomeTrashDir()
	if err != nil {
		return fmt.Errorf("error determining user trash directory: %w", err)
	}

	// FIXME Move into function and allow specifying permissions depending on
	// where the directory is located.

	// Assuming that the parent directories should already exist, we don't
	// invoke MkdirAll here. Since we only support user trash for now, we'll
	// accordingly set permissions only for our current user.
	if err := os.Mkdir(trashDir, 0700); err != nil && !os.IsExist(err) {
		return fmt.Errorf("error creating user trash directory: %w", err)
	}
	if err := os.Mkdir(filepath.Join(trashDir, "files"), 0700); err != nil && !os.IsExist(err) {
		return fmt.Errorf("error creating 'files' subdirectory for user trash directory: %w", err)
	}
	if err := os.Mkdir(filepath.Join(trashDir, "info"), 0700); err != nil && !os.IsExist(err) {
		return fmt.Errorf("error creating 'info' subdirectory for user trash directory: %w", err)
	}

	for _, path := range paths {
		// Avoid running into weird errors and there isn't anything to do
		// either way.
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		baseName := filepath.Base(path)

		// FIXME Once we move across partitions, we need to check for
		// permissions beforehand, as it might make the file unreachable otherwise.

		targetPath := filepath.Join(trashDir, baseName)

		// We need to check whether the trash already contains a file with this
		// name, since deleted files from different directories often have the
		// same name. An example would be .gitignore files, they always have
		// the same basename and therefore always the same trash path.
		// We simply count up in this case. Since we've got the info file, we
		// can map back to the original name later on.
		if exists, err := fileExists(targetPath); err != nil {
			return err
		} else if exists {
			extension := filepath.Ext(baseName)
			baseNameNoExtension := strings.TrimSuffix(baseName, extension)
			for i := uint64(1); i != 0; i = i + 1 {
				targetPath = filepath.Join(trashDir, fmt.Sprintf("%s.%d%s", baseNameNoExtension, i, extension))

				if exists, err := fileExists(targetPath); err == nil && !exists {
					// We found a valid name.
					break
				}
			}
		}

		if err := os.Rename(path, targetPath); err != nil {
			// Moving across different filesystems causes os.Rename to fail.
			// Therefore we need to do a costly copy followed by a remove.
			if linkErr, ok := err.(*os.LinkError); ok && linkErr.Err.Error() == "invalid cross-device link" {
				var fsStats unix.Statfs_t
				if err := unix.Statfs(trashDir, &fsStats); err == nil {
					trashDirFsType := fsStats.Type
					if err := unix.Statfs(path, &fsStats); err == nil {
						if trashDirFsType != fsStats.Type {
							if err := copyLib.Copy(path, targetPath); err != nil {
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
	}

WRITE_FILE_INFO:
	// FIXME Write file info

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
