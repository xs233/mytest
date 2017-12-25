package log

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	Debug("Hello World!2")
	Error("EEEEEE")
	ORM.Info("GET DEEEEEE")
	for i := 0; i < 100000; i++ {
		fmt.Print("")
	}
}
