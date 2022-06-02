package main

import (
	"flag"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ryota0624/proto-actor-sandbox/actor/hero/loader"
	"github.com/ryota0624/proto-actor-sandbox/actor/utility"
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

	loggerProps := actor.PropsFromProducer(func() actor.Actor {
		return &utility.LoggingActor{}
	})
	loggerPid := context.Spawn(loggerProps)
	context.Send(loggerPid, utility.InfoLog("Hello Logger"))

	println("Start loader")
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
