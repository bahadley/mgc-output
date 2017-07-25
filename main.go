package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"unicode"
)

const (
	HeartbeatEvent = "H"
	Query          = "Q"
	Verdict        = "V"
)

var (
	mgcDataPath string
	node        string
	eventType   string
)

type Tuple struct {
	Node      string
	EventType string
	EventTime string
	SeqNo     int
	Delay     string
	Verdict   string
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
		mgcTuples[i].SeqNo, err = strconv.Atoi(items[2])
		check(err)
		if mgcTuples[i].EventType == Query {
			mgcTuples[i].Delay = items[3]
		} else if mgcTuples[i].EventType == HeartbeatEvent {
			mgcTuples[i].Delay = items[3]
		} else if mgcTuples[i].EventType == Verdict {
			mgcTuples[i].Verdict = items[3]
		}
	}

	output(mgcTuples)
}

func output(mgcTuples []Tuple) {
	for _, t := range mgcTuples {
		if t.Node == node {
			if eventType == HeartbeatEvent && t.EventType == HeartbeatEvent {
				fmt.Printf("%s,%s\n", t.EventTime, t.Delay)
			} else if eventType == Query && t.EventType == Query {
				fmt.Printf("%s,%s\n", t.EventTime, t.Delay)
			} else if t.EventType == Verdict {
				fmt.Printf("%s,%s\n", t.EventTime, t.Verdict)
			}
		}
	}
}

func init() {
	flag.StringVar(&mgcDataPath, "file", "",
		"mgc data file")

	flag.StringVar(&node, "node", "",
		"node id")

	flag.StringVar(&eventType, "type", "",
		"event type [(H)eartbeat, (Q)uery, (V)erdict]")
}
