/*
DESC: Generates Resume and Cover Letter
Author: Joshua Haupt
Last Modified: 08-17-2018
*/


package app

import (
  "io/ioutil"
	"strings"
	"strconv"
  "fmt"
  "../date"
  "../copyfiles"
  "../cmd"
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
const ALL_FILENAME_BEGIN = "Joshua_Haupt_Cover_Resume_CV_"
const RESUME_FILENAME_BEGIN = "Joshua_Haupt_Resume_"
const COVER_FILENAME_BEGIN = "Joshua_Haupt_Cover_"
const CV_FILENAME_BEGIN = "Joshua_Haupt_CV_"
const DEFAULT_CONTACT = "To whom it may concern"
const CONTACT_ENDING = " or to whom it may concern"
const LOCAL = "I am currently located in the St. Louis area, and I am also receptive to relocation."
const DISTANT = "I am currently located in the St. Louis area, however, I am receptive to relocation."



type App struct {
  Option int
	MailTo string
  MailFrom string
	Subject string
	EmailPass string
	Heading string
	Company string
	Contact string
	Position string
	Source string
	Note string
	Note_tex string
	Note_email string
	Note_text string
	Local bool
	ReloLine string
	ReloLine_email string
	Skill1 string
	Skill1_tex string
	Skill1_email string
	Skill1_text string
	Skill2 string
	Skill2_tex string
	Skill2_email string
	Skill2_text string
	Url string
	Test bool
	Ref bool
	KvMap_tex map[string]string
	KvMap_email map[string]string
	KvMap_text map[string]string
	Attachments []string
}


func GetAttachments(appl *App) []string {
  for _, attachment := range appl.Attachments {
    fmt.Println(attachment)
  }
  return appl.Attachments
}

/*
DESC: parses flag string values to generate App object values
IN: the flag values as strings and the App object
OUT: nil on success
*/
func PharseFlags(localFlag, testFlag, optionFlag string, appl *App) error {

	// Function Level Variables
	var err error

  appl.Option, err = strconv.Atoi(optionFlag)

	if appl.Contact == "" {
		appl.Contact = DEFAULT_CONTACT
	} else {
		appl.Contact = appl.Contact + CONTACT_ENDING
	}

	if appl.Heading == "" && appl.Position != "" && appl.Source != "" {
		appl.Heading = "I am excited about the possibility of joining your organization in the position of " + appl.Position + ", as advertised on " + appl.Source + "." // default heading
	} else {
		if appl.Option <= 6 && appl.Option > 0 {
			panic("heading undefined")
		}
	}

	// If additional note is present, add a newline at the end of it
	if appl.Note != "" {
		appl.Note_tex = appl.Note + " \\newline"
		appl.Note_email = "<div><p style=\"text-align:left\";>" + appl.Note + "</p></div>"
	}

	// If additional skill is present, add a \item before it
	if appl.Skill1 != "" {
		appl.Skill1_tex = "\\item " + appl.Skill1
		appl.Skill1_email = "<li>" + appl.Skill1 + "</li>"
		appl.Skill1_text = "- " + appl.Skill1
	}

	// If additional skill is present, add a \item before it
	if appl.Skill2 != "" {
		appl.Skill2_tex = "\\item " + appl.Skill2
		appl.Skill2_email = "<li>" + appl.Skill2 + "</li>"
		appl.Skill2_text = "- " + appl.Skill2
	}


	appl.Local, err = parseBool(localFlag, true)
	if err != nil {
		panic(err)
	}

	if appl.Local == true {
		appl.ReloLine = LOCAL
	} else if appl.Local == false {
		appl.ReloLine = DISTANT
	} else {
		panic("Local undefined")
	}
	appl.ReloLine_email = "<div><p style=\"text-align:left;\">" + appl.ReloLine + "</p></div>"


	if appl.Option <= 6 && appl.Option > 0 {
		appl.KvMap_tex = map[string]string{"[COMPANY_NAME]": strings.Replace(appl.Company, "&", "\\&", -1), "[COMPANY_CONTACT]": strings.Replace(appl.Contact, "&", "\\&", -1),
			"[POSITION_NAME]": strings.Replace(appl.Position, "&", "\\&", -1), "[HEADING]": strings.Replace(appl.Heading, "&", "\\&", -1), "[POSITION_SOURCE]": strings.Replace(appl.Source, "&", "\\&", -1),
			"[ADDITIONAL_SKILL_1]": strings.Replace(appl.Skill1_tex, "&", "\\&", -1), "[ADDITIONAL_SKILL_2]": strings.Replace(appl.Skill2_tex, "&", "\\&", -1),
			"[ADDITIONAL_NOTE]": strings.Replace(appl.Note_tex, "&", "\\&", -1), "[RELOCATION]": strings.Replace(appl.ReloLine, "&", "\\&", -1)}

		appl.KvMap_text = map[string]string{"[COMPANY_NAME]": appl.Company, "[COMPANY_CONTACT]": appl.Contact, "[POSITION_NAME]": appl.Position,
			"[HEADING]": appl.Heading, "[POSITION_SOURCE]": appl.Source, "[ADDITIONAL_SKILL_1]": appl.Skill1_text, "[ADDITIONAL_SKILL_2]": appl.Skill2_text,
			"[ADDITIONAL_NOTE]": appl.Note, "[CURRENT_DATE]": date.Get_date("email"), "[RELOCATION]": appl.ReloLine}
	}

	if appl.MailTo != "" && appl.EmailPass != ""  && appl.MailFrom != "" {
		appl.KvMap_email = map[string]string{"[COMPANY_NAME]": appl.Company, "[COMPANY_CONTACT]": appl.Contact, "[POSITION_NAME]": appl.Position,
			"[HEADING]": appl.Heading, "[POSITION_SOURCE]": appl.Source, "[ADDITIONAL_SKILL_1]": appl.Skill1_email, "[ADDITIONAL_SKILL_2]": appl.Skill2_email,
			"[ADDITIONAL_NOTE]": appl.Note_email, "[CURRENT_DATE]": date.Get_date("email"), "[RELOCATION]": appl.ReloLine_email}


		if appl.Subject == "" && appl.Position != "" {
			appl.Subject = "Joshua Haupt appllication for " + appl.Position + " position at " + appl.Company // default subject
		}

	}

	appl.Test, err = parseBool(testFlag, false)
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
DESC: Generates Cover, Resume, and CV PDFs
IN: App object app
company name: company,
OUT: the current date as a string
*/
func Build_pdf(appl *App) error {

  var err error
	var cmdArgs []string
	var cmdArgs_cover []string
	var cmdArgs_resume []string
  var cmdArgs_CV []string

  if appl.Option <= 6 && appl.Option > 0 {
    contents, err := Replace_strings(TEX_COVER_TEMPL, appl.KvMap_tex)
  	if err != nil {
  		panic(err)
  	}
  	err = write_file("cover.tex", contents)
  }

  switch {
  case appl.Option == 1: // Everything w/ ref included as one file
    go appl.text_cover()
    cmdArgs = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_all_ref\".tex"}
  case appl.Option == 2: // Everything w/ ref as seperate files
    go appl.text_cover()
    cmdArgs_cover = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_cover\".tex"}
    cmdArgs_resume = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_resume\".tex"}
    cmdArgs_CV = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_CV_ref\".tex"}
  case appl.Option == 3: // Cover + CV w/ ref included as one file
    go appl.text_cover()
    cmdArgs = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_cover_CV_ref\".tex"}
  case appl.Option == 4: // Cover + CV w/ ref as seperate files
    go appl.text_cover()
    cmdArgs_cover = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_cover\".tex"}
    cmdArgs_CV = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_CV_ref\".tex"}
  case appl.Option == 5: // Cover + Resume included as one file
    go appl.text_cover()
    cmdArgs = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_cover_resume\".tex"}
  case appl.Option == 6: // Cover + Resume as seperate files
    go appl.text_cover()
    cmdArgs_cover = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_cover\".tex"}
    cmdArgs_resume = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_resume\".tex"}
  case appl.Option == 7: // CV w/ ref
    cmdArgs = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_CV_ref\".tex"}
  case appl.Option == 8: // CV w/0 ref
    cmdArgs = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_CV\".tex"}
  default: // just the resume
    cmdArgs = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_resume\".tex"}
  }

  if len(cmdArgs) > 0 {
    err = cmd.Run_cmd(TEX_CMD, cmdArgs)
    if err != nil {
      panic(err)
    }
  }

  if len(cmdArgs_cover) > 0 {
    err = cmd.Run_cmd(TEX_CMD, cmdArgs_cover)
    if err != nil {
      panic(err)
    }
  }

  if len(cmdArgs_resume) > 0 {
    err = cmd.Run_cmd(TEX_CMD, cmdArgs_resume)
    if err != nil {
      panic(err)
    }
  }

  if len(cmdArgs_CV) > 0 {
    err = cmd.Run_cmd(TEX_CMD, cmdArgs_CV)
    if err != nil {
      panic(err)
    }
  }

  err = rename_files(appl)
  if err != nil {
    panic(err)
  }

	return nil
}


/*
DESC: Renames PDF cover letters and resumes
IN: App object app
OUT: nill on success
*/
//TODO: functionize and move text_cover() calls
func rename_files(appl *App) error {

	var err error

	// Replace spaces in company name with _
	companyName := strings.Replace(appl.Company, " ", "_", -1)

  switch {
  case appl.Option == 1: // Everything w/ ref included as one file
    newName := ALL_FILENAME_BEGIN + companyName + "_" + date.Get_date("fileName") + ".pdf"
    err = copyfiles.Copy_file("main_all_ref.pdf", newName)
    if err != nil {
      panic(err)
    }
    appl.Attachments = append(appl.Attachments, newName)
  case appl.Option == 2: // Everything w/ ref as seperate files
    // COVER
    newName_cover := COVER_FILENAME_BEGIN + companyName + "_" + date.Get_date("fileName") + ".pdf"
    err = copyfiles.Copy_file("main_cover.pdf", newName_cover)
    if err != nil {
      panic(err)
    }
    appl.Attachments = append(appl.Attachments, newName_cover)
    // RESUME
    newName_resume := RESUME_FILENAME_BEGIN + companyName + "_" + date.Get_date("fileName") + ".pdf"
    err = copyfiles.Copy_file("main_resume.pdf", newName_resume)
    if err != nil {
      panic(err)
    }
    appl.Attachments = append(appl.Attachments, newName_resume)
    // CV w/ ref
    newName_CV := CV_FILENAME_BEGIN + companyName + "_" + date.Get_date("fileName") + ".pdf"
    err = copyfiles.Copy_file("main_CV_ref.pdf", newName_CV)
    if err != nil {
      panic(err)
    }
    appl.Attachments = append(appl.Attachments, newName_CV)
  case appl.Option == 3: // Cover + CV w/ ref included as one file
    newName := ALL_FILENAME_BEGIN + companyName + "_" + date.Get_date("fileName") + ".pdf"
    err = copyfiles.Copy_file("main_cover_CV_ref.pdf", newName)
    if err != nil {
      panic(err)
    }
    appl.Attachments = append(appl.Attachments, newName)
  case appl.Option == 4: // Cover + CV w/ ref as seperate files
    // COVER
    newName_cover := COVER_FILENAME_BEGIN + companyName + "_" + date.Get_date("fileName") + ".pdf"
    err = copyfiles.Copy_file("main_cover.pdf", newName_cover)
    if err != nil {
      panic(err)
    }
    appl.Attachments = append(appl.Attachments, newName_cover)
    // CV w/ ref
    newName_CV := CV_FILENAME_BEGIN + companyName + "_" + date.Get_date("fileName") + ".pdf"
    err = copyfiles.Copy_file("main_CV_ref.pdf", newName_CV)
    if err != nil {
      panic(err)
    }
    appl.Attachments = append(appl.Attachments, newName_CV)
  case appl.Option == 5: // Cover + Resume included as one file
    newName := ALL_FILENAME_BEGIN + companyName + "_" + date.Get_date("fileName") + ".pdf"
    err = copyfiles.Copy_file("main_cover_resume.pdf", newName)
    if err != nil {
      panic(err)
    }
    appl.Attachments = append(appl.Attachments, newName)
  case appl.Option == 6: // Cover + Resume as seperate files
    // COVER
    newName_cover := COVER_FILENAME_BEGIN + companyName + "_" + date.Get_date("fileName") + ".pdf"
    err = copyfiles.Copy_file("main_cover.pdf", newName_cover)
    if err != nil {
      panic(err)
    }
    appl.Attachments = append(appl.Attachments, newName_cover)
    // RESUME
    newName_resume := RESUME_FILENAME_BEGIN + companyName + "_" + date.Get_date("fileName") + ".pdf"
    err = copyfiles.Copy_file("main_resume.pdf", newName_resume)
    if err != nil {
      panic(err)
    }
    appl.Attachments = append(appl.Attachments, newName_resume)
  case appl.Option == 7: // CV w/ ref
    // CV w/ ref
    newName_CV := CV_FILENAME_BEGIN + companyName + "_" + date.Get_date("fileName") + ".pdf"
    err = copyfiles.Copy_file("main_CV_ref.pdf", newName_CV)
    if err != nil {
      panic(err)
    }
    appl.Attachments = append(appl.Attachments, newName_CV)
  case appl.Option == 8: // CV w/0 ref
    // CV w/o ref
    newName_CV := CV_FILENAME_BEGIN + companyName + "_" + date.Get_date("fileName") + ".pdf"
    err = copyfiles.Copy_file("main_CV.pdf", newName_CV)
    if err != nil {
      panic(err)
    }
    appl.Attachments = append(appl.Attachments, newName_CV)
  default: // just the resume
    // RESUME
    newName_resume := RESUME_FILENAME_BEGIN + companyName + "_" + date.Get_date("fileName") + ".pdf"
    err = copyfiles.Copy_file("main_resume.pdf", newName_resume)
    if err != nil {
      panic(err)
    }
    appl.Attachments = append(appl.Attachments, newName_resume)
  }

  // print list of generated files
  for i, _ := range(appl.Attachments){
    fmt.Printf("Generated file: %s\n", appl.Attachments[i])
  }

	return nil
}


/*
DESC: Reads a file into a buffer and replaces strings key in map with coresponding values
IN: input file: in_file, key value map kvMap
OUT: new string contents: contents, nill on success
*/
func Replace_strings(inFile string, kvMap map[string]string) (string, error) {

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
DESC: generates copyable plain text version of cover letter
IN: Application object app
OUT: nill on success
*/
func (appl App) text_cover() error {
	contents, err := Replace_strings(TEXT_COVER_TEMPL, appl.KvMap_text)
	if err != nil {
		panic(err)
	}
	err = write_file("cover.txt", contents)
	if err != nil {
		panic(err)
	}

	return nil
}
