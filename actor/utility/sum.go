package utility

import (
	"github.com/asynkron/protoactor-go/actor"
)

type SumActor struct {
	Receiver *actor.PID
	Current  int
}

func NewSumActor(receiver *actor.PID) *SumActor {
	return &SumActor{Receiver: receiver, Current: 0}
}

type SumNumber struct {
	Int int
}

type GetCurrentCalculated struct{}

type Calculated struct {
	Result int
}

func (state *SumActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case SumNumber:
		state.Current = state.Current + msg.Int
	case GetCurrentCalculated:
		context.Respond(Calculated{Result: state.Current})
	}
}
