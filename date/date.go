/*
DESC: Gets the current date in the specified format
Author: Joshua Haupt
Last Modified: 07-13-2018
*/

package date

import (
	"time"
)

/*
 Constant Declarations
*/
const FILENAME_TIME_FMT = "01_02_2006"
const EMAIL_TIME_FMT = "January _2, 2006"
const LOG_TIME_FMT = "01-02-2006"

/*
DESC: Gets the current date
IN: the desired time format: format
OUT: the current date as a string in the specified format
*/
func Get_date(format string) string {

	currentDate := time.Now()

	switch {
	case format == "fileName":
		return string(currentDate.Format(FILENAME_TIME_FMT))
	case format == "email":
		return string(currentDate.Format(EMAIL_TIME_FMT))
	default:
		return string(currentDate.Format(LOG_TIME_FMT))
	}
}
