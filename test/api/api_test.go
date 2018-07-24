package api_test

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	testhook "github.com/sirupsen/logrus/hooks/test"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"

	"github.com/vbogretsov/sendmail/api"
	"github.com/vbogretsov/sendmail/app"

	"github.com/vbogretsov/sendmail/test/api/client"
	"github.com/vbogretsov/sendmail/test/api/fixture"
	"github.com/vbogretsov/sendmail/test/api/loader"
	"github.com/vbogretsov/sendmail/test/api/sender"
)

const (
	timeout = time.Second * 5
	qname   = "sendmail"
)

var amqpurl = flag.String(
	"amqpurl",
	"amqp://guest:guest@localhost",
	"AMQP broker URL")

func wait(action func() bool) error {
	c := make(chan bool)

	go func() {
		for {
			if action() {
				break
			}
		}
		close(c)
	}()

	select {
	case <-c:
		return nil
	case <-time.After(timeout):
		return errors.New("test timed out")
	}
}

func unmarshal(v interface{}) ([]fixture.JsonError, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	resp := []fixture.JsonError{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func TestApi(t *testing.T) {
	conn, err := amqp.Dial(*amqpurl)
	require.Nil(t, err)
	defer conn.Close()

	cli, err := client.New(conn, qname)
	require.Nil(t, err)
	defer cli.Close()

	lr := loader.New()
	sd := sender.New()
	ap := app.New(lr, sd)

	cnt, err := api.New(ap, qname, false, conn)
	require.Nil(t, err)
	defer cnt.Close()

	go func() {
		cnt.Start()
	}()

	logrus.SetOutput(ioutil.Discard)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	log := testhook.NewLocal(logrus.StandardLogger())

	for _, fx := range fixture.Fixtures {
		t.Run(fx.Name, func(t *testing.T) {
			cli.Send(fx.Request)

			wait(func() bool {
				return len(log.Entries) > 0
			})

			require.Len(t, log.Entries, 1)

			data, err := unmarshal(log.LastEntry().Data["error"])
			require.Nil(t, err)

			require.Equal(t, fx.Errors, data)
			log.Reset()
		})
	}
}
