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
	// arch/x86/boot/bzImage
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

	r.GET("/clean", func(c *gin.Context) {
		for i := 1; i < 5; i++ {
			hamctl.RemoveApp(utils.DbuildPrefix + utils.ControllerContext + "-" + strconv.Itoa(i))
			for j := 1; j < 10; j++ {
				hamctl.RemoveApp(utils.DbuildPrefix + utils.WorkerContext + "-" + strconv.Itoa(i) + "-" + strconv.Itoa(j))
			}
		}
	})

	r.GET("/submit", func(c *gin.Context) {
		n := c.Query("n")
		repo := c.Query("repo")
		branch := c.Query("branch")
		command := c.Query("command")
		output := c.Query("output")
		s3endpoint := c.Query("s3endpoint")
		s3bucket := c.Query("s3bucket")
		s3access := c.Query("s3access")
		s3secret := c.Query("s3secret")
		err := StartController(fmt.Sprintf("%s %s %s %s %s %s %s %s %s %s", n, repo, branch, output, s3endpoint, s3bucket, s3access, s3secret, os.Getenv("HAMCTLCONFIG"), command))
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
