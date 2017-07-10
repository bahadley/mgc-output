package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var (
	mgcDataPath string
)

type Tuple struct {
	Timestamp int64
	Delay     float64
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
		if len(l) == 0 {
			continue
		}

		vals := strings.Split(l, ",")
		mgcTuples[i].Timestamp, err = strconv.ParseInt(vals[0], 10, 64)
		check(err)
		mgcTuples[i].Delay, err = strconv.ParseFloat(vals[1], 64)
		check(err)

		fmt.Printf("%d,%d\n", mgcTuples[i].Timestamp, mgcTuples[i].Delay)
	}
}

func init() {
	flag.StringVar(&mgcDataPath, "input", "",
		"mgc data file path")
}
