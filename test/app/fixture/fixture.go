package fixture

import (
	"github.com/vbogretsov/go-validation"

	"github.com/vbogretsov/sendmail/app"
	"github.com/vbogretsov/sendmail/model"

	"github.com/vbogretsov/sendmail/test/app/loader"
)

const (
	errorStrCannotBeBlank  = "cannot be blank"
	errorStrInvalidEmail   = "invalid email"
	errorMissingRecipients = "missing recipients"
	errorInvalidBodyType   = "invalid body type"
)

var (
	defaultArgs = map[string]interface{}{"Username": "user@mail.com"}
	bodyTypes   = []interface{}{"text/plain", "text/html"}
)

var (
	errTemplateLangBlank = validation.StructError{
		Field: "templateLang",
		Errors: []error{
			validation.Error{Message: errorStrCannotBeBlank},
		},
	}
	errTemplateNameBlank = validation.StructError{
		Field: "templateName",
		Errors: []error{
			validation.Error{Message: errorStrCannotBeBlank},
		},
	}
	errMissingRecipients = validation.StructError{
		Field: "",
		Errors: []error{
			validation.Error{Message: errorMissingRecipients},
		},
	}
	errInvalidEmail = validation.StructError{
		Field: "Email",
		Errors: []error{
			validation.Error{
				Message: errorStrInvalidEmail,
			},
		},
	}
	errTemplateNotFound = validation.Error{
		Message: "cannot load template",
		Params: validation.Params{
			"lang":  loader.Lang,
			"name":  "xxx",
			"cause": "template not found",
		},
	}
	errBodyBlank = validation.StructError{
		Field: "Body",
		Errors: []error{validation.Error{
			Message: errorStrCannotBeBlank,
		}},
	}
	errSubjectBlank = validation.StructError{
		Field: "Subject",
		Errors: []error{validation.Error{
			Message: errorStrCannotBeBlank,
		}},
	}
)

type Fixture struct {
	Name    string
	Request model.Request
	Result  error
}

var Fixtures = []Fixture{
	{
		Name: "ErrorIfMissingTemplateLang",
		Request: model.Request{
			TemplateName: loader.TemplateValid,
			TemplateArgs: defaultArgs,
			To: []model.Address{
				{
					Email: "to1@mail.com", Name: "",
				},
			},
		},
		Result: app.ArgumentError{
			Err: validation.Errors([]error{errTemplateLangBlank}),
		},
	},
	{
		Name: "ErrorIfMissingTemplateName",
		Request: model.Request{
			TemplateLang: loader.Lang,
			TemplateArgs: defaultArgs,
			To: []model.Address{
				{
					Email: "to1@mail.com",
					Name:  "",
				},
			},
		},
		Result: app.ArgumentError{
			Err: validation.Errors([]error{errTemplateNameBlank}),
		},
	},
	{
		Name: "ErrorIfMissingRecipients",
		Request: model.Request{
			TemplateLang: loader.Lang,
			TemplateName: loader.TemplateValid,
			TemplateArgs: defaultArgs,
		},
		Result: app.ArgumentError{
			Err: validation.Errors([]error{errMissingRecipients}),
		},
	},
	{
		Name: "ErrorIfRecipientMissingEmail",
		Request: model.Request{
			TemplateLang: loader.Lang,
			TemplateName: loader.TemplateValid,
			TemplateArgs: defaultArgs,
			To: []model.Address{
				{
					Email: "",
					Name:  "",
				},
			},
		},
		Result: app.ArgumentError{
			Err: validation.Errors([]error{
				validation.StructError{
					Field: "to",
					Errors: []error{
						validation.SliceError{
							Index:  0,
							Errors: []error{errInvalidEmail},
						},
					},
				},
			}),
		},
	},
	{
		Name: "ErrorIfRecipientEmailInvalid",
		Request: model.Request{
			TemplateLang: loader.Lang,
			TemplateName: loader.TemplateValid,
			TemplateArgs: defaultArgs,
			To: []model.Address{
				{
					Email: "to1.mail.com",
					Name:  "",
				},
			},
		},
		Result: app.ArgumentError{
			Err: validation.Errors([]error{
				validation.StructError{
					Field: "to",
					Errors: []error{
						validation.SliceError{
							Index:  0,
							Errors: []error{errInvalidEmail},
						},
					},
				},
			}),
		},
	},
	{
		Name: "ErrorIfTemplateNotFound",
		Request: model.Request{
			TemplateLang: loader.Lang,
			TemplateName: "xxx",
			TemplateArgs: defaultArgs,
			To: []model.Address{
				{
					Email: "to1@mail.com",
					Name:  "",
				},
			},
		},
		Result: app.ArgumentError{
			Err: validation.Errors([]error{errTemplateNotFound}),
		},
	},
	{
		Name: "ErrorIfMissingBodyType",
		Request: model.Request{
			TemplateLang: loader.Lang,
			TemplateName: loader.TemplateMissingBodyType,
			TemplateArgs: defaultArgs,
			To: []model.Address{
				{
					Email: "to1@mail.com", Name: "",
				},
			},
		},
		Result: validation.Errors([]error{
			validation.StructError{
				Field: "BodyType",
				Errors: []error{
					validation.Error{
						Message: errorStrCannotBeBlank,
					},
					validation.Error{
						Message: errorInvalidBodyType,
						Params: validation.Params{
							"unsupported": "",
							"supported":   bodyTypes,
						},
					},
				},
			},
		}),
	},
	{
		Name: "ErrorIfInvalidBodyType",
		Request: model.Request{
			TemplateLang: loader.Lang,
			TemplateName: loader.TemplateInvalidBodyType,
			TemplateArgs: defaultArgs,
			To: []model.Address{
				{
					Email: "to1@mail.com", Name: "",
				},
			},
		},
		Result: validation.Errors([]error{
			validation.StructError{
				Field: "BodyType",
				Errors: []error{
					validation.Error{
						Message: errorInvalidBodyType,
						Params: validation.Params{
							"unsupported": "text/xxx",
							"supported":   bodyTypes,
						},
					},
				},
			},
		}),
	},
	{
		Name: "ErrorIfMissingBody",
		Request: model.Request{
			TemplateLang: loader.Lang,
			TemplateName: loader.TemplateMissingBody,
			TemplateArgs: defaultArgs,
			To: []model.Address{
				{
					Email: "to1@mail.com", Name: "",
				},
			},
		},
		Result: validation.Errors([]error{errBodyBlank}),
	},
	{
		Name: "ErrorIfMissingSubject",
		Request: model.Request{
			TemplateLang: loader.Lang,
			TemplateName: loader.TemplateMissingSubject,
			TemplateArgs: defaultArgs,
			To: []model.Address{
				{
					Email: "to1@mail.com", Name: "",
				},
			},
		},
		Result: validation.Errors([]error{errSubjectBlank}),
	},
	{
		Name: "ErrorIfMissingFrom",
		Request: model.Request{
			TemplateLang: loader.Lang,
			TemplateName: loader.TemplateMissingFrom,
			TemplateArgs: defaultArgs,
			To: []model.Address{
				{
					Email: "to1@mail.com", Name: "",
				},
			},
		},
		Result: validation.Errors([]error{
			validation.StructError{
				Field:  "From",
				Errors: []error{errInvalidEmail},
			},
		}),
	},
	{
		Name: "ErrorIfMissingFromEmail",
		Request: model.Request{
			TemplateLang: loader.Lang,
			TemplateName: loader.TemplateMissingFromEmail,
			TemplateArgs: defaultArgs,
			To: []model.Address{
				{
					Email: "to1@mail.com", Name: "",
				},
			},
		},
		Result: validation.Errors([]error{
			validation.StructError{
				Field:  "From",
				Errors: []error{errInvalidEmail},
			},
		}),
	},
	{
		Name: "ErrorIfInvalidFromEmail",
		Request: model.Request{
			TemplateLang: loader.Lang,
			TemplateName: loader.TemplateInvalidFromEmail,
			TemplateArgs: defaultArgs,
			To: []model.Address{
				{
					Email: "to1@mail.com", Name: "",
				},
			},
		},
		Result: validation.Errors([]error{
			validation.StructError{
				Field:  "From",
				Errors: []error{errInvalidEmail},
			},
		}),
	},
}
