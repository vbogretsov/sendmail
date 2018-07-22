package sender

import (
	"fmt"

	"github.com/vbogretsov/sendmail/app"
	"github.com/vbogretsov/sendmail/app/sender/sendgrid"
)

type factory func(url, key string) (app.Sender, error)

var senders = map[string]factory{
	"sendgrid": sendgrid.New,
}

// New creates a new sender.
func New(name, url, key string) (app.Sender, error) {
	fn, ok := senders[name]
	if !ok {
		return nil, fmt.Errorf("unsupported sender %s", name)
	}
	return fn(url, key)
}

// Providers gets lost of available providers.
func Providers() []string {
	providers := []string{}
	for k, _ := range senders {
		providers = append(providers, k)
	}
	return providers
}
