package app

import (
	"github.com/vbogretsov/go-validation"
	"github.com/vbogretsov/go-validation/rule"

	"github.com/vbogretsov/sendmail/model"
)

const (
	errorStrRequired       = "cannot be blank"
	errorMissingRecipients = "missing recipients"
	errorInvalidEmail      = "invalid email"
	errorInvalidBodyType   = "invalid body type"
)

var (
	strRequired = rule.StrRequired(errorStrRequired)
	strEmail    = rule.StrEmail(errorInvalidEmail)
	bodyTypes   = rule.In([]interface{}{"text/plain", "text/html"}, errorInvalidBodyType)
)

func self(v interface{}) interface{} {
	return v
}

func requestTemplateLang(v interface{}) interface{} {
	return &((v.(*model.Request)).TemplateLang)
}

func requestTemplateName(v interface{}) interface{} {
	return &((v.(*model.Request)).TemplateName)
}

func requestTo(v interface{}) interface{} {
	return &((v.(*model.Request)).To)
}

func requestCc(v interface{}) interface{} {
	return &((v.(*model.Request)).Cc)
}

func requestBcc(v interface{}) interface{} {
	return &((v.(*model.Request)).Bcc)
}

func addressEmail(v interface{}) interface{} {
	return &(v.(*model.Address)).Email
}

func addressIter(v interface{}, i int) interface{} {
	return &((*v.(*[]model.Address))[i])
}

func messageSubject(v interface{}) interface{} {
	return &(v.(*model.Message)).Subject
}

func messageFrom(v interface{}) interface{} {
	return &(v.(*model.Message)).From
}

func messageBodyType(v interface{}) interface{} {
	return &(v.(*model.Message)).BodyType
}

func messageBody(v interface{}) interface{} {
	return &(v.(*model.Message)).Body
}

func atLeastOneRecipient(v interface{}) error {
	req := v.(*model.Request)
	if len(req.To) == 0 && len(req.Cc) == 0 && len(req.Bcc) == 0 {
		return validation.Error{Message: errorMissingRecipients}
	}
	return nil
}

var addressRule = validation.Struct(&model.Address{}, "yaml", []validation.Field{
	{
		Attr:  addressEmail,
		Rules: []validation.Rule{strEmail},
	},
})

var requestRule = validation.Struct(&model.Request{}, "json", []validation.Field{
	{
		Attr:  requestTemplateLang,
		Rules: []validation.Rule{strRequired},
	},
	{
		Attr:  requestTemplateName,
		Rules: []validation.Rule{strRequired},
	},
	{
		Attr: requestTo,
		Rules: []validation.Rule{
			rule.SliceEach(addressIter, []validation.Rule{addressRule}),
		},
	},
	{
		Attr: requestCc,
		Rules: []validation.Rule{
			rule.SliceEach(addressIter, []validation.Rule{addressRule}),
		},
	},
	{
		Attr: requestBcc,
		Rules: []validation.Rule{
			rule.SliceEach(addressIter, []validation.Rule{addressRule}),
		},
	},
	{
		Attr:  self,
		Rules: []validation.Rule{validation.Func(atLeastOneRecipient)},
	},
})(nil)

var messageRule = validation.Struct(&model.Message{}, "yaml", []validation.Field{
	{
		Attr:  messageSubject,
		Rules: []validation.Rule{strRequired},
	},
	{
		Attr:  messageBodyType,
		Rules: []validation.Rule{strRequired, bodyTypes},
	},
	{
		Attr:  messageBody,
		Rules: []validation.Rule{strRequired},
	},
	{
		Attr:  messageFrom,
		Rules: []validation.Rule{addressRule},
	},
	{
		Attr:  self,
		Rules: []validation.Rule{},
	},
})(nil)
