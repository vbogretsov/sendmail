package sendgrid

import (
	"encoding/json"
	"fmt"
	"net/http"

	api "github.com/sendgrid/sendgrid-go"
	log "github.com/sirupsen/logrus"

	"github.com/vbogretsov/sendmail/app"
	"github.com/vbogretsov/sendmail/model"
)

const v3URL = "/v3/mail/send"

type content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type personalization struct {
	Subject string          `json:"subject"`
	To      []model.Address `json:"to"`
	Cc      []model.Address `json:"cc,omitempty"`
	Bcc     []model.Address `json:"bcc,omitempty"`
}

type message struct {
	From            model.Address   `json:"from"`
	Personalization personalization `json:"personalization"`
	Content         content         `json:"content"`
}

// Sender represents a SendGrid sender.
type Sender struct {
	url string
	key string
}

// New creates new SendGrid sender.
func New(url, key string) (app.Sender, error) {
	s := Sender{
		url: url,
		key: key,
	}
	return &s, nil
}

// Send sends an email via SendGrid API.
func (s *Sender) Send(msg model.Message) error {
	data := message{
		From: msg.From,
		Personalization: personalization{
			Subject: msg.Subject,
			To:      msg.To,
			Cc:      msg.Cc,
			Bcc:     msg.Bcc,
		},
		Content: content{
			Type:  msg.BodyType,
			Value: msg.Body,
		},
	}

	args, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	request := api.GetRequest(s.key, v3URL, s.url)
	request.Method = http.MethodPost
	request.Body = args

	resp, err := api.API(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"request":      data,
			"responseCode": resp.StatusCode,
			"responseBody": resp.Body,
		}).Error("sendgrid call failed")

		return fmt.Errorf("sendgrid api error %d", resp.StatusCode)
	}

	return nil
}
