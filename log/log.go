/*
DESC: Writes to log file
Author: Joshua Haupt
Last Modified: 09-22-2018
*/


package log

import (
	"os"
	"strings"
  "strconv"
	"../app"
	"../date"
)


/*
 Constant Declarations
*/
const APP_LOG_FILE = "app_log.csv"
const FOLLOW_UP_LOG_FILE = "follow_up.csv"


/*
DESC: Writes to CSV log file
IN: company name: companyName, position: position, company contact: contactPtr, position AD source: source, position AD RUL: url, mail to: email
OUT: nill on success
*/
func Log_app(appl *app.App) error {

	var logFile string
	var logString string

	if appl.Option != 10 {
		logFile = APP_LOG_FILE

		logString = date.Get_date("log") + "," + strconv.Itoa(appl.Option) + "," + strings.Replace(appl.Company, ",", "_", -1) + "," +
			strings.Replace(appl.Position, ",", "_", -1) + "," + strings.Replace(appl.PositionID, ",", "_", -1) + "," + strings.Replace(appl.Contact, ",", "_", -1) + "," +
			strings.Replace(appl.Source, ",", "_", -1) + "," + strings.Replace(appl.Heading, ",", "_", -1) + "," +
			strings.Replace(appl.Note1, ",", "_", -1) + "," + strings.Replace(appl.Note2, ",", "_", -1) + "," + strings.Replace(appl.Skill1, ",", "_", -1) + "," +
			strings.Replace(appl.Skill2, ",", "_", -1) + "," + strings.Replace(appl.Skill3, ",", "_", -1) + "," + strconv.FormatBool(appl.Local) + "," + appl.Url + "," + appl.MailTo + "\n"
	} else if appl.Option == 10 {
		logFile = FOLLOW_UP_LOG_FILE

		logString = date.Get_date("log") + "," + strconv.Itoa(appl.Option) + "," + strings.Replace(appl.Company, ",", "_", -1) + "," +
			strings.Replace(appl.Position, ",", "_", -1) + "," + strings.Replace(appl.PositionID, ",", "_", -1) + "," + strings.Replace(appl.Contact, ",", "_", -1) + "," +
			strings.Replace(appl.WhenApplied, ",", "_", -1) + "," + strings.Replace(appl.Heading, ",", "_", -1) + "," +
			strings.Replace(appl.Note1, ",", "_", -1)  + "," + appl.MailTo + "\n"
	} else {
		panic("Failed to write log")
	}

	log, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}



	_, err = log.Write([]byte(logString))
	if err != nil {
		panic(err)
	}

	err = log.Close()
	if err != nil {
		panic(err)
	}

	return nil
}
