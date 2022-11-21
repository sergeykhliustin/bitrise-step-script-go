package main

import (
	"fmt"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/creack/pty"
	"io"
	"os"
	"os/exec"
)

func main() {
	logger := log.NewLogger()
	logger.EnableDebugLog(true)
	logger.Debugf("Started")
	runner := os.Getenv("runner_bin")
	content := os.Getenv("content")

	if runner == "" {
		logger.Errorf("runner_bin is empty")
		os.Exit(1)
	}

	f, err := os.Create("._script_cont")

	if err != nil {
		logger.Errorf(err.Error())
		os.Exit(1)
	}

	_, err = f.WriteString(content)
	if err != nil {
		logger.Errorf(err.Error())
		os.Exit(1)
	}
	_ = f.Close()

	fmt.Printf("Content %s", content)

	_, err = exec.LookPath(runner)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	cmd := exec.Command(runner, "._script_cont")
	file, err := pty.Start(cmd)
	defer func() { _ = file.Close() }()
	logger.Debugf("Start copy bytes")

	_, copyErr := io.Copy(os.Stdout, file)
	logger.Debugf("End copy bytes")
	if err != nil {
		logger.Errorf(err.Error())
		os.Exit(1)
	}
	if copyErr != nil {
		logger.Errorf(err.Error())
		os.Exit(1)
	}
	if err != nil {
		logger.Errorf(err.Error())
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		logger.Errorf(err.Error())
		os.Exit(1)
	}
}
