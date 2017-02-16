package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

type record map[string]string

func main() {
	var data, tglob, tname, df, action, delim, sep string
	flag.StringVar(&data, "d", "", "Data file")
	flag.StringVar(&df, "df", "csv", "Data format") // todo: tsv, kv, xml, sub
	flag.StringVar(&tglob, "tg", "", "Template file glob")
	flag.StringVar(&tname, "tn", "", "Name of the template to apply.")
	flag.StringVar(&action, "out", "", `Output action; for email: "from@example.com Subject goes here", "SMTP" env var defines addr+port, "_rcpt" in data file is To-address `)
	flag.StringVar(&delim, "delim", "", `Left and right template place holder delimiters, space separated, default "{{ }}"`)
	flag.StringVar(&sep, "sep", "", `Column separator in -df csv (default ',') or key-value separator in -df kv (default '=')`)
	flag.Parse()

	// template functions
	var funcs = template.FuncMap{
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
	}
	// set up root template
	tt := template.New("")
	if delim != "" {
		lr := strings.Split(delim, " ")
		tt = tt.Delims(lr[0], lr[1])
	}
	tt, err := tt.Funcs(funcs).ParseGlob(tglob)
	if err != nil {
		panic(err)
	}
	tt.Option("missingkey=zero") // print empty string if not found

	// set the main template
	if all := tt.Templates(); len(all) == 1 && tname == "" {
		tname = all[0].ParseName // use the template if it's the only one
	}
	tmpl := tt.Lookup(tname)
	if tmpl == nil {
		fmt.Printf("Template '%s' not found\r\n", tname)
		os.Exit(1)
	}

	// set up the input
	var in io.Reader

	if data == "" {
		in = os.Stdin // read from stdin, whatever is piped through
	} else {
		// data file path is given, use that instead of stdin
		file, err := os.Open(data)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		in = file
	}

	rec := make(record)
	handleRec := getAction(action, rec) // how to process data, can fill 'rec' with options

	// run parser and handler each in their own go-routine
	// they pass data through channel c
	c := make(chan record)
	done := make(chan struct{})
	opt := make(record)
	if sep != "" {
		opt = record{"_sep": sep}
	}
	switch df {
	case "csv":
		go parseCSV(in, c, opt)
	case "tsv":
		go parseCSV(in, c, record{"_sep": "\t"})
	// case "kv":
	// 	go parseKeyValues(in, c, opt)
	}
	go handleRec(tmpl, rec, c, done)
	<-done // wait till handleRec is done
}
