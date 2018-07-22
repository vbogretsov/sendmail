package app

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"text/template"

	"github.com/vbogretsov/go-validation"
	"gopkg.in/yaml.v2"

	"github.com/vbogretsov/sendmail/model"
)

// ErrLoadTemplate defines error message if template not found.
var ErrLoadTemplate = "cannot load template"

var (
	// ParamTemplateLang defines key name for error parameter template lang.
	ParamTemplateLang = "lang"
	// ParamTemplateName defines key name for error parameter template name.
	ParamTemplateName = "name"
	// ParamTemplateCause defines key name for  error caused template load fail.
	ParamTemplateCause = "cause"
)

// ArgumentError represents an error caused by user input.
type ArgumentError struct {
	Err error
}

// Error returns string representation of an argument error.
func (e ArgumentError) Error() string {
	return e.Err.Error()
}

// Loader represents interface of templates loader.
type Loader interface {
	// Load loads a template with the language and name provided.
	Load(lang, name string) (io.Reader, error)
}

// Sender represent inteface for an email sender.
type Sender interface {
	Send(model.Message) error
}

// App represents a maild application.
type App struct {
	loader Loader
	sender Sender
}

// New creates a new mail app.
func New(loader Loader, sender Sender) *App {
	return &App{
		loader: loader,
		sender: sender,
	}
}

// SendMail build email from template and sends it.
func (ap *App) SendMail(req model.Request) error {
	if err := requestRule(&req); err != nil {
		return ArgumentError{err}
	}

	body, err := ap.loader.Load(req.TemplateLang, req.TemplateName)
	if err != nil {
		e := validation.Error{
			Message: ErrLoadTemplate,
			Params: validation.Params{
				ParamTemplateLang:  req.TemplateLang,
				ParamTemplateName:  req.TemplateName,
				ParamTemplateCause: err.Error(),
			},
		}

		return ArgumentError{Err: validation.Errors([]error{e})}
	}

	key := fmt.Sprintf("%s-%s", req.TemplateLang, req.TemplateName)

	text, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	tml, err := template.New(key).Parse(string(text))
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := tml.Execute(buf, req.TemplateArgs); err != nil {
		return err
	}

	msg := model.Message{
		To:  []model.Address{},
		Cc:  []model.Address{},
		Bcc: []model.Address{},
	}
	if err := yaml.Unmarshal(buf.Bytes(), &msg); err != nil {
		return err
	}

	for _, rec := range req.To {
		msg.To = append(msg.To, rec)
	}

	for _, rec := range req.Cc {
		msg.Cc = append(msg.Cc, rec)
	}

	for _, rec := range req.Bcc {
		msg.Bcc = append(msg.Bcc, rec)
	}

	if err := messageRule(&msg); err != nil {
		return err
	}

	return ap.sender.Send(msg)
}
