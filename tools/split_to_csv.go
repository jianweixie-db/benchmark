// Read the source data file and prepare the csv file for spark usage.
package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	kRoughTokenCount = 512

	// Deduced.
	kRoughCharCount = kRoughTokenCount * 6
)

func main() {
	inputs := readData()
	records := prepareRowsForInputs(inputs)

	w := csv.NewWriter(os.Stdout)
	w.WriteAll(records) // calls Flush internally

	if err := w.Error(); err != nil {
		panic(err)
	}
}

func readData() string {
	bs, err := ioutil.ReadFile("data/shakespeare.txt")
	if err != nil {
		panic(err)
	}

	return string(bs)
}

func prepareRowsForInputs(inputs string) [][]string {
	outs := make([][]string, 0)

	index := 0
	for {
		if index+kRoughCharCount > len(inputs) {
			break
		}

		row := string(inputs[index : index+kRoughCharCount])
		row = strings.ReplaceAll(row, "\n", "")
		index += kRoughCharCount

		outs = append(outs, []string{row})
	}

	fmt.Fprintf(os.Stderr, "procssed %v records in total\n", len(outs))
	return outs
}
