package api

import (
	"encoding/json"

	"github.com/labstack/gommon/random"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/vbogretsov/go-validation"
	jsonerr "github.com/vbogretsov/go-validation/json"

	"github.com/vbogretsov/sendmail/app"
	"github.com/vbogretsov/sendmail/model"
)

const idsize = 32

type ErrorMarshaler func(error) interface{}

// Api represents sendmail AMQP API.
type Api struct {
	ap      *app.App
	ch      *amqp.Channel
	rq      <-chan amqp.Delivery
	requeue bool
}

// New creates new Api.
func New(ap *app.App, qname string, requeue bool, cn *amqp.Connection) (*Api, error) {
	var err error

	ch, err := cn.Channel()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			ch.Close()
		}
	}()

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
		return nil, err
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
		return nil, err
	}

	err = ch.QueueBind(
		qe.Name, // queue name
		qname,   // routing key
		qname,   // exchange
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &Api{ap: ap, ch: ch, requeue: requeue, rq: reqs}, nil
}

// Start starts listening for email requests.
func (self *Api) Start() {
	for i := range self.rq {
		reqid := random.String(idsize)

		err := sendmail(reqid, self.ap, i)
		if err != nil {
			switch err.(type) {
			case app.ArgumentError:
				log.WithFields(log.Fields{
					"id": reqid,
					"error": marshal(
						err.(app.ArgumentError).Err.(validation.Errors)),
				}).Error("invalid request")
				i.Nack(false, false)
			case app.TemplateError:

				log.WithFields(log.Fields{
					"id": reqid,
					"error": marshal(
						err.(app.TemplateError).Err.(validation.Errors)),
				}).Error("invalid template")
				i.Nack(false, self.requeue)
			default:
				log.WithFields(log.Fields{
					"id":    reqid,
					"error": err,
				}).Error("unable to send email")
				i.Nack(false, self.requeue)
			}
		} else {
			i.Ack(false)
		}
	}
}

// Close closes the underlying channel.
func (self *Api) Close() error {
	return self.ch.Close()
}

func sendmail(reqid string, ap *app.App, i amqp.Delivery) error {
	req := model.Request{}

	if err := json.Unmarshal(i.Body, &req); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"request": req,
		"id":      reqid,
	}).Debug("request received")

	if err := ap.SendMail(req); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"request": req,
		"id":      reqid,
	}).Debug("request completed")

	return nil
}

func marshal(err validation.Errors) json.Marshaler {
	return jsonerr.New(err, jsonerr.DefaultFormatter, jsonerr.DefaultJoiner)
}
