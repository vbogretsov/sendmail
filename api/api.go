package api

import (
	"encoding/json"

	"github.com/vbogretsov/go-validation"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	jsonerr "github.com/vbogretsov/go-validation/json"

	"github.com/vbogretsov/sendmail/app"
	"github.com/vbogretsov/sendmail/model"
)

// Run starts consuming client requests.
func Run(ap *app.App, url, qname string) error {
	log.Debugf("connecting AMQP broker %s", url)

	cn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	defer cn.Close()

	ch, err := cn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		qname,    // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	qe, err := ch.QueueDeclare(
		qname, // name
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		qe.Name, // queue name
		qname,   // routing key
		qname,   // exchange
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return err
	}

	reqs, err := ch.Consume(
		qe.Name, // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return err
	}

	log.Debug("AMQP broker connected")

	for i := range reqs {
		err := sendmail(ap, i)
		if err != nil {
			switch err.(type) {
			case app.ArgumentError:
				e := jsonerr.New(
					err.(validation.Errors),
					jsonerr.DefaultFormatter,
					jsonerr.DefaultJoiner)
				log.WithFields(log.Fields{
					"error": e,
				}).Error("invalid request parameters")
			default:
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Fatal("unable to send email")
			}
			i.Nack(false, true)
		} else {
			i.Ack(false)
		}
	}

	return nil
}

func sendmail(ap *app.App, i amqp.Delivery) error {
	req := model.Request{}

	if err := json.Unmarshal(i.Body, &req); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"request": req,
	}).Debug("received request")

	if err := ap.SendMail(req); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"request": req,
	}).Debug("mail sent")

	return nil
}
