package utility

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
)

type Map struct {
	Data any
}

type MapActor[From any, To any] struct {
	Fn       func(From) (To, error)
	Receiver *actor.PID
}

func NewMapper[From any, To any](receiver *actor.PID, fn func(From) (To, error)) *actor.Props {
	return actor.PropsFromProducer(func() actor.Actor {
		return &MapActor[From, To]{
			Fn:       fn,
			Receiver: receiver,
		}
	})
}

func NewNonErrorMapper[From any, To any](receiver *actor.PID, fn func(From) To) *actor.Props {
	return actor.PropsFromProducer(func() actor.Actor {
		return &MapActor[From, To]{
			Fn: func(f From) (To, error) {
				r := fn(f)
				return r, nil
			},
			Receiver: receiver,
		}
	})
}

func (state *MapActor[T, C]) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case T:
		mapped, err := state.Fn(msg)
		if err != nil {
			fmt.Printf("%+v\n", err)
			return
		}

		context.Send(state.Receiver, mapped)
	}
}
