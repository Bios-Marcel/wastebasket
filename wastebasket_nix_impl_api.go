//go:build !windows && !darwin && !nix_wrapper

package wastebasket

func Trash(paths ...string) error {
	return customImplTrash(paths...)
}

func Empty() error {
	return customImplEmpty()
}
