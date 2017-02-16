package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type recHandler func(t *template.Template, opt record, c <-chan record, done chan<- struct{})

func getAction(action string, rec record) recHandler {
	isEmail := strings.Contains(action, "@")
	if isEmail {
		rec["_addr"] = os.Getenv("SMTP")
		parts := strings.SplitN(action, " ", 2) // e.g. -out "noreply@example.com Important information"
		rec["_from"] = parts[0]
		if len(parts) > 1 {
			rec["_subject"] = parts[1]
		}
		return email
	}
	return print
}

func copyRecord(dest, src record) {
	for key, value := range src {
		dest[key] = value
	}
}

func print(t *template.Template, opt record, c <-chan record, done chan<- struct{}) {
	for rec := range c {
		copyRecord(rec, opt)
		t.Execute(os.Stdout, rec)
	}
	close(done)
}

func email(t *template.Template, opt record, c <-chan record, done chan<- struct{}) {
	var buf bytes.Buffer
	for rec := range c {
		copyRecord(rec, opt)
		t.Execute(&buf, rec)
		err := sendEmail(rec["_addr"], rec["_from"], rec["_rcpt"], rec["_subject"], buf.String())
		if err != nil {
			panic(err)
		}
		buf.Reset()
	}
	close(done)
}

func sendEmail(addr, from string, rcpt string, subject, msg string) error {
	body := "From: " + from + "\r\nTo: " + rcpt + "\r\nSubject: " + subject + "\r\n\r\n" + msg
	fmt.Println(addr, nil, from, strings.Split(rcpt, ","), body)
	fmt.Println("========================================")
	return nil
	//return smtp.SendMail(addr, nil, from, strings.Split(rcpt, ","), []byte(body))
}
