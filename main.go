package main

import (
	"bufio"
	"fmt"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/creack/pty"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	logger := log.NewLogger()
	runner := os.Getenv("runner_bin")
	content := os.Getenv("content")
	tmpDir := os.Getenv("TMPDIR")
	scriptPath := tmpDir + "/._script_cont"
	workingDir := os.Getenv("working_dir")
	scriptFilePath := os.Getenv("script_file_path")
	isDebug := os.Getenv("is_debug")
	logger.EnableDebugLog(isDebug == "yes")
	timestampEnabled := os.Getenv("timestamp") == "yes"

	logger.Debugf("==> Start")

	if runner == "" {
		logger.Errorf("runner_bin is empty")
		os.Exit(1)
	}

	if workingDir != "" {
		err := os.Chdir(workingDir)
		if err != nil {
			logger.Errorf(" [!] Failed to switch to working directory: %s", workingDir)
		}
	}

	if scriptFilePath != "" {
		logger.Debugf("==> Script (tmp) save path specified: %s", scriptFilePath)
		scriptPath = scriptFilePath
	}

	f, err := os.Create(scriptPath)

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

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	runnerArgs := strings.Fields(runner)
	runnerBin := runnerArgs[0]
	args := runnerArgs[1:]
	args = append(args, scriptPath)

	cmd := exec.Command(runnerBin, args...)
	file, err := pty.Start(cmd)
	defer func() { _ = file.Close() }()
	logger.Debugf("Start logging")
	if timestampEnabled {
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			timestamp := []byte(time.Now().Format("15:04:05") + " ")
			text := append(scanner.Bytes(), []byte("\n")...)
			_, _ = os.Stdout.Write(append(timestamp, text...))
		}
	} else {
		_, copyErr := io.Copy(os.Stdout, file)
		if copyErr != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}
	}
	logger.Debugf("End end logging")
	if err != nil {
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
	os.Exit(0)
}
