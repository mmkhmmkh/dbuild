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
	// http://localhost:8080/submit?n=2&command=cp%20minimalconfig%20.config%20%26%26%20make%20-j8%20bzImage&repo=https://github.com/liva/minimal-linux.git&branch=master
	// http://localhost:8080/submit?n=2&command=cp%20minimalconfig%20.config%20%26%26%20make%20-j8%20bzImage&repo=https://github.com/liva/minimal-linux.git&branch=master&output=arch/x86/boot/bzImage&s3endpoint=https://kise-thr-nd-1.sotoon.cloud&s3bucket=delivery-neda-mmkh&s3access=0e54adbb5ddc063081dfd9212c0234744e607c12&s3secret=e8334fa7da55e28d370a6bd4bd8341b3c8e31b0fd1342285597f38713b83708d70a06ca88a4fdbad
	// http://localhost:8080/submit?n=2&command=cp%20minimalconfig%20.config%20%26%26%20make%20-j12%20bzImage%20CC%3Ddistcc&repo=https://github.com/liva/minimal-linux.git&branch=master&output=arch/x86/boot/bzImage&s3endpoint=kise-thr-nd-1.sotoon.cloud&s3bucket=delivery-neda-mmkh&s3access=0e54adbb5ddc063081dfd9212c0234744e607c12&s3secret=e8334fa7da55e28d370a6bd4bd8341b3c8e31b0fd1342285597f38713b83708d70a06ca88a4fdbad
	// 1 2 https://github.com/liva/minimal-linux.git master arch/x86/boot/bzImage https://kise-thr-nd-1.sotoon.cloud delivery-neda-mmkh 0e54adbb5ddc063081dfd9212c0234744e607c12 e8334fa7da55e28d370a6bd4bd8341b3c8e31b0fd1342285597f38713b83708d70a06ca88a4fdbad {"ApiKey":"dca8733d-bb05-4be1-8e91-fc0b80f820f0","DefaultOrganization":{"id":258,"name":"vbhammk","available_clusters":null}} cp minimalconfig .config && make -j12 bzImage CC=distcc
	// 1 2 https://github.com/liva/minimal-linux.git master arch/x86/boot/bzImage https://kise-thr-nd-1.sotoon.cloud delivery-neda-mmkh 0e54adbb5ddc063081dfd9212c0234744e607c12 e8334fa7da55e28d370a6bd4bd8341b3c8e31b0fd1342285597f38713b83708d70a06ca88a4fdbad {"ApiKey":"dca8733d-bb05-4be1-8e91-fc0b80f820f0","DefaultOrganization":{"id":258,"name":"vbhammk","available_clusters":null}} cp minimalconfig .config && make -j12 bzImage CC=distcc
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
