package tmpl

import (
	"fmt"
	"testing"
)

func Test_ReadFile(t *testing.T) {

	b, err := GetFile("Makefile")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", b)
}
