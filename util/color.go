package util

import (
	"fmt"
)

type ColorCode string

const (
	Red   ColorCode = "31"
	Green           = "32"
	Yello           = "33"
	Blue            = "34"
	Pink            = "35"
	Cyan            = "36"
)

type Color struct{}

func (c Color) colorize(row string, color ColorCode) string {
	return fmt.Sprintf("%s%s%s", "\033["+color+"m", row, "\033[0m")
}

func (c Color) Red(row string) string {
	return c.colorize(row, Red)
}

func (c Color) Green(row string) string {
	return c.colorize(row, Green)
}

func (c Color) Yello(row string) string {
	return c.colorize(row, Yello)
}

func (c Color) Blue(row string) string {
	return c.colorize(row, Blue)
}

func (c Color) Pink(row string) string {
	return c.colorize(row, Pink)
}

func (c Color) Cyan(row string) string {
	return c.colorize(row, Cyan)
}
