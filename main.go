package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"text/template"
)

func main() {
	var data, tmpl, action, apply string
	flag.StringVar(&data, "d", "", "Data file")
	flag.StringVar(&tmpl, "tg", "", "Template file glob")
	flag.StringVar(&apply, "tn", "", "Name of the template to apply.")
	flag.StringVar(&action, "out", "", `Output action; for email: "from@example.com Subject goes here", "SMTP" env var defines addr+port, "_rcpt" in data file is To-address `)
	flag.Parse()

	var funcs = template.FuncMap{
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
	}
	tt, err := template.New("").Funcs(funcs).ParseGlob(tmpl)
	if err != nil {
		panic(err)
	}
	dd, err := os.Open(data)
	if err != nil {
		panic(err)
	}
	defer dd.Close()

	actionEmail := strings.Contains(action, "@")
	var addr, from, subject string
	if actionEmail {
		addr = os.Getenv("SMTP")
		parts := strings.SplitN(action, " ", 2) // e.g. -out "noreply@example.com Important information"
		from = parts[0]
		if len(parts) > 1 {
			subject = parts[1]
		}
	}

	var n int
	var hdr []string
	m := make(map[string]string)
	scanner := bufio.NewScanner(dd)
	for scanner.Scan() {
		row := scanner.Text()
		fields := strings.Split(row, ",")
		if n == 0 {
			n = 1
			hdr = make([]string, len(fields))
			for i, value := range fields {
				hdr[i] = value
			}
			continue
		}
		for i, value := range fields {
			m[hdr[i]] = value
		}
		m["_line"] = fmt.Sprintf("%d", n) // line number, from 1

		if action == "" {
			tt.ExecuteTemplate(os.Stdout, apply, m)
		} else if strings.Contains(action, "@") {
			to := strings.Split(m["_rcpt"], ",")
			var buf bytes.Buffer
			tt.ExecuteTemplate(&buf, apply, m)
			email(addr, from, to, subject, buf.String())
		}
		n++
	}
}

func email(addr, from string, to []string, subject, msg string) error {
	body := "From: " + from + "\r\nTo: " + strings.Join(to, ",") + "\r\nSubject: " + subject + "\r\n\r\n" + msg
	// fmt.Println(addr, from, to, body)
	// return nil
	return smtp.SendMail(addr, nil, from, to, []byte(body))
}
