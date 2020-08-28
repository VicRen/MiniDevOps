package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/vicren/minidevops/core/app/task"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/vicren/minidevops/core/app/exporter"
	"github.com/vicren/minidevops/core/common/serial"

	"github.com/vicren/minidevops/core"
)

func main() {
	flag.Parse()

	config := &core.Config{
		Features: []*any.Any{
			serial.MarshalAny(&task.Config{}),
			serial.MarshalAny(&exporter.Config{
				Address: "0.0.0.0:10988",
			}),
		},
	}
	s, err := core.New(config)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Start(); err != nil {
		log.Fatalf("failed to start: %v", err)
	}
	defer s.Close()

	runtime.GC()

	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-osSignals
	}
}
