package distcc

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

const (
	shBinPath = "/bin/sh"
)

func Compile(dir string, command string, workers []string) error {
	var args []string
	args = append(args, fmt.Sprintf("export DISTCC_POTENTIAL_HOSTS=\"localhost %s\"", strings.Join(workers, " ")), ";")
	args = append(args, "cd", dir, ";")
	args = append(args, "pump", "make", command, "CC=distcc")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	cmd := exec.Command(shBinPath, fmt.Sprintf("-c '%s'", strings.Join(args, " ")))

	fmt.Println("Running: ", shBinPath, fmt.Sprintf("-c '%s'", strings.Join(args, " ")))

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed starting distcc with args (%s) with error %s", args, err)
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
