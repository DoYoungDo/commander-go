package commandergo

import (
	"fmt"
	"testing"
)

func TestNewCommand(t *testing.T) {
	cmd, err := newCommandWithNameAndArg("cmd [arg...]")
	fmt.Println(cmd, err)
}
