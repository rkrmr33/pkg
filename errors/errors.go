package errors

import "github.com/rkrmr33/pkg/log"

// Must calls log.Fatal in case err is not nil
func Must(err error) {
	if err != nil {
		log.G().Fatalw("fatal error", "err", err)
	}
}

// MustV like Must but returns a value
func MustV[T any](v T, err error) T {
	Must(err)
	return v
}
