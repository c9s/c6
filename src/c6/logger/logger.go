package logger

import "fmt"
import "os"

func Warn(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
}

func Info(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+msg, args...)
}
