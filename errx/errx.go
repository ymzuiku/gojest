package errx

import (
	"errors"
)

func UnwrapAll(err error) error {
	last := err
	for i := 0; i < 99; i++ {
		e := errors.Unwrap(last)
		if e == nil {
			return last
		}
		last = e
	}
	return err
}
