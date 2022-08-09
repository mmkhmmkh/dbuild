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

func Initialize(hamctlconfig string) error {
	if hamctlconfig == "" {
		return fmt.Errorf("no HAMCTLCONFIG")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	err = os.WriteFile(homeDir+"/.hamctlconfig", []byte(hamctlconfig), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func RemoveApp(appName string) error {
	var args []string
	args = append(args, "apps", "del")
	args = append(args, appName)
	args = append(args, "-n", dbuildNamespace)
	args = append(args, "-c", clusterName)
	args = append(args, "-o", organizationName)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	cmd := exec.Command(hamctlBinPath, args...)

	fmt.Println("Running: ", hamctlBinPath, strings.Join(args, " "))

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
		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	}

	scanner := bufio.NewScanner(stdoutIn)

	scanner.Split(split)

	buf := make([]byte, 2)
	scanner.Buffer(buf, bufio.MaxScanTokenSize)
	//fmt.Println("BEGINS")
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		if strings.Contains(line, "Enter app name to confirm") {
			io.WriteString(stdinOut, fmt.Sprintf("%s\r", appName))
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		if strings.Contains(line, "✗") && !strings.Contains(line, "Enter app name to confirm") {
			return fmt.Errorf("hamctl error: %s", strings.TrimSpace(strings.ReplaceAll(line, "✗ ", "")))
		} else if strings.Contains(line, "deleted successfully") {
			return nil
		}
	}

	return fmt.Errorf("hamctl unexpected behaviour")
}

// CreateApp creates new app on Hamravesh PaaS
func CreateApp(appName string, repoName string, branchName string, buildContext string, command string, arguments string) error {
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

	fmt.Println("Running: ", hamctlBinPath, strings.Join(args, " "))

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
	//fmt.Println("BEGINS")
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		if strings.Contains(line, fmt.Sprintf("➤ %s", repoName)) {
			//fmt.Println("About to Accept " + line)
			io.WriteString(stdinOut, "\r")
			break
		} else if strings.Contains(line, repoName) {
			//fmt.Println("HERE")
			io.WriteString(stdinOut, "\x1B[B")
		}
		//time.Sleep(100 * time.Millisecond)
	}

	//fmt.Println("ENDS")

	//time.Sleep(500 * time.Millisecond)

	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
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
		//fmt.Println(line)
		if strings.Contains(line, "Dockerfile address") {
			io.WriteString(stdinOut, fmt.Sprintf("./%s.Dockerfile\r", buildContext))
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		if strings.Contains(line, "Build Context") {
			io.WriteString(stdinOut, fmt.Sprintf(".\r", buildContext))
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		if strings.Contains(line, "Auto deploy after push") {
			io.WriteString(stdinOut, "yes\r")
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		if strings.Contains(line, "Service port") {
			io.WriteString(stdinOut, "\r")
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		if strings.Contains(line, "command of your project") {
			io.WriteString(stdinOut, fmt.Sprintf("%s\r", command))
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		if strings.Contains(line, "Args of running command") {
			io.WriteString(stdinOut, fmt.Sprintf("%s\r", arguments))
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		if strings.Contains(line, "Health Check Path") {
			io.WriteString(stdinOut, "\r")
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		if strings.Contains(line, "Plan 1") {
			io.WriteString(stdinOut, "\x1B[B")
			io.WriteString(stdinOut, "\x1B[B")
			io.WriteString(stdinOut, "\x1B[B")
			io.WriteString(stdinOut, "\x1B[B")
			io.WriteString(stdinOut, "\r")
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "✗") {
			return fmt.Errorf("hamctl error: %s", strings.TrimSpace(strings.ReplaceAll(line, "✗ ", "")))
		} else if strings.Contains(line, "created successfully") {
			return nil
		}
	}

	return fmt.Errorf("hamctl unexpected behaviour")
}
