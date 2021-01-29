package main

import (
	"fmt"
	"strconv"
	"io/ioutil"
)

func main() {
	wasmOpts := SSVMWasmOptions {
		wasmFile: "/root/yomo-flow-ssvm-example/triple/pkg/triple.wasm",
	}
	result := Run(wasmOpts, []string{"+6.913980e+015"})
	if s, err := strconv.ParseFloat(string(result), 64); err == nil {
		fmt.Println(s)
	}


	/**
	 *  Tensorflow sample
	 */
	wasmOpts = SSVMWasmOptions {
		wasmFile: "/root/yomo-flow-ssvm-example/triple/pkg/tensorflow.wasm",
		cmdEnvs: map[string]string {
			"LD_LIBRARY_PATH": "/root/ssvm-tensorflow-lib", // required by ssvm-tensorflow
		},
		reactor: true,
		dirs: map[string]string {
			"host_path": "guest_path",
		},
		envs: map[string]string {
			"APP_TOKEN": "sample case",
		},
	}

	content, _ := ioutil.ReadFile("./Hamburger.jpg")
	// content, _ := ioutil.ReadFile("./Ramen.jpg")

	result = RunTensorflow(wasmOpts, content)
	fmt.Println(string(result))
}
