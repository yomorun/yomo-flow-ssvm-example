package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type SSVMWasmOptions struct {
	wasmFile string // wasm file path
	cmdEnvs map[string]string // command line env
	reactor bool // ssvm option
	dirs map[string]string // ssvm option
	envs map[string]string // ssvm option
}

func Run(wasmOptions SSVMWasmOptions, args interface{}) []byte {
	// Prerequisite:
	// Build ssvm
	// https://github.com/second-state/ssvm
	cmd := prepareCmd("ssvm", wasmOptions)
	return run(cmd, args)
}

func RunTensorflow(wasmOptions SSVMWasmOptions, args interface{}) []byte {
	// Prerequisite:
	// Build ssvm-tensorflow and download the required shared libraries
	// Then LD_LIBRARY_PATH should be passed via cmdEnvs
	// https://github.com/second-state/ssvm-tensorflow
	cmd := prepareCmd("ssvm-tensorflow", wasmOptions)
	return run(cmd, args)
}

func run(cmd exec.Cmd, args interface{}) []byte {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		switch v := args.(type) {
		case []string:
			x, _ := json.Marshal(v)
			io.WriteString(stdin, string(x))
		case []byte:
			stdin.Write(v)
		}
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	return out
}

func prepareCmd(cmdName string, wasmOptions SSVMWasmOptions) exec.Cmd {
	cmdArgs := []string{wasmOptions.wasmFile}
	if wasmOptions.reactor {
		cmdArgs = append(cmdArgs, "--reactor")
	}
	for dk, dv := range wasmOptions.dirs {
		cmdArgs = append(cmdArgs, "--dir", fmt.Sprintf("%s:'%s'", dk, dv))
	}
	for ek, ev := range wasmOptions.envs {
		cmdArgs = append(cmdArgs, "--env", fmt.Sprintf("%s='%s'", ek, ev))
	}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Env = os.Environ()
	for ek, ev := range wasmOptions.cmdEnvs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s='%s'", ek, ev))
	}

	return *cmd
}

