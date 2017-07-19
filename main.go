package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
	"unicode"
)

const (
	HeartbeatEvent = "H"
	Query          = "Q"
	Verdict        = "V"
)

var (
	mgcDataPath string
)

type Tuple struct {
	Node      string
	EventType string
	EventTime string
	SeqNo     uint16
	Delay     time.Duration
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	flag.Parse()

	mgcRaw, err := ioutil.ReadFile(mgcDataPath)
	check(err)

	mgcLines := strings.Split(string(mgcRaw), "\n")
	mgcTuples := make([]Tuple, len(mgcLines))

	for i, l := range mgcLines {
		if i < 5 || len(l) == 0 {
			continue
		}

		vals := strings.Fields(l)
		mgcTuples[i].Node = strings.Trim(vals[0], "<>:")

		f := func(c rune) bool {
			return !unicode.IsLetter(c) && !unicode.IsNumber(c)
		}
		items := strings.FieldsFunc(vals[4], f)
		mgcTuples[i].EventType = items[0]
		mgcTuples[i].EventTime = items[1]

		if mgcTuples[i].EventType == HeartbeatEvent ||
			mgcTuples[i].EventType == Query {
			fmt.Printf("%s,%s,%s\n", mgcTuples[i].Node,
				mgcTuples[i].EventType, mgcTuples[i].EventTime)
		}
	}
}

func init() {
	flag.StringVar(&mgcDataPath, "input", "",
		"mgc data file path")
}
