package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	mgcDataPath string
)

type Tuple struct {
	Data string
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
		mgcTuples[i].Data = vals[4]
		check(err)

		fmt.Println(mgcTuples[i].Data)
	}
}

func init() {
	flag.StringVar(&mgcDataPath, "input", "",
		"mgc data file path")
}
