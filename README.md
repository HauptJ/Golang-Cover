# Golang-Cover
Golang cover letter generator

See the tutorial I wrote:
[Using Golang to Generate Custom Cover Letters](https://hauptj.com/2018/06/12/using-golang-to-generate-custom-cover-letters/)

**UPDATE:** I have now added sample template files.

## Usage:

#### Build:
`go build main.go`

#### Flags:
- --opt: option, see list below
- --company: company name
- --note: note to add at bottom
- --skill1, --skill2: additional skills to list in cover
- --contact: contact name
- --position: name of the position
- --source: source of advertisement of the position, eg. LinkedIn

#### Options:
- 1.) Everything w/ ref included as one file
- 2.) Everything w/ ref as seperate files
- 3.) Cover and CV w/ ref included as one file
- 4.) Cover and CV w/ ref as seperate files
- 5.) Cover and Resume included as one file
- 6.) Cover and Resume as seperate files
- 7.) CV w/ ref
- 8.) CV w/0 ref
- default: Just the resume

**Example: ** ```.\main.exe --company "Some Company" --opt 1 --note "I am a developer who loves to code." --source "your company's website" --position "Software Engineer" --contact "Some Person"```
