package errors

import "github.com/rkrmr33/pkg/log"

// Must calls log.Fatal in case err is not nil
func Must(err error) {
	if err != nil {
		log.G().Fatalw("fatal error", "err", err)
	}
}
