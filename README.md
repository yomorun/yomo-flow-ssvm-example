![ssvm-yomo](https://raw.githubusercontent.com/yomorun/yomo-flow-ssvm-example/main/yomo-ssvm-2.png)

# yomo-flow-ssvm-example

This examples represents how to write a [yomo-flow](https://yomo.run/flow) with [ssvm](https://www.secondstate.io/ssvm)

## Prerequisites

Install the [SSVM shared library and SSVM-go](https://github.com/second-state/ssvm-go)

## Compile wasm file

```bash
# Install wasm-pack tool
curl https://rustwasm.github.io/wasm-pack/installer/init.sh -sSf | sh

cd triple
wasm-pack build
```

will get the wasm file at `triple/pkg/` directory

## Run this flow

1. Run this `yomo-flow`

```bash
go run main.go ./triple/pkg/triple_lib_bg.wasm triple
```

Then start your [yomo source](https://yomo.run/source) and [yomo zipper](https://yomo.run/zipper)

2. Run `yomo-zipper`

```bash
yomo wf run workflow.yaml
```

3. Run `yomo-source-example`

```bash
gh repo clone yomorun/yomo-source-example
cd yomo-source-example
YOMO_ZIPPER_ENDPOINT=localhost:9999 go run main.go
```

you'll get: 

```bash
Ξ _wrk/yomo-sink-ssvm → go run main.go
2021/01/26 00:14:20 Starting sink server...
2021/01/26 00:14:20 ✅ Listening on 0.0.0.0:4141
+6.913980e+015
+9.036794e+015
^Csignal: interrupt
```

## Code

The `triple/src/lib.rs` source file will be built into `pkg/triple_lib_bg.wasm`, which is then passed to `vm.RunWasmFile` function as a command argument.

The definition `RunWasmFile(path string, funcname string, params ...interface{}) ([]interface{}, error)` shows that the second parameter is the calling function name which is also passed as a command argument. Other params will be passed to the calling function in the wasm. Please attention that the type of params and return values can only be i32, i64, f32, f64 currently. [Check source](https://github.com/second-state/ssvm-go/blob/master/ssvm/vm.go)

The other function for calling wasm function is `RunWasmFileWithDataAndWASI` which can receive bytes array as parameters. There is an [example](https://github.com/second-state/ssvm-go/blob/master/examples/go_PassBytes/pass_bytes.go) for this.


# Additional information

## Why WebAssembly in edge computing
In an edge computing framework (e.g., YOMO), we often need to execute user submitted code in the dataflow to handle application-specific logic. Since the user submitted code is potentially unsafe and poorly tested, it is crucial that they are executed in a sandbox for security and cross-platform compatibility. WebAssembly is an application sandbox that is very light and very fast—making it an ideal choice for resource-constrained and real-time edge computing scenarios.

Compared with prevailing application containers like Docker, WebAssembly provides a higher level of abstraction, and hence higher productivity, for developers. WebAssembly can [deploy functions instantly](https://www.secondstate.io/articles/getting-started-with-rust-function/), instead of launching an operating system and then a language runtime, making it suitable for real-time and high-performance applications.

## Why SSVM

The [SSVM (Second State VM)](https://github.com/second-state/SSVM) is a popular WebAssembly VM optimized for high-performance applications on the server-side. With advanced AOT (Ahead of Time compiler) support, the SSVM is already one of the fastest WebAssembly VMs.

> Reference: [A Lightweight Design for High-performance Serverless Computing](https://arxiv.org/abs/2010.07115), published on IEEE Software, Jan 2021.

A key differentiator of the SSVM is its support for WebAssembly extensions. For example, the SSVM  supports an extension for Tensorflow and other AI frameworks. Developers can write AI inference functions using [a simple Rust API](https://crates.io/crates/ssvm_tensorflow_interface), and then run the function securely at native speed on CPU / GHPU / AI chips from the SSVM. That is an excellent fit for edge computing.

## YOMO and SSVM

Through its support for the WebAssembly System Interface (WASI), the SSVM can be easily started and managed by the golang-based YOMO host environment, making it a good choice for YOMO extensions. In the future, the SSVM will also support a golang SDK based on WebAssembly interface types specification and make golang / WebAssembly interoperations seamless. SSVM also supports high-performance AI inference, which is commonly required in edge computing settings. YOMO and SSVM is a perfect match.
