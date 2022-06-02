package main

import (
	"flag"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"github.com/asynkron/protoactor-go/router"
	"github.com/asynkron/protoactor-go/stream"
	"github.com/ryota0624/proto-actor-sandbox/model"
	"strconv"
	"time"

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

	heroWeightSumPid := context.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return utility.NewSumActor(loggerPid)
	}))

	mapToHeroWightSumPid := context.Spawn(
		utility.NewMapper[model.Hero, utility.SumNumber](heroWeightSumPid, func(hero model.Hero) (utility.SumNumber, error) {
			weight, err := strconv.Atoi(hero.Weight)
			if err != nil {
				return utility.SumNumber{}, err
			}
			return utility.SumNumber{
				Int: weight,
			}, nil
		}))

	context.Send(loggerPid, utility.InfoLog("Hello Logger"))

	mapToLoggerPid := context.Spawn(router.NewRoundRobinGroup(
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
			context.Send(mapToHeroWightSumPid, hero)
		}
	}()

	composedHeroProcessor := context.Spawn(router.NewBroadcastGroup(
		mapToLoggerPid,
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

	f := context.RequestFuture(heroWeightSumPid, utility.GetCurrentCalculated{}, 3*time.Second)
	result, err := f.Result()
	if err != nil {
		fmt.Printf("err=%+v\n", err)
	} else {
		fmt.Printf("result=%+v\n", result)

	}
	sys.Shutdown()
}
