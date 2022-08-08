package main

import (
	"fmt"
	"os"
	"time"
)

// main is entry for worker node. args: [id repo cmd]
func main() {
	fmt.Println("#############################")
	fmt.Println("##      dbuild Worker      ##")
	fmt.Println("##   By Mahdi Khancherli   ##")
	fmt.Println("#############################")

	if len(os.Args) != 4 {
		fmt.Printf("[WORKER] [ERROR] Wrong args count (%v, %v).", os.Args, len(os.Args))
		return
	}

	fmt.Println("READY")
	time.Sleep(20 * time.Second)
	fmt.Println("DONE")
}
