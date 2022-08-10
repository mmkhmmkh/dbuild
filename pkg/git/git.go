package git

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

const (
	gitBinPath = "/usr/bin/git"
)

func CloneRepo(url, branch string, dir string) error {
	var args []string
	args = append(args, "clone")
	args = append(args, "-b", branch)
	args = append(args, "--depth", "1")
	args = append(args, url)
	args = append(args, dir)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	cmd := exec.Command(gitBinPath, args...)

	fmt.Println("Running: ", gitBinPath, strings.Join(args, " "))

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed starting git with args (%s) with error %s", args, err)
	}

	go func() {
		for sig := range c {
			fmt.Println("got signal ", sig)

			if err := cmd.Process.Kill(); err != nil {
				fmt.Printf("failed to kill process: %v", err)
			}
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
