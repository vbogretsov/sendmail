package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/vbogretsov/sendmail/api"
	"github.com/vbogretsov/sendmail/app"
	"github.com/vbogretsov/sendmail/app/loader"
	"github.com/vbogretsov/sendmail/app/sender"
)

const (
	name = "sendmail"
	desc = "email microservice"
)

var Version = ""

const (
	helpProviderName = "SMTP provider name [%v]"
	helpProviderURL  = "SMTP provider API endpoint"
	helpProviderKey  = "SMTP provider authorization key"
	helpAMQPURL      = "AMQP brocker URL"
	helpAMQPQName    = "AMQP quee listening name"
	helpTemplatePath = "templates root location"
	helpLogLevel     = "log level [%v]"
)

type argT struct {
	provider struct {
		Name *string
		URL  *string
		Key  *string
	}
	amqp struct {
		URL   *string
		QName *string
	}
	template struct {
		Path *string
	}
	log struct {
		Level *string
	}
}

var (
	args      = argT{}
	parser    = argparse.NewParser(fmt.Sprintf("%s %s", name, Version), desc)
	logLevels = []string{
		"panic",
		"fatal",
		"error",
		"warn",
		"info",
		"debug",
	}
)

func init() {
	args.provider.Name = parser.String(
		"",
		"provider-name",
		&argparse.Options{
			Required: true,
			Help:     fmt.Sprintf(helpProviderName, sender.Providers()),
		})
	args.provider.URL = parser.String(
		"",
		"provider-url",
		&argparse.Options{
			Required: true,
			Help:     helpProviderURL,
		})
	args.provider.Key = parser.String(
		"",
		"provider-key",
		&argparse.Options{
			Required: true,
			Help:     helpProviderKey,
		})
	args.template.Path = parser.String(
		"",
		"templates-path",
		&argparse.Options{
			Required: true,
			Help:     helpTemplatePath,
		})
	args.amqp.URL = parser.String(
		"",
		"amqp-url",
		&argparse.Options{
			Required: false,
			Default:  "amqp://guest:guest@localhost",
			Help:     helpAMQPURL,
		})
	args.amqp.QName = parser.String("", "amqp-qname", &argparse.Options{
		Required: false,
		Default:  name,
		Help:     helpAMQPURL,
	})
	args.log.Level = parser.Selector(
		"",
		"log-level",
		logLevels,
		&argparse.Options{
			Required: false,
			Default:  "info",
			Help:     fmt.Sprintf(helpLogLevel, logLevels),
		})
}

func run() error {
	if err := parser.Parse(os.Args); err != nil {
		return err
	}

	lr, err := loader.New(*args.template.Path)
	if err != nil {
		return err
	}

	sr, err := sender.New(
		*args.provider.Name,
		*args.provider.URL,
		*args.provider.Key)
	if err != nil {
		return err
	}

	ap := app.New(lr, sr)

	lv, err := log.ParseLevel(*args.log.Level)
	if err != nil {
		return err
	}
	log.SetLevel(lv)

	log.SetFormatter(&log.JSONFormatter{})

	cn, err := amqp.Dial(*args.amqp.URL)
	if err != nil {
		return err
	}
	defer cn.Close()

	cnt, err := api.New(ap, *args.amqp.QName, true, cn)
	if err != nil {
		return err
	}
	defer cnt.Close()

	cnt.Start()
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
