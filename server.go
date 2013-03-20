package speedy

import (
	"code.google.com/p/go.net/spdy"
)

func GetV() int {
	return spdy.Version
}
