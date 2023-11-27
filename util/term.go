package util

import (
	"golang.org/x/term"
)

func MustGetTermSize(fd uintptr) int {
	width, _, err := term.GetSize(int(fd))

	if err != nil {
		panic(err)
	}

	return width
}
