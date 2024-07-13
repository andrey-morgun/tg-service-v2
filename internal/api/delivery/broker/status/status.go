package status

import (
	"encoding/json"
	"github.com/andReyM228/lib/log"
	"github.com/andReyM228/lib/rabbit"
	"time"
)

type Handler struct {
	log         log.Logger
	serviceName string
	date        time.Time
	rabbit      rabbit.Rabbit
}

func NewHandler(log log.Logger, serviceName string, rabbit rabbit.Rabbit) Handler {
	return Handler{
		log:         log,
		serviceName: serviceName,
		rabbit:      rabbit,
		date:        time.Now().UTC(),
	}
}

func (h Handler) CheckStatus(request []byte) error {
	var req rabbit.RequestModel
	if err := json.Unmarshal(request, &req); err != nil {
		return err
	}

	return h.rabbit.Reply(req.ReplyTopic, 200, response{
		Name:      h.serviceName,
		BuildDate: h.date,
	})
}
