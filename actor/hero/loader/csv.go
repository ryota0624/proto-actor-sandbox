package loader

import (
	"bufio"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/jszwec/csvutil"
	"github.com/ryota0624/proto-actor-sandbox/model"
	"os"
	"strings"
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

		var header string
		if csvScanner.Scan() {
			header = csvScanner.Text()
		}

		for csvScanner.Scan() {
			csv := strings.Join([]string{header, csvScanner.Text()}, "\n")
			fmt.Printf("%s\n", csv)

			var records []model.Hero
			if err := csvutil.Unmarshal([]byte(csv), &records); err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			context.Send(msg.Receiver, records)
		}

		if err := csvFile.Close(); err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	}
}
