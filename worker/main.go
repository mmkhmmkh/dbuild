package main

import (
	"fmt"
	"github.com/mmkhmmkh/dbuild/pkg/hamctl"
	"os"
	"strings"
	"time"
)

// main is entry for worker node. args: [id repo cmd env]
func main() {
	fmt.Println("#############################")
	fmt.Println("##      dbuild Worker      ##")
	fmt.Println("##   By Mahdi Khancherli   ##")
	fmt.Println("#############################")

	if len(os.Args) != 2 {
		fmt.Printf("[WORKER] [ERROR] Wrong args count.\n")
		return
	}

	args := strings.Split(os.Args[1], " ")

	if len(args) != 4 {
		fmt.Printf("[WORKER] [ERROR] Wrong args count.\n")
		return
	}

	env := args[3]
	err := hamctl.Initialize(env)
	if err != nil {
		fmt.Printf("[WORKER] [ERROR] %v\n", err)
		return
	}

	fmt.Println("READY")
	time.Sleep(20 * time.Second)
	fmt.Println("DONE")
}
