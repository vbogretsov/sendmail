package sender

import (
	"sync"

	"github.com/vbogretsov/sendmail/model"
)

type Sender struct {
	mutex sync.Mutex
	Inbox []model.Message
	Error error
}

func New() *Sender {
	return &Sender{Inbox: []model.Message{}}
}

func (self *Sender) Send(msg model.Message) error {
	if self.Error != nil {
		return self.Error
	}

	self.mutex.Lock()
	self.Inbox = append(self.Inbox, msg)
	self.mutex.Unlock()
	return nil
}
