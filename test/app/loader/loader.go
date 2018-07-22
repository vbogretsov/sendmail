package loader

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/vbogretsov/sendmail/app"
)

const ExpectedBody = `Hello %s!
This is test body!
`

const Lang = "en"

const (
	TemplateValid            = "valid"
	TemplateMissingBodyType  = "invalid-missing-body-type"
	TemplateInvalidBodyType  = "invalid-body-type"
	TemplateMissingBody      = "invalid-missing-body"
	TemplateMissingSubject   = "invalid-missing-subject"
	TemplateMissingFrom      = "invalid-missing-from"
	TemplateMissingFromEmail = "invalid-missing-from-email"
	TemplateInvalidFromEmail = "invalid-from-email"
)

const msg = `
From:
  Email: user@mail.com
  Name: Sender
Subject: Subject
BodyType: "text/plain"
Body: |
  Hello {{.Username}}!
  This is test body!
`

const msgInvalidBodyType = `
From:
  Email: user@mail.com
  Name: LevelUp
Subject: Subject
BodyType: "text/xxx"
Body: |
  Hello {{.Username}}!
  This is test body!
`

const msgMissingBodyType = `
From:
  Email: user@mail.com
  Name: LevelUp
Subject: Subject
Body: |
  Hello {{.Username}}!
  This is test body!
`

const msgMissingBody = `
From:
  Email: user@mail.com
  Name: LevelUp
Subject: Subject
BodyType: "text/plain"
`

const msgMissingSubject = `
From:
  Email: user@mail.com
  Name: LevelUp
BodyType: "text/plain"
Body: |
  Hello {{.Username}}!
  This is test body!
`

const msgMissingFrom = `
Subject: Subject
BodyType: "text/plain"
Body: |
  Hello {{.Username}}!
  This is test body!
`

const msgMissingFromEmail = `
From:
  Name: LevelUp
Subject: Subject
BodyType: "text/plain"
Body: |
  Hello {{.Username}}!
  This is test body!
`

const msgInvalidFromEmail = `
From:
  Email: user.mail.com
  Name: LevelUp
Subject: Subject
BodyType: "text/plain"
Body: |
  Hello {{.Username}}!
  This is test body!
`

type loader struct {
	templates map[string]string
}

func New() app.Loader {
	return &loader{templates: map[string]string{
		fmt.Sprintf("%s-%s", Lang, TemplateValid):            msg,
		fmt.Sprintf("%s-%s", Lang, TemplateMissingBodyType):  msgMissingBodyType,
		fmt.Sprintf("%s-%s", Lang, TemplateInvalidBodyType):  msgInvalidBodyType,
		fmt.Sprintf("%s-%s", Lang, TemplateMissingBody):      msgMissingBody,
		fmt.Sprintf("%s-%s", Lang, TemplateMissingSubject):   msgMissingSubject,
		fmt.Sprintf("%s-%s", Lang, TemplateMissingFrom):      msgMissingFrom,
		fmt.Sprintf("%s-%s", Lang, TemplateMissingFromEmail): msgMissingFromEmail,
		fmt.Sprintf("%s-%s", Lang, TemplateInvalidFromEmail): msgInvalidFromEmail,
	}}
}

func (self *loader) Load(lang, name string) (io.Reader, error) {
	id := fmt.Sprintf("%s-%s", lang, name)
	text, ok := self.templates[id]
	if !ok {
		return nil, errors.New("template not found")
	}
	return strings.NewReader(text), nil
}
