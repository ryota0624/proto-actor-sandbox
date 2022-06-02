package main

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
)

type Hello struct{ Who string }
type HelloActor struct{}

func (state *HelloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case Hello:
		fmt.Printf("Hello %v\n", msg.Who)
	}
}

func main() {
	sys := actor.NewActorSystem()
	context := actor.NewRootContext(sys, nil)
	props := actor.PropsFromProducer(func() actor.Actor { return &HelloActor{} })
	pid := context.Spawn(props)
	context.Send(pid, Hello{Who: "Roger"})
	var a int
	fmt.Scan(&a)
	println("shutdown")
	sys.Shutdown()
}
