package util

import (
	"golang.org/x/term"
)

func MustGetTermSize() int {
	width, _, err := term.GetSize(0)

	if err != nil {
		panic(err)
	}

	return width
}
