package util

import (
	"os"

	"golang.org/x/term"
)

var (
	Stdin = os.Stdin.Fd()
)

func MustGetTermSize() int {
	width, _, err := term.GetSize(int(Stdin))

	if err != nil {
		panic(err)
	}

	return width
}
