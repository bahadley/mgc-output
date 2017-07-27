package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

const (
	HeartbeatEvent = "H"
	Query          = "Q"
	Verdict        = "V"

	HeaderLines = 5
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

type MgcTuples []Tuple

func (slice MgcTuples) Len() int {
	return len(slice)
}

func (slice MgcTuples) Less(i, j int) bool {
	l, erl := strconv.Atoi(slice[i].EventTime)
	check(erl)
	r, err := strconv.Atoi(slice[j].EventTime)
	check(err)
	return l < r
}

func (slice MgcTuples) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
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
	mgcTuples := make(MgcTuples, len(mgcLines)-HeaderLines-1)
	idx := 0

	for i, l := range mgcLines {
		if i < HeaderLines || len(l) == 0 {
			continue
		}

		vals := strings.Fields(l)
		mgcTuples[idx].Node = strings.Trim(vals[0], "<>:")

		f := func(c rune) bool {
			return !unicode.IsLetter(c) && !unicode.IsNumber(c)
		}
		items := strings.FieldsFunc(vals[4], f)
		mgcTuples[idx].EventType = items[0]
		mgcTuples[idx].EventTime = items[1]
		mgcTuples[idx].SeqNo, err = strconv.Atoi(items[2])
		check(err)
		if mgcTuples[idx].EventType == Query {
			mgcTuples[idx].Delay = items[3]
		} else if mgcTuples[idx].EventType == HeartbeatEvent {
			mgcTuples[idx].Delay = items[3]
		} else if mgcTuples[idx].EventType == Verdict {
			mgcTuples[idx].Verdict = items[3]
		}

		idx++
	}

	sort.Sort(mgcTuples)
	output(mgcTuples)
}

func output(mgcTuples []Tuple) {
	for _, t := range mgcTuples {
		if t.Node == node {
			if eventType == HeartbeatEvent && t.EventType == HeartbeatEvent {
				fmt.Printf("%s,%s\n", t.EventTime, t.Delay)
			} else if eventType == Query && t.EventType == Query {
				fmt.Printf("%s,%s\n", t.EventTime, t.Delay)
			} else if eventType == Verdict && t.EventType == Verdict {
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
