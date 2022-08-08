package main

import (
	"fmt"
	"github.com/mmkhmmkh/dbuild/pkg/hamctl"
	"github.com/mmkhmmkh/dbuild/pkg/utils"
	"os"
	"strconv"
	"time"
)

func StartWorker(controllerID, workerID string, arguments string) error {
	fmt.Printf("[CTRL] Starting new worker...\n")
	err := hamctl.CreateApp(utils.DbuildPrefix+utils.WorkerContext+"-"+controllerID+"-"+workerID, utils.DbuildRepo, utils.DbuildBranch, utils.WorkerContext, "/dbuild/bin/"+utils.WorkerContext, workerID+" "+arguments)
	if err != nil {
		return err
	}

	fmt.Printf("[CTRL] New worker started successfuly.\n")
	return nil
}

// main is entry for controller node. args: [id n repo cmd]
func main() {
	fmt.Println("#############################")
	fmt.Println("##    dbuild Controller    ##")
	fmt.Println("##   By Mahdi Khancherli   ##")
	fmt.Println("#############################")

	if len(os.Args) != 5 {
		fmt.Printf("[CTRL] [ERROR] Wrong args count (%v, %v).", os.Args, len(os.Args))
		return
	}

	id := os.Args[1]
	n, _ := strconv.Atoi(os.Args[2])
	repo := os.Args[3]
	command := os.Args[4]

	fmt.Printf("[CTRL] Running %v workers...\n", n)
	for i := 1; i < n; i++ {
		err := StartWorker(id, strconv.Itoa(i), fmt.Sprintf("%s %s", repo, command))
		if err != nil {
			fmt.Printf("[CTRL] [WORKER] [ERROR] %v\n", err)
			return
		}
	}

	fmt.Println("READY")
	time.Sleep(60 * time.Second)
	fmt.Println("DONE")

	fmt.Printf("[CTRL] Removing %v workers...\n", n)
	for i := 1; i < n; i++ {
		err := hamctl.RemoveApp(utils.DbuildPrefix + utils.WorkerContext + "-" + id + "-" + strconv.Itoa(i))
		if err != nil {
			fmt.Println("[CTRL] [WORKER] [ERROR] Failed to remove worker.")
		}
	}

	fmt.Printf("[CTRL] Removing controller...\n")
	err := hamctl.RemoveApp(utils.DbuildPrefix + utils.ControllerContext + "-" + id)
	if err != nil {
		fmt.Println("[CTRL] [ERROR] Failed to remove controller.")
	}
}
