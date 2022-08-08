package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mmkhmmkh/dbuild/pkg/hamctl"
	"strconv"
)

const (
	dbuildRepo        = "mmkhmmkh/dbuild"
	controllerContext = "controller"
	workerContext     = "worker"
	dbuildBranch      = "master"
	dbuildPrefix      = "dbuild-"
)

var controllersCount int

func StartController(arguments string) error {
	fmt.Printf("[ORCH] Starting new controller...\n")
	err := hamctl.CreateApp(dbuildPrefix+controllerContext+"-"+strconv.Itoa(controllersCount+1), dbuildRepo, dbuildBranch, controllerContext, "/dbuild/bin/controller", arguments)
	if err != nil {
		return err
	}

	fmt.Printf("[ORCH] New controller started successfuly.\n")
	controllersCount++
	return nil
}

// main is entry for orchestrator node
func main() {

	fmt.Println("#############################")
	fmt.Println("##   dbuild Orchestrator   ##")
	fmt.Println("##   By Mahdi Khancherli   ##")
	fmt.Println("#############################")

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/submit", func(c *gin.Context) {
		n := c.Query("n")
		repo := c.Query("repo")
		cmd := c.Query("cmd")
		err := StartController(fmt.Sprintf("%s %s %s", n, repo, cmd))
		if err != nil {
			fmt.Printf("[ORCH] [CONTROLLER] [ERROR] %v\n", err)
			return
		}
	})

	err := r.Run("localhost:8080")
	if err != nil {
		fmt.Printf("[ORCH] [ERROR] %v\n", err)
	}

	fmt.Println("[ORCH] Bye!")

}
