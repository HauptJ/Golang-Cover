/*
DESC: Generates a custom cover letter
Author: Joshua Haupt
Last Modified: 07-07-2018
TODO: break into multiple files
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"gopkg.in/gomail.v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
	"strconv"
)


/*
 Constant Declarations
*/
const EMAIL = "josh@hauptj.com"
const EMAIL_SMTP = "smtp.gmail.com"
const TEX_COVER_TEMPL = "cover_template.tex"
const EMAIL_COVER_TEMPL = "email_cover_template.html"
const TEXT_COVER_TEMPL = "text_cover_template.txt"
const EMAIL_TEMPL = "email_template.html"
const TEX_CMD = "pdflatex"
const FILENAME_TIME_FMT = "01_02_2006"
const EMAIL_TIME_FMT = "January _2, 2006"
const LOG_TIME_FMT = "01-02-2006"
const ALL_FILENAME_BEGIN = "Joshua_Haupt_Cover_Resume_CV_"
const RESUME_FILENAME_BEGIN = "Joshua_Haupt_Resume_CV_"
const COVER_FILENAME_BEGIN = "Joshua_Haupt_Cover_"
const LOG_FILE = "app_log.csv"
const DEFAULT_CONTACT = "To whom it may concern"
const LOCAL = "I am currently located in the St. Louis area, and I am also receptive to relocation."
const DISTANT = "I am currently located in the St. Louis area, however, I am receptive to relocation."


type APP struct {
	cover string
	emailAddr string
	subject string
	emailPass string
	heading string
	company string
	contact string
	position string
	source string
	note string
	note_tex string
	note_email string
	note_text string
	local bool
	reloLine string
	reloLine_email string
	skill1 string
	skill1_tex string
	skill1_email string
	skill1_text string
	skill2 string
	skill2_tex string
	skill2_email string
	skill2_text string
	url string
	test bool
	ref bool
	kvMap_tex map[string]string
	kvMap_email map[string]string
	kvMap_text map[string]string
	attachments []string
}


/*
DESC: parses flag string values to generate APP object values
IN: the flag values as strings and the APP object
OUT: nil on success
*/
func pharseFlags(coverFlag, emailAddrFlag, subjectFlag, emailPassFlag, headingFlag, companyFlag, contactFlag, positionFlag, sourceFlag, noteFlag, localFlag,
	 skillFlag1, skillFlag2, urlFlag, testFlag, refFlag string, app *APP) error {

	// Function Level Variables
	var err error

	if contactFlag == "" {
		app.contact = DEFAULT_CONTACT
	} else {
		app.contact = contactFlag + " or to whom it may concern"
	}

	if headingFlag == "" && positionFlag != "" && sourceFlag != "" {
		app.heading = "I am excited about the possibility of joining your organization in the position of " + app.position + ", as advertised on " + app.source + "." // default heading
	} else if headingFlag != "" {
		app.heading = headingFlag
	} else {
		if coverFlag == "incl" || coverFlag == "sep" {
			panic("heading undefined")
		}
	}

	// If additional note is present, add a newline at the end of it
	if noteFlag != "" {
		app.note_tex = app.note + " \\newline"
		app.note_email = "<div><p style=\"text-align:left\";>" + app.note + "</p></div>"
	}

	// If additional skill is present, add a \item before it
	if skillFlag1 != "" {
		app.skill1_tex = "\\item " + app.skill1
		app.skill1_email = "<li>" + app.skill1 + "</li>"
		app.skill1_text = "- " + app.skill1
	}

	// If additional skill is present, add a \item before it
	if skillFlag2 != "" {
		app.skill2_tex = "\\item " + app.skill2
		app.skill2_email = "<li>" + app.skill2 + "</li>"
		app.skill2_text = "- " + app.skill2
	}


	app.local, err = parseBool(localFlag, true)
	if err != nil {
		panic(err)
	}

	if app.local == true {
		app.reloLine = LOCAL
	} else if app.local == false {
		app.reloLine = DISTANT
	} else {
		panic("Local undefined")
	}
	app.reloLine_email = "<div><p style=\"text-align:left;\">" + app.reloLine + "</p></div>"


	if coverFlag == "incl" || coverFlag == "sep" {
		app.kvMap_tex = map[string]string{"[COMPANY_NAME]": strings.Replace(app.company, "&", "\\&", -1), "[COMPANY_CONTACT]": strings.Replace(app.contact, "&", "\\&", -1),
			"[POSITION_NAME]": strings.Replace(app.position, "&", "\\&", -1), "[HEADING]": strings.Replace(app.heading, "&", "\\&", -1), "[POSITION_SOURCE]": strings.Replace(app.source, "&", "\\&", -1),
			"[ADDITIONAL_SKILL_1]": strings.Replace(app.skill1_tex, "&", "\\&", -1), "[ADDITIONAL_SKILL_2]": strings.Replace(app.skill2_tex, "&", "\\&", -1),
			"[ADDITIONAL_NOTE]": strings.Replace(app.note_tex, "&", "\\&", -1), "[RELOCATION]": strings.Replace(app.reloLine, "&", "\\&", -1)}

		app.kvMap_text = map[string]string{"[COMPANY_NAME]": app.company, "[COMPANY_CONTACT]": app.contact, "[POSITION_NAME]": app.position,
			"[HEADING]": app.heading, "[POSITION_SOURCE]": app.source, "[ADDITIONAL_SKILL_1]": app.skill1_text, "[ADDITIONAL_SKILL_2]": app.skill2_text,
			"[ADDITIONAL_NOTE]": app.note, "[CURRENT_DATE]": get_date("email"), "[RELOCATION]": app.reloLine}
	}

	if emailAddrFlag != "" && emailPassFlag != "" {
		app.kvMap_email = map[string]string{"[COMPANY_NAME]": app.company, "[COMPANY_CONTACT]": app.contact, "[POSITION_NAME]": app.position,
			"[HEADING]": app.heading, "[POSITION_SOURCE]": app.source, "[ADDITIONAL_SKILL_1]": app.skill1_email, "[ADDITIONAL_SKILL_2]": app.skill2_email,
			"[ADDITIONAL_NOTE]": app.note_email, "[CURRENT_DATE]": get_date("email"), "[RELOCATION]": app.reloLine_email}


		if subjectFlag == "" && positionFlag != "" {
			app.subject = "Joshua Haupt application for " + app.position + " position at " + app.company // default subject
		} else if subjectFlag != "" {
			app.subject = subjectFlag
		} else {
			flag.PrintDefaults()
			os.Exit(1)
		}

	}

	app.test, err = parseBool(testFlag, false)
	if err != nil {
		panic(err)
	}

	app.ref, err = parseBool(refFlag, true)
	if err != nil {
		panic(err)
	}

	return nil
}


/*
DESC: parses a string for a bool value and if a blank string is provided, returns specified default
IN: the string to pharse and a default bool value
OUT: the bool value of the parsed string
*/
func parseBool(input string, deflt bool) (output bool, err error) {
	if input == "" { //default if nothing specified
		output = deflt
	} else if input != "" {
		output, err = strconv.ParseBool(input)
		if err != nil {
			panic(err)
		}
	} else {
		panic("input undefined")
	}
	return output, nil
}


/*
DESC: runs a bash command
IN: command name: cmdName, command arguments cmdArgs
OUT: nill on success
SOURCE: https://nathanleclaire.com/blog/2014/12/29/shelled-out-commands-in-golang/
*/
func run_cmd(cmdName string, cmdArgs []string) error {

	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for CMD", err)
		panic(err)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("command output | %s\n", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting CMD", err)
		panic(err)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for CMD", err)
		panic(err)
	}

	return nil
}


/*
DESC: Gets the current date
IN: the desired time format: format
OUT: the current date as a string in the specified format
*/
func get_date(format string) string {

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


/*
DESC: Generates Cover, Resume, and CV PDFs
IN: APP object app
company name: company,
OUT: the current date as a string
*/
func build_pdf(app *APP) error {

	var cmdArgs []string
	var cmdArgs_cover []string
	var cmdArgs_resume []string

	contents, err := replace_strings(TEX_COVER_TEMPL, app.kvMap_tex)
	if err != nil {
		panic(err)
	}
	err = write_file("cover.tex", contents)

	if app.cover == "incl" {
		if app.ref == false {
			cmdArgs = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_all\".tex"}
		} else {
			cmdArgs = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_all_ref\".tex"}
		}

		err = run_cmd(TEX_CMD, cmdArgs)
		if err != nil {
			panic(err)
		}

	} else if app.cover == "sep" {

		if app.ref == false {
			cmdArgs_resume = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_resume\".tex"}
		} else {
			cmdArgs_resume = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_resume_ref\".tex"}
		}

		cmdArgs_cover = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_cover\".tex"}

		err := run_cmd(TEX_CMD, cmdArgs_resume)
		if err != nil {
			panic(err)
		}

		err = run_cmd(TEX_CMD, cmdArgs_cover)
		if err != nil {
			panic(err)
		}

	} else { // Just the CV

		if app.ref == false {
			cmdArgs = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_resume\".tex"}
		} else {
			cmdArgs = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_resume_ref\".tex"}
		}

		err = run_cmd(TEX_CMD, cmdArgs)
		if err != nil {
			panic(err)
		}

	}

	return nil
}


/*
SOURCE: https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
DESC: copy_file copies a file from src to dst. If src and dst files exist, and are
   the same, then return success. Otherise, attempt to create a hard link
   between the two files. If that fail, copy the file contents from src to dst.
*/
func copy_file(src, dst string) (err error) {

	sfi, err := os.Stat(src)
	if err != nil {
		return
	}

	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("copy_file: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}

	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("copy_file: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}

	if err = os.Link(src, dst); err == nil {
		return
	}

	err = copy_file_contents(src, dst)
	return
}


/*
SOURCE: https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
DESC: copy_fileContents copies the contents of the file named src to the file named
   by dst. The file will be created if it does not already exist. If the
   destination file exists, all it's contents will be replaced by the contents
   of the source file.
*/
func copy_file_contents(src, dst string) (err error) {

	in, err := os.Open(src)
	if err != nil {
		return
	}

	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}

	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()

	return
}


/*
DESC: Renames PDF cover letters and resumes
IN: APP object app
OUT: nill on success
*/
func rename_files(app *APP) error {

	var err error

	// Replace spaces in company name with _
	companyName := strings.Replace(app.company, " ", "_", -1)

	if app.cover == "incl" {

		newName := ALL_FILENAME_BEGIN + companyName + "_" + get_date("fileName") + ".pdf"
		fmt.Printf("%s: %s\n", "Output File", newName)

		if app.ref == false {
			err = copy_file("main_all.pdf", newName)
			if err != nil {
				panic(err)
			}
		} else {
				err = copy_file("main_all_ref.pdf", newName)
				if err != nil {
					panic(err)
				}
		}

		app.attachments = append(app.attachments, newName)

	} else if app.cover == "sep" {

		newName_resume := RESUME_FILENAME_BEGIN + companyName + "_" + get_date("fileName") + ".pdf"
		fmt.Printf("%s: %s\n", "Output File", newName_resume)

		if app.ref == false {
			err = copy_file("main_resume.pdf", newName_resume)
			if err != nil {
				panic(err)
			}
		} else {
				err = copy_file("main_resume_ref.pdf", newName_resume)
				if err != nil {
					panic(err)
				}
		}

		app.attachments = append(app.attachments, newName_resume)

		newName_cover := COVER_FILENAME_BEGIN + companyName + "_" + get_date("fileName") + ".pdf"
		fmt.Printf("%s: %s\n", "Output File", newName_cover)

		err = copy_file("main_cover.pdf", newName_cover)
		if err != nil {
			panic(err)
		}

		app.attachments = append(app.attachments, newName_cover)

	} else { // Just the CV

		newName_resume := RESUME_FILENAME_BEGIN + companyName + "_" + get_date("fileName") + ".pdf"
		fmt.Printf("%s: %s\n", "Output File", newName_resume)

		if app.ref == false {
			err = copy_file("main_resume.pdf", newName_resume)
			if err != nil {
				panic(err)
			}
		} else {
				err = copy_file("main_resume_ref.pdf", newName_resume)
				if err != nil {
					panic(err)
				}
		}

		app.attachments = append(app.attachments, newName_resume)
	}

	return nil
}


/*
DESC: Reads a file into a buffer and replaces strings key in map with coresponding values
IN: input file: in_file, key value map kvMap
OUT: new string contents: contents, nill on success
*/
func replace_strings(inFile string, kvMap map[string]string) (string, error) {

	// read original file contents into buffer
	buffer, err := ioutil.ReadFile(inFile)
	if err != nil {
		panic(err)
	}
	contents := string(buffer)

	// replace strings using a map
	// keys k are the old values being replaced
	// values v are the values replacing the coresponding keys
	for k, v := range kvMap {
		contents = strings.Replace(contents, k, v, -1)
	}

	return contents, nil
}


/*
DESC: writes to new file
IN: output file: outFile, string of contents to write: contents
OUT: nill on success
*/
func write_file(outFile, contents string) error {

	err := ioutil.WriteFile(outFile, []byte(contents), 0644)
	if err != nil {
		panic(err)
	}

	return nil
}


/*
DESC: sends email via SMTP
IN: APP object app
OUT: nill on success
*/
func send_email(app *APP) error {

	var body string
	var err error

	if app.cover == "incl" || app.cover == "sep" {

		body, err = replace_strings(EMAIL_COVER_TEMPL, app.kvMap_email)
		if err != nil {
			panic(err)
		}

	} else { //send plain email
		body, err = replace_strings(EMAIL_TEMPL, app.kvMap_email)
		if err != nil {
			panic(err)
		}
	}

	m := gomail.NewMessage()
	m.SetHeader("From", EMAIL)
	m.SetHeader("To", app.emailAddr)
	m.SetHeader("Subject", app.subject)
	m.SetBody("text/html", body)

	for _, attachment := range app.attachments {
		m.Attach(attachment)
	}

	d := gomail.NewDialer(EMAIL_SMTP, 587, EMAIL, app.emailPass)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return nil
}


/*
DESC: generates copyable plain text version of cover letter
IN: Application object app
OUT: nill on success
*/
func text_cover(app *APP) error {
	contents, err := replace_strings(TEXT_COVER_TEMPL, app.kvMap_text)
	if err != nil {
		panic(err)
	}
	err = write_file("cover.txt", contents)
	if err != nil {
		panic(err)
	}

	return nil
}


/*
DESC: Writes to CSV log file
IN: company name: companyName, position: position, company contact: contactPtr, position AD source: source, position AD RUL: url, mail to: email
OUT: nill on success
*/
func log_app(app *APP) error {

	log, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	logString := get_date("log") + "," + app.cover + "," + strings.Replace(app.company, ",", "_", -1) + "," +
		strings.Replace(app.position, ",", "_", -1) + "," + strings.Replace(app.contact, ",", "_", -1) + "," +
		strings.Replace(app.source, ",", "_", -1) + "," + strings.Replace(app.heading, ",", "_", -1) + "," +
		strings.Replace(app.note, ",", "_", -1) + "," + strings.Replace(app.skill1, ",", "_", -1) + "," +
		strings.Replace(app.skill2, ",", "_", -1) + "," + strconv.FormatBool(app.local) + "," + app.url + "," + app.emailAddr + "\n"

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


func main() {

	fmt.Printf("%s: %s\n", "Current date", get_date("email"))

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

	app := APP{cover: *coverPtr, emailAddr: *emailAddrPtr, emailPass: *emailPassPtr, company: *companyPtr, position: *positionPtr, source: *sourcePtr,
		note: *notePtr, skill1: *skillPtr1, skill2: *skillPtr2, url: *urlPtr}

	err := pharseFlags(*coverPtr, *emailAddrPtr, *subjectPtr, *emailPassPtr, *headingPtr, *companyPtr, *contactPtr, *positionPtr, *sourcePtr, *notePtr,
		*localPtr, *skillPtr1, *skillPtr2, *urlPtr, *testPtr, *refPtr, &app)
	if err != nil {
		panic(err)
	}

	err = build_pdf(&app)
	if err != nil {
		panic(err)
	}

	err = text_cover(&app)
	if err != nil {
		panic(err)
	}

	err = rename_files(&app)
	if err != nil {
		panic(err)
	}

	// send Email
	if app.emailAddr != "" && app.emailPass != "" {

		err = send_email(&app)
		if err != nil {
			panic(err)
		} else {
			fmt.Printf("%s: %s %s: %s\n", "Email sent to", app.emailAddr, "Subject", app.subject)
		}

	}

	//finally log application to file
	if app.test == false {
		err = log_app(&app)
		if err != nil {
			panic(err)
		}
	} else if app.test == true {
		fmt.Println("Application not logged")
	} else {
		panic("test undefined")
	}

}
