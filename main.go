/*
DESC: Driver - takes in CLI flags
Author: Joshua Haupt
Last Modified: 08-15-2018
*/

package main

import (
	"./app"
	"./date"
	"./email"
	"./log"
	"flag"
	"fmt"
	"os"
)

func main() {

	fmt.Printf("%s: %s\n", "Current date", date.Get_date("email"))
	optionPtr := flag.String("opt", "", "[REQUIRED] file option")
	mailToPtr := flag.String("to", "", "[REQUIRED w/ --email] mail to address")
	subjectPtr := flag.String("subject", "", "[OPTIONAL] email subject")
	mailFromPtr := flag.String("from", "", "[OPTIONAL] mail from address")
	emailPassPtr := flag.String("pass", "", "[REQUIRED w/ --email] email account password")
	headingPtr := flag.String("head", "", "[OPTIONAL] custom heading message")
	companyPtr := flag.String("company", "", "[REQUIRED] company name")
	contactPtr := flag.String("contact", "", "[OPTIONAL] contact name")
	positionPtr := flag.String("position", "", "[REQUIRED w/o --head] position name")
	sourcePtr := flag.String("source", "", "[REQUIRED w/o --head] position source")
	notePtr := flag.String("note", "", "[OPTIONAL] additional note")
	localPtr := flag.String("local", "", "[OPTIONAL] is the position local")
	skillPtr1 := flag.String("skill1", "", "[OPTIONAL] additional skill 1")
	skillPtr2 := flag.String("skill2", "", "[OPTIONAL] additional skill 2")
	urlPtr := flag.String("url", "", "[OPTIONAL] URL to postion AD")
	testPtr := flag.String("test", "", "[OPTIONAL] test build not to be logged")
	flag.Parse()

	// Make sure required company name is present if cover is generated
	if *companyPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	appl := app.App{MailTo: *mailToPtr, MailFrom: *mailFromPtr, EmailPass: *emailPassPtr, Company: *companyPtr,
		Position: *positionPtr, Source: *sourcePtr, Contact: *contactPtr, Note: *notePtr, Skill1: *skillPtr1,
		Skill2: *skillPtr2, Url: *urlPtr, Subject: *subjectPtr, Heading: *headingPtr}

	err := app.PharseFlags(*localPtr, *testPtr, *optionPtr, &appl)
	if err != nil {
		panic(err)
	}

	err = app.Build_pdf(&appl)
	if err != nil {
		panic(err)
	}

	// send Email
	if appl.MailTo != "" && appl.EmailPass != "" && appl.MailFrom != "" {

		err = email.Send_email(&appl)
		if err != nil {
			panic(err)
		} else {
			fmt.Printf("%s: %s %s: %s\n", "Email sent to", appl.MailTo, "Subject", appl.Subject)
		}

	}

	//finally log application to file
	if appl.Test == false {
		err = log.Log_app(&appl)
		if err != nil {
			panic(err)
		}
	} else if appl.Test == true {
		fmt.Println("Application not logged")
	} else {
		panic("test undefined")
	}

}
