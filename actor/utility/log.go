package utility

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
)

type LogLevel string

const (
	INFO LogLevel = "INFO"
)

type Log struct {
	Level  LogLevel
	Format string
	Data   []any
}

func NewLog(level LogLevel, format string, data []any) *Log {
	return &Log{Level: level, Format: format, Data: data}
}

func InfoLog(str string) Log {
	return Log{Level: INFO, Format: str, Data: nil}
}

func InfoLogf(format string, data []any) Log {
	return Log{Level: INFO, Format: format, Data: data}
}

type LoggingActor struct{}

func (state *LoggingActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case Log:
		fmt.Printf("[level:%s] ", msg.Level)
		fmt.Printf("- [logger_id:%s] ", context.Self().Id)

		fmt.Printf(msg.Format, msg.Data...)
		fmt.Println()
	}
}
