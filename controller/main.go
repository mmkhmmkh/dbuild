package main

import (
	"fmt"
	"github.com/mmkhmmkh/dbuild/pkg/distcc"
	"github.com/mmkhmmkh/dbuild/pkg/git"
	"github.com/mmkhmmkh/dbuild/pkg/hamctl"
	"github.com/mmkhmmkh/dbuild/pkg/s3"
	"github.com/mmkhmmkh/dbuild/pkg/utils"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	CloneDirectory = "repo"
)

func StartWorker(controllerID, workerID string) error {

	workerName := utils.DbuildPrefix + utils.WorkerContext + "-" + controllerID + "-" + workerID
	fmt.Printf("[CTRL] Starting new worker (%s)...\n", workerName)

	hamctl.RemoveApp(workerName)

	//err := hamctl.CreateApp(utils.DbuildPrefix+utils.WorkerContext+"-"+controllerID+"-"+workerID, utils.DbuildRepo, utils.DbuildBranch, utils.WorkerContext, "/dbuild/bin/"+utils.WorkerContext, workerID+" "+arguments)
	err := hamctl.CreateApp(utils.DbuildPrefix+utils.WorkerContext+"-"+controllerID+"-"+workerID, utils.DbuildRepo, utils.DbuildBranch, utils.WorkerContext, "", "")
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

// main is entry for controller node. args: [id n repo branch output s3endpoint s3bucket s3access s3secret env commands...]
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

	if len(args) < 7 {
		fmt.Printf("[CTRL] [ERROR] Wrong args count.\n")
		return
	}

	id := args[0]
	n, _ := strconv.Atoi(args[1])
	repo := args[2]
	branch := args[3]
	output := args[4]
	s3endpoint := args[5]
	s3bucket := args[6]
	s3access := args[7]
	s3secret := args[8]
	env := args[9]
	var commands []string
	for i := 10; i < len(args); i++ {
		commands = append(commands, args[i])
	}
	command := strings.Join(commands, " ")

	err := hamctl.Initialize(env)
	if err != nil {
		fmt.Printf("[CTRL] [ERROR] %v\n", err)
		return
	}

	fmt.Printf("[CTRL] Running %v workers...\n", n)

	var workers []string

	for i := 1; i <= n; i++ {
		//err := StartWorker(id, strconv.Itoa(i), fmt.Sprintf("%s %s %s", repo, command, env))
		err := StartWorker(id, strconv.Itoa(i))
		if err != nil {
			fmt.Printf("[CTRL] [WORKER] [ERROR] %v\n", err)
			return
		}

		workers = append(workers, utils.DbuildPrefix+utils.WorkerContext+"-"+id+"-"+strconv.Itoa(i))
	}

	fmt.Printf("[CTRL] Workers created. Waiting for them to be ready...\n")

	time.Sleep(2 * time.Minute)

	fmt.Printf("[CTRL] Workers ready.\n")

	err = git.CloneRepo(repo, branch, CloneDirectory)
	if err != nil {
		fmt.Printf("[CTRL] [ERROR] %v\n", err)
		return
	}

	fmt.Printf("[CTRL] Git clone completed.\n")

	err = distcc.Compile(CloneDirectory, command, workers)
	if err != nil {
		fmt.Printf("[CTRL] [ERROR] %v\n", err)
	}

	fmt.Printf("[CTRL] Compiled!\n")

	s3client, err := s3.NewS3Client(s3endpoint, s3access, s3secret, true)
	if err != nil {
		fmt.Printf("[CTRL] [ERROR] %v\n", err)
	}

	fmt.Printf("[CTRL] S3 Connected.\n")

	err = s3.Upload(utils.DbuildPrefix+utils.ControllerContext+"-"+id, s3bucket, utils.DbuildDir+"/"+CloneDirectory+"/"+output, s3client)
	if err != nil {
		fmt.Printf("[CTRL] [ERROR] %v\n", err)
	}

	fmt.Printf("[CTRL] Uploaded.\n")

	time.Sleep(1 * time.Hour)

	//fmt.Printf("[CTRL] Removing %v workers...\n", n)
	//for i := 1; i < n; i++ {
	//	err := hamctl.RemoveApp(utils.DbuildPrefix + utils.WorkerContext + "-" + id + "-" + strconv.Itoa(i))
	//	if err != nil {
	//		fmt.Println("[CTRL] [WORKER] [ERROR] Failed to remove worker.")
	//	}
	//}
	//
	//gracefulShutdown(id)
}
