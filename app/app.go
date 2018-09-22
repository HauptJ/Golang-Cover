/*
DESC: Generates Resume, Cover Letter and Follow Up
Author: Joshua Haupt
Last Modified: 09-22-2018
*/


package app

import (
  "io/ioutil"
	"strings"
  "fmt"
  "../date"
  "../copyfiles"
  "../cmd"
  "../gcloud"
)


/*
 Constant Declarations
*/
const TEX_COVER_TEMPL = "cover_template.tex"
const TEXT_COVER_TEMPL = "text_cover_template.txt"
const TEXT_FOLLOW_UP_TEMPL = "text_follow_up_template.txt"
const TEX_CMD = "pdflatex"
const ALL_FILENAME_BEGIN = "Joshua_Haupt_Cover_Resume_CV_"
const RESUME_FILENAME_BEGIN = "Joshua_Haupt_Resume_"
const COVER_FILENAME_BEGIN = "Joshua_Haupt_Cover_"
const CV_FILENAME_BEGIN = "Joshua_Haupt_CV_"
const DEFAULT_CONTACT = "To whom it may concern"
const CONTACT_ENDING = " or to whom it may concern"
const LOCAL = "I am currently located in the St. Louis area, and I am also receptive to relocation."
const DISTANT = "I am currently located in the St. Louis area, however, I am receptive to relocation."
const GCLOUD_FILENAME = "cv.pdf"


type TexCover struct {
  	Note1_tex string
    Note2_tex string
    Skill1_tex string
    Skill2_tex string
    Skill3_tex string
    KvMap_tex map[string]string
}

type TextCover struct {
  	Note1_text string
    Note2_text string
    Skill1_text string
    Skill2_text string
    Skill3_text string
    KvMap_text map[string]string
}

type EmailCover struct {
  	Note1_email string
    Note2_email string
    ReloLine_email string
    Skill1_email string
    Skill2_email string
    Skill3_email string
    KvMap_email map[string]string
}

type Email struct {
    MailTo string
    MailFrom string
    Subject string
    EmailPass string
}

type GCS struct {
    GCUploadFile bool
    GCBucket string
    GCProjectID string
}

type FollowUp struct {
    WhenApplied string
    FollowUpRef string
    FollowUpRefInfo string
}

type Common struct {
  Heading string
  HeadingAdd string
	Company string
	Contact string
	Position string
  PositionID string
	Source string
	Note1 string
  Note2 string
	Local bool
	ReloLine string
	Skill1 string
	Skill2 string
  Skill3 string
	Url string
  Attachments []string
}

type Control struct {
  Option int
  Test bool
}


type App struct {
  TexCover
  TextCover
  EmailCover
  Email
  GCS   // Google Cloud Storage
  FollowUp
  Common
  Control
}


/*
DESC: parses flag string values to generate App object values
IN: the flag values as strings and the App object
OUT: nil on success
*/
func PharseFlags(appl *App) error {

	if appl.Contact == "" {
		appl.Contact = DEFAULT_CONTACT
	} else {
		appl.Contact = appl.Contact + CONTACT_ENDING
	}

	if appl.Heading == "" && appl.Position != "" && appl.Source != "" {
		appl.Heading = "I am excited about the possibility of joining your organization in the position of " + appl.Position + ", as advertised on " + appl.Source + ". " + appl.HeadingAdd// default heading
	} else {
		if appl.Option <= 6 && appl.Option > 0 {
			panic("heading undefined")
		}
	}

  if appl.Option == 10 && appl.Position != "" && appl.PositionID != "" {
    appl.FollowUpRef = "For reference, I have pasted the position title and ID number below."
    appl.FollowUpRefInfo = appl.Position + " - " + appl.PositionID
  }

  if appl.Option == 10 && appl.Position != "" && appl.PositionID == "" {
    appl.FollowUpRef = "For reference, I have pasted the position title below."
    appl.FollowUpRefInfo = appl.Position
  }

  if appl.Option == 10 && (appl.Position == "" || appl.Company == "") {
    panic("Not all of the required flags for a follow up were provided. Need: --position and --company")
  }

	// If additional note is present, add a newline at the end of it
	if appl.Note1 != "" {
		appl.Note1_tex = appl.Note1 + " \\newline"
		appl.Note1_email = "<div><p style=\"text-align:left\";>" + appl.Note1 + "</p></div>"
	}

  // If additional note is present, add a newline at the end of it
  if appl.Note2 != "" {
    appl.Note2_tex = appl.Note2 + " \\newline"
    appl.Note2_email = "<div><p style=\"text-align:left\";>" + appl.Note2 + "</p></div>"
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

  // If additional skill is present, add a \item before it
  if appl.Skill3 != "" {
    appl.Skill3_tex = "\\item " + appl.Skill3
    appl.Skill3_email = "<li>" + appl.Skill3 + "</li>"
    appl.Skill3_text = "- " + appl.Skill3
  }

  fmt.Println(appl.Local)

	if appl.Local == true {
		appl.ReloLine = LOCAL
  }
	if appl.Local == false {
		appl.ReloLine = DISTANT
  }

	appl.ReloLine_email = "<div><p style=\"text-align:left;\">" + appl.ReloLine + "</p></div>"


	if (appl.Option <= 6 && appl.Option > 0) || appl.Option == 10 {
		appl.KvMap_tex = map[string]string{"[COMPANY_NAME]": strings.Replace(appl.Company, "&", "\\&", -1), "[COMPANY_CONTACT]": strings.Replace(appl.Contact, "&", "\\&", -1),
			"[POSITION_NAME]": strings.Replace(appl.Position, "&", "\\&", -1), "[HEADING]": strings.Replace(appl.Heading, "&", "\\&", -1), "[POSITION_SOURCE]": strings.Replace(appl.Source, "&", "\\&", -1),
			"[ADDITIONAL_SKILL_1]": strings.Replace(appl.Skill1_tex, "&", "\\&", -1), "[ADDITIONAL_SKILL_2]": strings.Replace(appl.Skill2_tex, "&", "\\&", -1), "[ADDITIONAL_SKILL_3]": strings.Replace(appl.Skill3_tex, "&", "\\&", -1),
			"[ADDITIONAL_NOTE1]": strings.Replace(appl.Note1_tex, "&", "\\&", -1),	"[ADDITIONAL_NOTE2]": strings.Replace(appl.Note2_tex, "&", "\\&", -1), "[RELOCATION]": strings.Replace(appl.ReloLine, "&", "\\&", -1)}

		appl.KvMap_text = map[string]string{"[COMPANY_NAME]": appl.Company, "[COMPANY_CONTACT]": appl.Contact, "[POSITION_NAME]": appl.Position,
			"[HEADING]": appl.Heading, "[POSITION_SOURCE]": appl.Source, "[ADDITIONAL_SKILL_1]": appl.Skill1_text, "[ADDITIONAL_SKILL_2]": appl.Skill2_text,  "[ADDITIONAL_SKILL_3]": appl.Skill3_text,
			"[ADDITIONAL_NOTE1]": appl.Note1, "[ADDITIONAL_NOTE2]": appl.Note2, "[CURRENT_DATE]": date.Get_date("email"), "[RELOCATION]": appl.ReloLine, "[WHEN_APPLIED]": appl.WhenApplied, "[FOLLOWUPREF]": appl.FollowUpRef, "[FOLLOWUPREFINFO]": appl.FollowUpRefInfo}
	}

	if appl.MailTo != "" && appl.EmailPass != ""  && appl.MailFrom != "" {
		appl.KvMap_email = map[string]string{"[COMPANY_NAME]": appl.Company, "[COMPANY_CONTACT]": appl.Contact, "[POSITION_NAME]": appl.Position,
			"[HEADING]": appl.Heading, "[POSITION_SOURCE]": appl.Source, "[ADDITIONAL_SKILL_1]": appl.Skill1_email, "[ADDITIONAL_SKILL_2]": appl.Skill2_email, "[ADDITIONAL_SKILL_3]": appl.Skill3_email,
			"[ADDITIONAL_NOTE1]": appl.Note1_email, "[ADDITIONAL_NOTE2]": appl.Note2_email, "[CURRENT_DATE]": date.Get_date("email"), "[RELOCATION]": appl.ReloLine_email, "[WHEN_APPLIED]": appl.WhenApplied, "[FOLLOWUPREF]": appl.FollowUpRef, "[FOLLOWUPREFINFO]": appl.FollowUpRefInfo}


		if appl.Subject == "" && appl.Position != "" {
			appl.Subject = "Joshua Haupt appllication for " + appl.Position + " position at " + appl.Company // default subject
		}

	}

	return nil
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
  case appl.Option == 9: // just the resume
    cmdArgs = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_resume\".tex"}
  default: // the follow up
    go appl.text_follow_up()
    cmdArgs = []string{"-synctex=1", "-interaction=nonstopmode", "\"main_CV_ref\".tex"}
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
    if appl.GCUploadFile == true {
      err = gcloud.GCUpload(appl.GCProjectID, appl.GCBucket, newName_CV, GCLOUD_FILENAME, true) // upload to Google Cloud Storage
      if err != nil {
        panic(err)
      }
    }
  case appl.Option == 8: // CV w/0 ref
    // CV w/o ref
    newName_CV := CV_FILENAME_BEGIN + companyName + "_" + date.Get_date("fileName") + ".pdf"
    err = copyfiles.Copy_file("main_CV.pdf", newName_CV)
    if err != nil {
      panic(err)
    }
    appl.Attachments = append(appl.Attachments, newName_CV)
    if appl.GCUploadFile == true {
      err = gcloud.GCUpload(appl.GCProjectID, appl.GCBucket, newName_CV, GCLOUD_FILENAME, true) // upload to Google Cloud Storage
      if err != nil {
        panic(err)
      }
    }
  case appl.Option == 9: // just the resume
    // RESUME
    newName_resume := RESUME_FILENAME_BEGIN + companyName + "_" + date.Get_date("fileName") + ".pdf"
    err = copyfiles.Copy_file("main_resume.pdf", newName_resume)
    if err != nil {
      panic(err)
    }
    appl.Attachments = append(appl.Attachments, newName_resume)
  default: // follow up
    // CV w/ ref
    newName_CV := CV_FILENAME_BEGIN + companyName + "_" + date.Get_date("fileName") + ".pdf"
    err = copyfiles.Copy_file("main_CV_ref.pdf", newName_CV)
    if err != nil {
      panic(err)
    }
    appl.Attachments = append(appl.Attachments, newName_CV)
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


/*
DESC: generates copyable plain text version of follow up message
IN: Application object app
OUT: nill on success
*/
func (appl App) text_follow_up() error {
	contents, err := Replace_strings(TEXT_FOLLOW_UP_TEMPL, appl.KvMap_text)
	if err != nil {
		panic(err)
	}
	err = write_file("follow_up.txt", contents)
	if err != nil {
		panic(err)
	}

	return nil
}
