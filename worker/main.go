package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// main is entry for worker node. args: [id repo cmd]
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

	if len(args) != 3 {
		fmt.Printf("[WORKER] [ERROR] Wrong args count.\n")
		return
	}

	fmt.Println("READY")
	time.Sleep(20 * time.Second)
	fmt.Println("DONE")
}
