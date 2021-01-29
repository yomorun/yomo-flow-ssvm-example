package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	y3 "github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/yomo/pkg/quic"
)

var (
	addr = "0.0.0.0:4141"
)

func main() {
	go serveSinkServer(addr)
	select {}
}

// serveSinkServer serves the Sink server over QUIC.
func serveSinkServer(addr string) {
	log.Print("Starting sink server...")
	quicServer := quic.NewServer(&quicServerHandler{})
	err := quicServer.ListenAndServe(context.Background(), addr)
	if err != nil {
		log.Printf("❌ Serve the yomo-ssvm-example on %s failure with err: %v", addr, err)
	}
}

type quicServerHandler struct {
}

func (s *quicServerHandler) Listen() error {
	return nil
}

func (s *quicServerHandler) Read(st quic.Stream) error {
	// decode the data via Y3 Codec.
	ch := y3.
		FromStream(st).
		Subscribe(0x10).
		OnObserve(onObserve)

	go func() {
		for item := range ch {
			// Invoke wasm
			val := triple(item.(noiseData).Noise)
			println(val)
		}
	}()

	return nil
}

type noiseData struct {
	Noise float64 `yomo:"0x11" fauna:"noise"` // Noise value
	Time  int64   `yomo:"0x12" fauna:"time"`  // Timestamp (ms)
	From  string  `yomo:"0x13" fauna:"from"`  // Source IP
}

func onObserve(v []byte) (interface{}, error) {
	// decode the data via Y3 Codec.
	data := noiseData{}
	err := y3.ToObject(v, &data)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return data, nil
}

func triple(i float64) float64 {
	f := fmt.Sprintf("%f", i)
	wasmOpts := SSVMWasmOptions {
		wasmFile: "/root/yomo-flow-ssvm-example/triple/pkg/triple.wasm",
	}
	result := Run(wasmOpts, []string{f})
	s, _ := strconv.ParseFloat(string(result), 64)
	return s
}

