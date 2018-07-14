/*
DESC: Driver - takes in CLI flags
Author: Joshua Haupt
Last Modified: 07-13-2018
*/


package main

import (
	"flag"
	"fmt"
	"os"
	"./app"
	"./date"
	"./email"
	"./log"
)


func main() {

	fmt.Printf("%s: %s\n", "Current date", date.Get_date("email"))

	coverPtr := flag.String("cover", "no", "[OPTIONAL] incl, sep or no cover letter")
	emailAddrPtr := flag.String("to", "", "[REQUIRED w/ --email] mail to address")
	subjectPtr := flag.String("subject", "", "[OPTIONAL] email subject")
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
	refPtr := flag.String("ref", "", "[OPTIONAL] include references")
	flag.Parse()

	// Make sure required company name is present if cover is generated
	if *companyPtr == "" && (*coverPtr == "incl" || *coverPtr == "sep") {
		flag.PrintDefaults()
		os.Exit(1)
	}

	appl := app.App{Cover: *coverPtr, EmailAddr: *emailAddrPtr, EmailPass: *emailPassPtr, Company: *companyPtr,
		Position: *positionPtr, Source: *sourcePtr, Contact: *contactPtr, Note: *notePtr, Skill1: *skillPtr1,
		Skill2: *skillPtr2, Url: *urlPtr, Subject: *subjectPtr, Heading: *headingPtr}

	err := app.PharseFlags(*localPtr, *testPtr, *refPtr, &appl)
	if err != nil {
		panic(err)
	}

	err = app.Build_pdf(&appl)
	if err != nil {
		panic(err)
	}

	// send Email
	if appl.EmailAddr != "" && appl.EmailPass != "" {

		err = email.Send_email(&appl)
		if err != nil {
			panic(err)
		} else {
			fmt.Printf("%s: %s %s: %s\n", "Email sent to", appl.EmailAddr, "Subject", appl.Subject)
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
