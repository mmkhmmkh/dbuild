package main

import (
	"fmt"
	"github.com/mmkhmmkh/dbuild/pkg/hamctl"
)

const (
	dbuildRepo        = "mmkhmmkh/dbuild"
	controllerContext = "controller"
	workerContext     = "worker"
	dbuildBranch      = "master"
	dbuildPrefix      = "dbuild-"
)

func StartController() error {
	fmt.Printf("[ORCH] Starting controller...\n")
	err := hamctl.CreateApp(dbuildPrefix+controllerContext, dbuildRepo, dbuildBranch, controllerContext)
	if err != nil {
		return err
	}

	fmt.Printf("[ORCH] Controller started successfuly.\n")

	return nil
}

// main is entry for orchestrator node
func main() {

	fmt.Println("#########################")
	fmt.Println("## dbuild Orchestrator ##")
	fmt.Println("## By Mahdi Khancherli ##")
	fmt.Println("#########################")

	err := StartController()
	if err != nil {
		fmt.Printf("[ORCH] [ERROR] %v\n", err)
		return
	}
}
