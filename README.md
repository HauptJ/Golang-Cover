# Golang-Cover
Golang cover letter generator

[![Build Status](https://travis-ci.org/HauptJ/Golang-Cover.svg?branch=master)](https://travis-ci.org/HauptJ/Golang-Cover)

See the tutorial I wrote:
[Using Golang to Generate Custom Cover Letters](https://hauptj.com/2018/06/12/using-golang-to-generate-custom-cover-letters/)

**UPDATE:** I have now added sample template files.

## Usage:

#### Build:
1. [Download and install LaTeX](https://www.latex-project.org/get/)
2. [Download, extract and copy Moderncv Classic into root repo directory](https://www.sharelatex.com/templates/cv-or-resume/moderncv-classic)
3. Fill out the templated files
4. Replace your name in `app/app.go`
5. Replace your **Email address** and **SMTP server address** in `email/email.go`
6. Download the Gomail package: `go get gopkg.in/gomail.v2`
7. Build the binary: `go build main.go`

#### Flags:
- `--opt`: **[required]** option, see list below
- `--company`: **[required]** company name
- `--note`: **[optional]** note to add at bottom
- `--skill1`, `--skill2`: **[optional]** additional skills to list in cover
- `--contact`: **[optional]** contact name
- `--position`: **[REQUIRED w/o --head]** name of the position
- `--head`: **[optional]** override default heading w/ custom one
- `--source`: **[REQUIRED w/o --head]** source of advertisement of the position, eg. LinkedIn
- `--to`: **[optional]** mail to address
- `--pass`: **[REQUIRED w/ --email]** email account password
- `--subject`: **[optional]** override default Email subject w/ custom one
- `--url`: **[optional]** URL to position ad
- `--test`: **[optional]** run test build which will not to be logged

#### Options:
- `1`.) Everything w/ ref included as one file
- `2`.) Everything w/ ref as seperate files
- `3`.) Cover and CV w/ ref included as one file
- `4`.) Cover and CV w/ ref as seperate files
- `5`.) Cover and Resume included as one file
- `6`.) Cover and Resume as seperate files
- `7`.) CV w/ ref
- `8`.) CV w/0 ref
- `default`: Just the resume

**Example: ** ```.\main.exe --company "Some Company" --opt 1 --note "I am a developer who loves to code." --source "your company's website" --position "Software Engineer" --contact "Some Person"```
