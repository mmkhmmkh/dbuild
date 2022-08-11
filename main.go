package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mmkhmmkh/dbuild/pkg/hamctl"
	"github.com/mmkhmmkh/dbuild/pkg/utils"
	"os"
	"strconv"
)

var controllersCount int

func StartController(arguments string) error {
	fmt.Printf("[ORCH] Starting new controller...\n")
	err := hamctl.CreateApp(utils.DbuildPrefix+utils.ControllerContext+"-"+strconv.Itoa(controllersCount+1), utils.DbuildRepo, utils.DbuildBranch, utils.ControllerContext, "/dbuild/bin/"+utils.ControllerContext, strconv.Itoa(controllersCount+1)+" "+arguments)
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
	// http://localhost:8080/submit?n=2&command=cp%20minimalconfig%20.config%20%26%26%20make%20-j8%20bzImage&repo=https://github.com/liva/minimal-linux.git&branch=master
	err := hamctl.Initialize(os.Getenv("HAMCTLCONFIG"))
	if err != nil {
		fmt.Printf("[ORCH] [ERROR] %v\n", err)
		return
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/submit", func(c *gin.Context) {
		n := c.Query("n")
		repo := c.Query("repo")
		branch := c.Query("branch")
		command := c.Query("command")
		err := StartController(fmt.Sprintf("%s %s %s %s %s", n, repo, branch, os.Getenv("HAMCTLCONFIG"), command))
		if err != nil {
			fmt.Printf("[ORCH] [CTRL] [ERROR] %v\n", err)
			return
		}
	})

	err = r.Run("localhost:8080")
	if err != nil {
		fmt.Printf("[ORCH] [ERROR] %v\n", err)
	}

	fmt.Println("[ORCH] Bye!")

}
