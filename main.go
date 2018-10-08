/*
DESC: Driver - takes in CLI flags
Author: Joshua Haupt
Last Modified: 10-08-2018
*/

package main

import (
	"./api"
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
	optionPtr := flag.Int("opt", 0, "[REQUIRED] file option")
	headingPtr := flag.String("head", "", "[OPTIONAL] custom heading message")
	headingAddPtr := flag.String("headAdd", "", "[OPTIONAL] custom heading message")
	companyPtr := flag.String("company", "", "[REQUIRED] company name")
	contactPtr := flag.String("contact", "", "[OPTIONAL] contact name")
	positionPtr := flag.String("position", "", "[REQUIRED w/o --head] position name")
	positionIDPtr := flag.String("positionID", "", "[REQUIRED w/o --head] position ID")
	sourcePtr := flag.String("source", "", "[REQUIRED w/o --head] position source")
	notePtr1 := flag.String("note1", "", "[OPTIONAL] additional note1")
	notePtr2 := flag.String("note2", "", "[OPTIONAL] additional note2")
	localPtr := flag.Bool("local", false, "[OPTIONAL] is the position local")
	skillPtr1 := flag.String("skill1", "", "[OPTIONAL] additional skill 1")
	skillPtr2 := flag.String("skill2", "", "[OPTIONAL] additional skill 2")
	skillPtr3 := flag.String("skill3", "", "[OPTIONAL] additional skill 3")
	urlPtr := flag.String("url", "", "[OPTIONAL] URL to postion AD")
	testPtr := flag.Bool("test", false, "[OPTIONAL] test build not to be logged")
	// EMAIL
	mailToPtr := flag.String("to", "", "[REQUIRED w/ --email] mail to address")
	subjectPtr := flag.String("subject", "", "[OPTIONAL] email subject")
	// Google Cloud Storage Specific
	gcUploadPtr := flag.Bool("upload", false, "[OPTIONAL] upload file to bucket")
	// follow up
	whenAppliedPtr := flag.String("applied", "Earlier this week", "when application was submitted")
	flag.Parse()

	// Make sure required company name is present if cover is generated
	if *companyPtr == "" || *optionPtr < 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	appl := app.App{
		app.TexCover{},
		app.TextCover{},
		app.EmailCover{},
		app.Email{MailTo: *mailToPtr, Subject: *subjectPtr},
		app.GCS{GCUploadFile: *gcUploadPtr},
		app.FollowUp{WhenApplied: *whenAppliedPtr},
		app.Common{Local: *localPtr, Company: *companyPtr,
			Position: *positionPtr, PositionID: *positionIDPtr, Source: *sourcePtr, Contact: *contactPtr, Note1: *notePtr1, Note2: *notePtr2, Skill1: *skillPtr1,
			Skill2: *skillPtr2, Skill3: *skillPtr3, Url: *urlPtr, Heading: *headingPtr, HeadingAdd: *headingAddPtr},
		app.Control{Option: *optionPtr, Test: *testPtr}}

	err := app.PharseFlags(&appl)
	if err != nil {
		panic(err)
	}

	err = app.Build_pdf(&appl)
	if err != nil {
		panic(err)
	}

	// send Email
	if appl.MailTo != "" && os.Getenv("EmailPass") != "" && os.Getenv("MailFrom") != "" {

		err = email.Send_email(&appl)
		if err != nil {
			panic(err)
		} else {
			fmt.Printf("%s: %s %s: %s\n", "Email sent to", appl.MailTo, "Subject", appl.Subject)
		}

	}

	// finally log application to external database and to local CSV file
	if appl.Test == false {
		err = log.Log_app(&appl)
		if err != nil {
			panic(err)
		}
		if appl.Option != 10 {
			err = api.SendApp(&appl)
			if err != nil {
				panic(err)
			}
		}
		// TODO: follow up API call
	} else if appl.Test == true {
		fmt.Println("Application not logged")
	} else {
		panic("test undefined")
	}

}
