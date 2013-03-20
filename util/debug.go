package util

import (
	"fmt"
	"os"
	"runtime"
)

func fileformat(file string) (filename string) {
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			filename = file[i+1:]
			break
		}
	}
	return filename
}

func Debug(format string, v ...interface{}) {
	if os.Getenv("DEBUG") == "" {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	file = fileformat(file)
	fmt.Printf("%s %v: %s\n", file, line, fmt.Sprintf(format, v...))
}
