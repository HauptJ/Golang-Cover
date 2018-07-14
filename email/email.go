/*
DESC: Writes to log file
Author: Joshua Haupt
Last Modified: 07-13-2018
*/


package email

import (
	"gopkg.in/gomail.v2"
	"../app"
)


/*
 Constant Declarations
*/
const EMAIL = "josh@hauptj.com"
const EMAIL_SMTP = "smtp.gmail.com"
const EMAIL_COVER_TEMPL = "email_cover_template.html"
const EMAIL_TEMPL = "email_template.html"


/*
DESC: sends email via SMTP
IN: App object app
OUT: nill on success
*/
func Send_email(appl *app.App) error {

	var body string
	var err error

	if appl.Cover == "incl" || appl.Cover == "sep" {

		body, err = app.Replace_strings(EMAIL_COVER_TEMPL, appl.KvMap_email)
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
	m.SetHeader("From", EMAIL)
	m.SetHeader("To", appl.EmailAddr)
	m.SetHeader("Subject", appl.Subject)
	m.SetBody("text/html", body)

	for _, attachment := range appl.Attachments {
		m.Attach(attachment)
	}

	d := gomail.NewDialer(EMAIL_SMTP, 587, EMAIL, appl.EmailPass)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return nil
}
