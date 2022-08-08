package hamctl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

const (
	hamctlBinPath    = "./tools/hamctl"
	dbuildNamespace  = "vbhammk-dbuild"
	clusterName      = "hamravesh-c11"
	organizationName = "vbhammk"
)

// CreateApp creates new app on Hamravesh PaaS
func CreateApp(appName string, repoName string, branchName string) error {
	var args []string
	args = append(args, "apps", "create")
	args = append(args, "-n", appName)
	args = append(args, "--namespace", dbuildNamespace)
	args = append(args, "-c", clusterName)
	args = append(args, "-o", organizationName)
	args = append(args, "-t", "github-repo")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	cmd := exec.Command(hamctlBinPath, args...)

	fmt.Println(hamctlBinPath, strings.Join(args, " "))

	stdoutIn, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed getting hamctl output with args (%s) with error %s", args, err)
	}

	stdinOut, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed getting hamctl input with args (%s) with error %s", args, err)
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed starting hamctl with args (%s) with error %s", args, err)
	}

	go func() {
		for sig := range c {
			fmt.Println("got signal ", sig)

			if err := cmd.Process.Kill(); err != nil {
				fmt.Printf("failed to kill process: %v", err)
			}
		}
	}()

	defer stdoutIn.Close()
	defer stdinOut.Close()

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
		//if i := bytes.IndexByte(data, ' '); i >= 0 {
		//	// We have a space-terminated word.
		//	return i + 1, data[0:i], nil
		//}
		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	}

	scanner := bufio.NewScanner(stdoutIn)
	//writer := bufio.NewWriter(stdinOut)

	scanner.Split(split)

	buf := make([]byte, 2)
	scanner.Buffer(buf, bufio.MaxScanTokenSize)
	fmt.Println("BEGINS")
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if strings.Contains(line, fmt.Sprintf("➤ %s", repoName)) {
			io.WriteString(stdinOut, "\r")
			break
		} else if strings.Contains(line, repoName) {
			io.WriteString(stdinOut, "\x1B[B")
		}
		//time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("ENDS")

	//time.Sleep(500 * time.Millisecond)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if strings.Contains(line, fmt.Sprintf("➤ %s", branchName)) {
			io.WriteString(stdinOut, "\r")
			break
		} else if strings.Contains(line, branchName) {
			io.WriteString(stdinOut, "\x1B[B")
		}
		//time.Sleep(100 * time.Millisecond)
	}

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if strings.Contains(line, "Dockerfile address") {
			io.WriteString(stdinOut, ".\r")
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if strings.Contains(line, "Build Context") {
			io.WriteString(stdinOut, ".\r")
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if strings.Contains(line, "Auto deploy after push") {
			io.WriteString(stdinOut, "no\r")
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if strings.Contains(line, "Service port") {
			io.WriteString(stdinOut, "\r")
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if strings.Contains(line, "Runnng command of your project") {
			io.WriteString(stdinOut, "\r")
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if strings.Contains(line, "Args of running command") {
			io.WriteString(stdinOut, "\r")
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if strings.Contains(line, "Health Check Path") {
			io.WriteString(stdinOut, "\r")
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	return nil
}
