package logger

// import "fmt"
import "log"

func Warn(msg string, args ...interface{}) {
	log.Fatalf(msg, args...)
}

func Info(msg string, args ...interface{}) {
	log.Printf(msg, args...)
}
