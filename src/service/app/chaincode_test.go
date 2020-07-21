package app

import (
	"fmt"
	"testing"
)

func TestFindUninstallChaincode(t *testing.T) {
	dirs, err := findAllChaincodeDir()
	fmt.Println(dirs, err)
}
