package main

import (
	"flag"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"github.com/asynkron/protoactor-go/router"
	"github.com/asynkron/protoactor-go/stream"
	"github.com/ryota0624/proto-actor-sandbox/model"

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
		return &utility.LoggingActor{
			Logger: log.New(log.DebugLevel, "logging-actor"),
		}
	})
	loggerPid := context.Spawn(loggerProps)
	context.Send(loggerPid, utility.InfoLog("Hello Logger"))

	mapperPid := context.Spawn(router.NewRoundRobinGroup(
		context.Spawn(utility.NewNonErrorMapper[any, utility.Log](loggerPid, func(from any) utility.Log {
			return utility.InfoLogf("mapper 1 %+v", []any{from})
		})),
		context.Spawn(utility.NewNonErrorMapper[any, utility.Log](loggerPid, func(from any) utility.Log {
			return utility.InfoLogf("mapper 2 %+v", []any{from})
		})),
		context.Spawn(utility.NewNonErrorMapper[any, utility.Log](loggerPid, func(from any) utility.Log {
			return utility.InfoLogf("mapper 3 %+v", []any{from})
		})),
	))

	heroBatchStream := stream.NewTypedStream[[]model.Hero](sys)
	heroStream := stream.NewTypedStream[model.Hero](sys)

	heroStreamLogger := log.New(log.InfoLevel, "hero-stream")

	go func() {
		for heroes := range heroBatchStream.C() {
			for _, hero := range heroes {
				context.Send(heroStream.PID(), hero)
			}
		}
	}()

	go func() {
		for hero := range heroStream.C() {
			heroStreamLogger.Info("received Hero", log.Object("hero", hero))
		}
	}()

	composedHeroProcessor := context.Spawn(router.NewBroadcastGroup(
		mapperPid,
		heroBatchStream.PID(),
	))

	props := actor.PropsFromProducer(func() actor.Actor { return &loader.HeroCSVLoader{} })
	pid := context.Spawn(props)
	context.Send(pid, loader.Load{
		CSVPath:  csvFilePath,
		Receiver: composedHeroProcessor,
	})

	var a int
	fmt.Scan(&a)
	println("shutdown")
	sys.Shutdown()
}
