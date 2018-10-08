# Golang-Cover
Golang application tracker, cover letter and follow up generator, and whatever else I can think of adding.

This repository contains a mirror of the Golang source files used in a private repository that I use for my resume and CV.

See the tutorial I wrote:
[Using Golang to Generate Custom Cover Letters](https://hauptj.com/2018/06/12/using-golang-to-generate-custom-cover-letters/)

[![Build Status](https://travis-ci.org/HauptJ/Golang-Cover.svg?branch=master)](https://travis-ci.org/HauptJ/Golang-Cover)

## Usage:

#### Build:
1. [Download and install LaTeX](https://www.latex-project.org/get/)
2. [Download, extract and copy Moderncv Classic into root repo directory](https://www.sharelatex.com/templates/cv-or-resume/moderncv-classic)
3. Fill out the templated files
4. Replace your name in `app/app.go`
5. Replace your **Email address** and **SMTP server address** in `email/email.go`
6. Download the following Go packages:
    1. `go get gopkg.in/gomail.v2`
    2. `go get cloud.google.com/go/storage`
    3. `go get golang.org/x/net/context`
7. Build the binary: `go build main.go`

#### ENV Vars:
- ##### Windows:
  - **Email:**
    - `$env:MailFrom="user@example.com"`
    - `$env:MailPass="password"`
    - `$env:EmailSMTP="smtp.gmail.com"` Email SMTP server
  - **GCP Storage:**
    - `$env:GCProjectID="project-123"`
    - `$env:GCBucket="bucket.com"`


- ##### Linux:
  - **Email:**
    - `export MailFrom="user@example.com"` Email account
    - `export MailPass="password"` Email account password
    - `export EmailSMTP="smtp.gmail.com"` Email SMTP server
  - **GCP Storage:**
    - `export GCProjectID="project-123"` the GCP bucket to upload content to
    - `export GCBucket="bucket.com"` the ID of the GCP project to use

#### GCP Default Application Auth:
  `gcloud auth application-default login`

#### Flags:
- `--opt`: **[required]** option, see list below
- `--company`: **[required]** company name
- `--note1`, `--note2`: **[optional]** notes to add at bottom
- `--skill1`, `--skill2`, `--skill3`: **[optional]** additional skills to list in cover
- `--contact`: **[optional]** contact name
- `--position`: **[REQUIRED w/o --head]** name of the position
- `--positionID`: **[optional]** position ID
- `--applied`: **[optional]** when application was submitted
- `--head`: **[optional]** override default heading w/ custom one
- `--headAdd`: **[optional]** Extend the default header
- `--source`: **[REQUIRED w/o --head]** source of advertisement of the position, eg. LinkedIn
- `--to`: **[optional]** mail to address
- `--subject`: **[optional]** override default Email subject w/ custom one
- `--url`: **[optional]** URL to position ad
- `--test`: **[optional]** run test build which will not to be logged
- `--upload`: **[optional]** upload file to Google Cloud storage bucket **NOTE** only applies to opt `7` and `8`

#### Options:
- `1`.) Everything w/ ref included as one file
- `2`.) Everything w/ ref as separate files
- `3`.) Cover and CV w/ ref included as one file
- `4`.) Cover and CV w/ ref as separate files
- `5`.) Cover and Resume included as one file
- `6`.) Cover and Resume as separate files
- `7`.) CV w/ ref
- `8`.) CV w/0 ref
- `9`.) Just the resume
- `10`.) Follow Up message


**Example: ** ```.\main.exe --company "Some Company" --opt 1 --note "I am a developer who loves to code." --source "your company's website" --position "Software Engineer" --contact "Some Person"```
