package distcc

import (
	"fmt"
	"github.com/mmkhmmkh/dbuild/pkg/utils"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

const (
	bashBinPath = "/bin/bash"
)

func Compile(dir string, command string, workers []string) error {
	defer utils.TimeTrack(time.Now(), "DISTCC")

	var args []string
	args = append(args, "(")
	args = append(args, fmt.Sprintf("export DISTCC_POTENTIAL_HOSTS=\"localhost %s\"", strings.Join(workers, " ")), ";")
	args = append(args, "export DISTCC_VERBOSE=1", ";")
	args = append(args, "export CC=distcc", ";")
	args = append(args, "cd", dir, ";")
	args = append(args, "pump --startup", ";")
	args = append(args, "mkdir -p /root/.distcc", ";")
	args = append(args, fmt.Sprintf("printf \"127.0.0.1\\n%s\" > /root/.distcc/hosts", strings.Join(workers, "\\n")), ";")
	commandParts := strings.Split(command, "&&")
	for _, commandPart := range commandParts {
		args = append(args, strings.TrimSpace(commandPart), ";")
	}
	args = append(args, "pump --shutdown")
	args = append(args, ")")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	cmd := exec.Command(bashBinPath, "-c", strings.Join(args, " "))

	fmt.Println("Running: ", bashBinPath, fmt.Sprintf("-c '%s'", strings.Join(args, " ")))

	//stderrIn, err := cmd.StderrPipe()
	//if err != nil {
	//	return fmt.Errorf("failed getting distcc err with args (%s) with error %s", args, err)
	//}
	//defer stderrIn.Close()

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

	//go func() {
	//	split := func(data []byte, atEOF bool) (advance int, token []byte, spliterror error) {
	//		if atEOF && len(data) == 0 {
	//			return 0, nil, nil
	//		}
	//		if i := bytes.IndexByte(data, '\n'); i >= 0 {
	//			// We have a full newline-terminated line.
	//			return i + 1, data[0:i], nil
	//		}
	//		if i := bytes.IndexByte(data, '\r'); i >= 0 {
	//			// We have a cr terminated line
	//			return i + 1, data[0:i], nil
	//		}
	//		if atEOF {
	//			return len(data), data, nil
	//		}
	//
	//		return 0, nil, nil
	//	}
	//	scanner := bufio.NewScanner(stderrIn)
	//	scanner.Split(split)
	//	buf := make([]byte, 2)
	//	scanner.Buffer(buf, bufio.MaxScanTokenSize)
	//	for scanner.Scan() {
	//		line := scanner.Text()
	//		fmt.Println(line)
	//	}
	//}()

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
