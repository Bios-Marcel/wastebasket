//go:build !freebsd && !openbsd && !netbsd && !linux && !windows

package wastebasket

func Query(options QueryOptions) (*QueryResult, error) {
	return nil, ErrPlatformNotSupported
}

func Trash(paths ...string) error {
	return ErrPlatformNotSupported
}

func Empty() error {
	return ErrPlatformNotSupported
}
