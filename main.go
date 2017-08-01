package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	HeartbeatEvent = "H"
	Query          = "Q"
	Verdict        = "V"

	HeaderLines = 5

	OutputHeader = "time,delay"
)

var (
	mgcDataPath string
	node        string
	eventType   string
	baseTime    time.Time
)

type Tuple struct {
	Node      string
	EventType string
	EventTime time.Time
	SeqNo     int
	Delay     string
	Verdict   string
}

type MgcTuples []Tuple

func (slice MgcTuples) Len() int {
	return len(slice)
}

func (slice MgcTuples) Less(i, j int) bool {
	return (slice[i].EventTime).Before(slice[j].EventTime)
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

		if unixNano, err := strconv.ParseInt(items[1], 10, 64); err == nil {
			mgcTuples[idx].EventTime = time.Unix(0, unixNano)
		} else {
			panic(err)
		}
		if i == HeaderLines {
			baseTime = mgcTuples[idx].EventTime
		}

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
	var outputted bool

	for _, t := range mgcTuples {
		if t.Node == node {
			elapsedTime := t.EventTime.Sub(baseTime)
			elapsedMillis := elapsedTime.Nanoseconds() / 1000000

			if !outputted {
				fmt.Println(OutputHeader)
				outputted = true
			}

			if eventType == HeartbeatEvent && t.EventType == HeartbeatEvent {
				fmt.Printf("%d,%s\n", elapsedMillis, t.Delay)
			} else if eventType == Query && t.EventType == Query {
				fmt.Printf("%d,%s\n", elapsedMillis, t.Delay)
			} else if eventType == Verdict && t.EventType == Verdict {
				fmt.Printf("%d,%s\n", elapsedMillis, t.Verdict)
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
