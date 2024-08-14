package pkg

import (
	"fmt"
)

type ColorInt int

const (
	_ ColorInt = 29 + iota + 1
	Red
	Green
	Yellow
	Blue
	Purple
	Cyan
	White
)

func Color(text string, color ColorInt) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, text)
}
