package app_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vbogretsov/sendmail/app"
	"github.com/vbogretsov/sendmail/model"

	"github.com/vbogretsov/sendmail/test/app/fixture"
	"github.com/vbogretsov/sendmail/test/app/loader"
	"github.com/vbogretsov/sendmail/test/app/sender"
)

func TestSendMail(t *testing.T) {
	lr := loader.New()
	sd := sender.New()
	ap := app.New(lr, sd)

	for _, fx := range fixture.Fixtures {
		t.Run(fx.Name, func(t *testing.T) {
			err := ap.SendMail(fx.Request)
			require.Equal(t, fx.Result, err)
		})
	}

	t.Run("MailSent", func(t *testing.T) {
		to := []model.Address{
			{
				Email: "user@mail.com",
				Name:  "",
			},
		}

		req := model.Request{
			TemplateLang: loader.Lang,
			TemplateName: loader.TemplateValid,
			TemplateArgs: map[string]interface{}{
				"Username": "SuperUser",
			},
			To: to,
		}

		err := ap.SendMail(req)
		require.Nil(t, err)
		require.Len(t, sd.Inbox, 1)

		exp := model.Message{
			Subject:  "Subject",
			From:     model.Address{Email: "user@mail.com", Name: "Sender"},
			BodyType: "text/plain",
			Body:     fmt.Sprintf(loader.ExpectedBody, "SuperUser"),
			To:       to,
			Cc:       []model.Address{},
			Bcc:      []model.Address{},
		}
		act := sd.Inbox[0]
		require.Equal(t, exp, act)
	})
}
