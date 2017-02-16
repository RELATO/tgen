package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type recParser func(in io.Reader, out chan<- record, opt record)

func parseCSV(in io.Reader, out chan<- record, opt record) {
	var n int
	var keys []string

	r := csv.NewReader(in)
	if opt != nil {
		if sep := opt["_sep"]; len(sep) == 1 {
			r.Comma = []rune(sep)[0]
		}
	}
	for {
		fields, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}
		if n == 0 {
			n = 1
			keys = make([]string, len(fields))
			for i, value := range fields {
				keys[i] = value
			}
			continue
		}
		rec := make(record)
		for i, value := range fields {
			rec[keys[i]] = value
		}
		rec["_line"] = strconv.Itoa(n) // line number, from 1
		n++
		out <- rec
	}
	close(out)
}

func parseKeyValues(in io.Reader, out chan<- record, opt record) {
	/*
	data file:
	name=John Smith
	city=New York
	name=Jane Jones
	city=San Francisco

	template file:
	Dear {{.name}}, we know you live in {{.city}}
	*/
}

func parseBlocks(in io.Reader, out chan<- record, opt record) {
	/*
	data file:
	John Smith
	New York
	Jane Jones
	Los Angeles

	opt: -len 2
	"_len": "2"

	template file:
	Hi {{._1}} from {{._2}}
	*/

	/*
	data file:
	John Smith
	New York
	---
	Jane Jones
	Los Angeles

	opt: -sep "---"
	"_sep": "---"

	template file:
	Hi {{._1}} from {{._2}}
	*/
}

func parseXML(in io.Reader, out chan<- record, opt record) {
}

func parseFixedFields(in io.Reader, out chan<- record, opt record) {
	/*
	data file:
	First    Last      City
	John     Smith     New York
	Joanna   McSmith   San Francisco

	opt:  -pos "0 9 19"
	or    -len "9 10 *"
	-hdr 1

	template file:
	{{._1 }} {{._Last}} {{.City}}
	*/
}