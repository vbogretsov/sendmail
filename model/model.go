package model

// Address represents an email address.
type Address struct {
	Email string `yaml:"Email" json:"email,omitempty"`
	Name  string `yaml:"Name" json:"name,omitempty"`
}

// Message represent an email message.
type Message struct {
	From     Address   `yaml:"From"`
	To       []Address `yaml:"To"`
	Cc       []Address `yaml:"Cc"`
	Bcc      []Address `yaml:"Bcc"`
	Subject  string    `yaml:"Subject"`
	BodyType string    `yaml:"BodyType"`
	Body     string    `yaml:"Body"`
}

// Request represents request parameters required to build and send an email.
type Request struct {
	TemplateLang string                 `json:"templateLang"`
	TemplateName string                 `json:"templateName"`
	TemplateArgs map[string]interface{} `json:"templateArgs"`
	To           []Address              `json:"to"`
	Cc           []Address              `json:"cc"`
	Bcc          []Address              `json:"bcc"`
}
