package client

import (
	"errors"
	"fmt"
)

var ErrRequired = errors.New("required condition failed")

func require(value bool, message string) { // TODO - Is it better to build strings and not use or pass a lazy function?
	if !value {
		panic(fmt.Errorf("%s: %w", message, ErrRequired))
	}
}
