/*
DESC: Writes to log file
Author: Joshua Haupt
Last Modified: 07-13-2018
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
const LOG_FILE = "app_log.csv"


/*
DESC: Writes to CSV log file
IN: company name: companyName, position: position, company contact: contactPtr, position AD source: source, position AD RUL: url, mail to: email
OUT: nill on success
*/
func Log_app(appl *app.App) error {

	log, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	logString := date.Get_date("log") + "," + appl.Cover + "," + strings.Replace(appl.Company, ",", "_", -1) + "," +
		strings.Replace(appl.Position, ",", "_", -1) + "," + strings.Replace(appl.Contact, ",", "_", -1) + "," +
		strings.Replace(appl.Source, ",", "_", -1) + "," + strings.Replace(appl.Heading, ",", "_", -1) + "," +
		strings.Replace(appl.Note, ",", "_", -1) + "," + strings.Replace(appl.Skill1, ",", "_", -1) + "," +
		strings.Replace(appl.Skill2, ",", "_", -1) + "," + strconv.FormatBool(appl.Local) + "," + appl.Url + "," + appl.EmailAddr + "\n"

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
