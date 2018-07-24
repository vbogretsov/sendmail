package fixture

import (
	"encoding/json"

	"github.com/vbogretsov/go-validation"
	jsonerr "github.com/vbogretsov/go-validation/json"

	"github.com/vbogretsov/sendmail/model"

	"github.com/vbogretsov/sendmail/test/api/loader"
)

const (
	errorCannotBeBlank       = "cannot be blank"
	errorInvalidEmail        = "invalid email"
	errorMissingRecipients   = "missing recipients"
	errorInvalidBodyType     = "invalid body type"
	errrorCannotLoadTemplate = "cannot load template"
)

var (
	defaultArgs = map[string]interface{}{"Username": "user@mail.com"}
	bodyTypes   = []interface{}{"text/plain", "text/html"}
)

type Fixture struct {
	Name    string
	Request model.Request
	Errors  []JsonError
}

func marshal(err validation.Errors) json.Marshaler {
	return jsonerr.New(err, jsonerr.DefaultFormatter, jsonerr.DefaultJoiner)
}

type Params map[string]interface{}

type JsonError struct {
	Path   string `json:"path,omitempty"`
	Error  string `json:"error"`
	Params Params `json:"params,omitempty"`
}

var Fixtures = []Fixture{
	{
		Name: "ErrorIfMissingTemplateLang",
		Request: model.Request{
			TemplateName: loader.TemplateValid,
			TemplateArgs: defaultArgs,
			To: []model.Address{
				{
					Email: "to1@mail.com",
				},
			},
		},
		Errors: []JsonError{
			{
				Path:  ".templateLang",
				Error: errorCannotBeBlank,
			},
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
				},
			},
		},
		Errors: []JsonError{
			{
				Path:  ".templateName",
				Error: errorCannotBeBlank,
			},
		},
	},
	{
		Name: "ErrorIfMissingRecipients",
		Request: model.Request{
			TemplateLang: loader.Lang,
			TemplateName: loader.TemplateValid,
			TemplateArgs: defaultArgs,
		},
		Errors: []JsonError{
			{
				Path:  ".",
				Error: errorMissingRecipients,
			},
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
				},
			},
		},
		Errors: []JsonError{
			{
				Path:  ".to[0].Email",
				Error: errorInvalidEmail,
			},
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
				},
			},
		},
		Errors: []JsonError{
			{
				Path:  ".to[0].Email",
				Error: errorInvalidEmail,
			},
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
				},
			},
		},
		Errors: []JsonError{
			{
				Error: errrorCannotLoadTemplate,
				Params: Params{
					"lang":  loader.Lang,
					"name":  "xxx",
					"cause": "template not found",
				},
			},
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
					Email: "to1@mail.com",
				},
			},
		},
		Errors: []JsonError{
			{
				Path:  ".BodyType",
				Error: errorCannotBeBlank,
			},
			{
				Path:  ".BodyType",
				Error: "invalid body type",
				Params: Params{
					"unsupported": "",
					"supported":   bodyTypes,
				},
			},
		},
	},
	{
		Name: "ErrorIfInvalidBodyType",
		Request: model.Request{
			TemplateLang: loader.Lang,
			TemplateName: loader.TemplateInvalidBodyType,
			TemplateArgs: defaultArgs,
			To: []model.Address{
				{
					Email: "to1@mail.com",
				},
			},
		},
		Errors: []JsonError{
			{
				Path:  ".BodyType",
				Error: errorInvalidBodyType,
				Params: Params{
					"unsupported": "text/xxx",
					"supported":   bodyTypes,
				},
			},
		},
	},
	{
		Name: "ErrorIfMissingBody",
		Request: model.Request{
			TemplateLang: loader.Lang,
			TemplateName: loader.TemplateMissingBody,
			TemplateArgs: defaultArgs,
			To: []model.Address{
				{
					Email: "to1@mail.com",
				},
			},
		},
		Errors: []JsonError{
			{
				Path:  ".Body",
				Error: errorCannotBeBlank,
			},
		},
	},
	{
		Name: "ErrorIfMissingSubject",
		Request: model.Request{
			TemplateLang: loader.Lang,
			TemplateName: loader.TemplateMissingSubject,
			TemplateArgs: defaultArgs,
			To: []model.Address{
				{
					Email: "to1@mail.com",
				},
			},
		},
		Errors: []JsonError{
			{
				Path:  ".Subject",
				Error: errorCannotBeBlank,
			},
		},
	},
	{
		Name: "ErrorIfMissingFrom",
		Request: model.Request{
			TemplateLang: loader.Lang,
			TemplateName: loader.TemplateMissingFrom,
			TemplateArgs: defaultArgs,
			To: []model.Address{
				{
					Email: "to1@mail.com",
				},
			},
		},
		Errors: []JsonError{
			{
				Path:  ".From.Email",
				Error: errorInvalidEmail,
			},
		},
	},
	{
		Name: "ErrorIfMissingFromEmail",
		Request: model.Request{
			TemplateLang: loader.Lang,
			TemplateName: loader.TemplateMissingFromEmail,
			TemplateArgs: defaultArgs,
			To: []model.Address{
				{
					Email: "to1@mail.com",
				},
			},
		},
		Errors: []JsonError{
			{
				Path:  ".From.Email",
				Error: errorInvalidEmail,
			},
		},
	},
	{
		Name: "ErrorIfInvalidFromEmail",
		Request: model.Request{
			TemplateLang: loader.Lang,
			TemplateName: loader.TemplateInvalidFromEmail,
			TemplateArgs: defaultArgs,
			To: []model.Address{
				{
					Email: "to1@mail.com",
				},
			},
		},
		Errors: []JsonError{
			{
				Path:  ".From.Email",
				Error: errorInvalidEmail,
			},
		},
	},
}
