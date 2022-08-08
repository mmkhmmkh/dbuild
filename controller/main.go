package main

import (
	"fmt"
	"github.com/mmkhmmkh/dbuild/pkg/hamctl"
	"github.com/mmkhmmkh/dbuild/pkg/utils"
	"os"
	"strconv"
	"strings"
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

func gracefulShutdown(id string) {
	fmt.Printf("[CTRL] Removing controller...\n")
	err := hamctl.RemoveApp(utils.DbuildPrefix + utils.ControllerContext + "-" + id)
	if err != nil {
		fmt.Println("[CTRL] [ERROR] Failed to remove controller.")
	}
}

// main is entry for controller node. args: [id n repo cmd env]
func main() {
	fmt.Println("#############################")
	fmt.Println("##    dbuild Controller    ##")
	fmt.Println("##   By Mahdi Khancherli   ##")
	fmt.Println("#############################")

	if len(os.Args) != 2 {
		fmt.Printf("[CTRL] [ERROR] Wrong args count.\n")
		return
	}

	args := strings.Split(os.Args[1], " ")

	if len(args) != 5 {
		fmt.Printf("[CTRL] [ERROR] Wrong args count.\n")
		return
	}

	id := args[0]
	n, _ := strconv.Atoi(args[1])
	repo := args[2]
	command := args[3]
	env := args[4]

	err := hamctl.Initialize(env)
	if err != nil {
		fmt.Printf("[CTRL] [ERROR] %v\n", err)
		return
	}

	fmt.Printf("[CTRL] Running %v workers...\n", n)
	for i := 1; i <= n; i++ {
		err := StartWorker(id, strconv.Itoa(i), fmt.Sprintf("%s %s %s", repo, command, env))
		if err != nil {
			fmt.Printf("[CTRL] [WORKER] [ERROR] %v\n", err)
			return
		}
	}

	fmt.Println("READY")
	time.Sleep(120 * time.Second)
	fmt.Println("DONE")

	fmt.Printf("[CTRL] Removing %v workers...\n", n)
	for i := 1; i < n; i++ {
		err := hamctl.RemoveApp(utils.DbuildPrefix + utils.WorkerContext + "-" + id + "-" + strconv.Itoa(i))
		if err != nil {
			fmt.Println("[CTRL] [WORKER] [ERROR] Failed to remove worker.")
		}
	}

	gracefulShutdown(id)
}
