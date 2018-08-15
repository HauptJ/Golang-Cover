# Golang-Cover
Golang cover letter generator

See the tutorial I wrote:
[Using Golang to Generate Custom Cover Letters](https://hauptj.com/2018/06/12/using-golang-to-generate-custom-cover-letters/)

**UPDATE:** I have now added sample template files.

## Usage:

#### Flags:
- --cover: cover is included **[incl]**, seperate **[sep]** or not generated at all **[no]**
- --company: company name
- --note: note to add at bottom
- --skill1, --skill2: additional skills to list in cover
- --contact: contact name
- --position: name of the position
- --source: source of advertisement of the position, eg. LinkedIn


**Example: ** ```.\main.exe --company "Some Company" --cover sep --note "I am a developer who loves to code." --source "your company's website" --position "Software Engineer" --contact "Some Person"```
