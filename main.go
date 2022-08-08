package main

import (
	"fmt"
	"github.com/mmkhmmkh/dbuild/pkg/hamctl"
)

// main is entry for orchestrator node
func main() {
	err := hamctl.CreateApp("dbuild-worker-1", "mmkhmmkh/dbuild", "master")
	if err != nil {
		fmt.Println(err)
	}
}
