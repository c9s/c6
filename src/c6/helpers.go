package c6

import (
	"fmt"
	"runtime"
)

func unimplemented(featureName string) {
	_, file, line, _ := runtime.Caller(1)
	panic(fmt.Errorf("%s is unimplemented at %s line %d", featureName, file, line))
}
