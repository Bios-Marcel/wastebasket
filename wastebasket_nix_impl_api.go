//go:build !windows && !darwin

package wastebasket

func Trash(paths ...string) error {
	return customImplTrash(paths...)
}

func Empty() error {
	return customImplEmpty()
}
