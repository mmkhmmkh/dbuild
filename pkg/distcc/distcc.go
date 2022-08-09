package distcc

import (
	"bufio"
	"bytes"
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
	commandParts := strings.Split(command, "&&")
	for i, commandPart := range commandParts {
		args = append(args, "pump", strings.TrimSpace(commandPart), "CC=distcc")
		if i != len(commandParts)-1 {
			args = append(args, ";")
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	cmd := exec.Command(shBinPath, fmt.Sprintf("-c '%s'", strings.Join(args, " ")))

	fmt.Println("Running: ", shBinPath, fmt.Sprintf("-c '%s'", strings.Join(args, " ")))

	stdoutIn, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed getting distcc output with args (%s) with error %s", args, err)
	}
	defer stdoutIn.Close()

	err = cmd.Start()
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

	go func() {
		split := func(data []byte, atEOF bool) (advance int, token []byte, spliterror error) {
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}
			if i := bytes.IndexByte(data, '\n'); i >= 0 {
				// We have a full newline-terminated line.
				return i + 1, data[0:i], nil
			}
			if i := bytes.IndexByte(data, '\r'); i >= 0 {
				// We have a cr terminated line
				return i + 1, data[0:i], nil
			}
			if atEOF {
				return len(data), data, nil
			}

			return 0, nil, nil
		}
		scanner := bufio.NewScanner(stdoutIn)
		scanner.Split(split)
		buf := make([]byte, 2)
		scanner.Buffer(buf, bufio.MaxScanTokenSize)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
