package errors

import "github.com/rkrmr33/pkg/log"

// Must calls log.Fatal in case err is not nil
func Must(err error) {
	if err != nil {
		log.G().Fatalw("fatal error", "err", err)
	}
}

// MustV like Must but returns the value
func MustV[T any](val T, err error) T {
	Must(err)
	return val
}

// Drop drops the first return falue and returns the second one
func Drop[F any, T any](f F, last T) T {
	return last
}

// Drop2 drops the first two return values and returns the last one
func Drop2[F any, S any, T any](f F, s S, last T) T {
	return last
}
