package jarvis

import (
	"fmt"
	"testing"
)

func Test_refConfig(t *testing.T) {

	for _, ref := range []string{
		"master",
		"develop",
		"feat/xxxx",
	} {
		fmt.Println(_refConfig(ref))

	}
}
