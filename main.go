package main

import (
	"flag"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ryota0624/proto-actor-sandbox/actor/hero/loader"
	"os"
)

func main() {
	flag.Parse()
	csvFilePath := flag.Arg(0)
	if csvFilePath == "" {
		println("require arg")
		os.Exit(1)
	}

	sys := actor.NewActorSystem()
	context := actor.NewRootContext(sys, nil)
	props := actor.PropsFromProducer(func() actor.Actor { return &loader.HeroCSVLoader{} })
	pid := context.Spawn(props)
	context.Send(pid, loader.Load{
		CSVPath:  csvFilePath,
		Receiver: nil,
	})
	var a int
	fmt.Scan(&a)
	println("shutdown")
	sys.Shutdown()
}
