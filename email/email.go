/*
DESC: Writes to log file
Author: Joshua Haupt
Last Modified: 0-22-2018
*/


package email

import (
	"gopkg.in/gomail.v2"
	"../app"
	"fmt"
	"os"
)


/*
 Constant Declarations
*/
const EMAIL_SMTP = "smtp.gmail.com"
const EMAIL_COVER_TEMPL = "email_cover_template.html"
const EMAIL_TEMPL = "email_template.html"
const EMAIL_FOLLOW_UP_TEMPL = "email_follow_up_template.html"


/*
DESC: sends email via SMTP
IN: App object app
OUT: nill on success
*/
func Send_email(appl *app.App) error {

	var body string
	var err error

	if appl.Option <= 6 && appl.Option > 0 {

		body, err = app.Replace_strings(EMAIL_COVER_TEMPL, appl.KvMap_email)
		if err != nil {
			panic(err)
		}

	} else if appl.Option == 10 {

		body, err = app.Replace_strings(EMAIL_FOLLOW_UP_TEMPL, appl.KvMap_email)
		if err != nil {
			panic(err)
		}

	} else { //send plain email
		body, err = app.Replace_strings(EMAIL_TEMPL, appl.KvMap_email)
		if err != nil {
			panic(err)
		}
	}

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("MailFrom"))
	m.SetHeader("To", appl.MailTo)
	m.SetHeader("Subject", appl.Subject)
	m.SetBody("text/html", body)

	for _, attachment := range appl.Attachments {
		m.Attach(attachment)
		fmt.Printf("Attached file: %s\n", attachment)
	}

	if len(appl.Attachments) <= 0 { // abort program if there are no attachments
		os.Exit(2)
	}

	d := gomail.NewDialer(EMAIL_SMTP, 587, os.Getenv("MailFrom"), os.Getenv("EmailPass"))

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return nil
}
