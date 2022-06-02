package loader

import (
	"bufio"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"os"
)

type Load struct {
	CSVPath  string
	Receiver *actor.PID
}
type HeroCSVLoader struct{}

func (state *HeroCSVLoader) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case Load:
		csvFile, err := os.Open(msg.CSVPath)
		if err != nil {
			println(err)
			return
		}
		csvScanner := bufio.NewScanner(csvFile)
		csvScanner.Split(bufio.ScanLines)
		for csvScanner.Scan() {
			fmt.Println(csvScanner.Text())
		}

		if err := csvFile.Close(); err != nil {
			println(err)
			return
		}
	}
}
