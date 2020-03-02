package sdkcm

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type NatHandler func(msg *NatMessage) error

type NatClient interface {
	Client() *nats.Conn
	Connect() error
	Subscribe(subject string, handler NatHandler) (*nats.Subscription, error)
	Publish(subject string, data interface{}) error
	Request(subject string, data interface{}) (*NatMessage, error)
	RequestRaw(subject string, data []byte) (*NatMessage, error)
	RequestReply(subject string, data []byte) (*Response, error)
	Marshal(data interface{}) ([]byte, error)
}

type natClient struct {
	client *nats.Conn
	logger *logrus.Logger
	cf     SDKConfig
}

func NewNatClient(logger *logrus.Logger, cf SDKConfig) *natClient {
	return &natClient{logger: logger, cf: cf}
}

func (s *natClient) Connect() error {
	var err error
	s.client, err = nats.Connect(s.cf.NatURL())
	if err != nil {
		return err
	}

	logger.Infof(`🎉 Connected to NAT server on "%s" !`, s.cf.NatURL())

	return nil
}

func (s *natClient) Client() *nats.Conn {
	return s.client
}

func (s *natClient) Subscribe(subject string, handler NatHandler) (*nats.Subscription, error) {
	sub, err := s.client.Subscribe(subject, func(msg *nats.Msg) {
		natMessage := &NatMessage{
			msg:    msg,
			logger: s.logger,
		}

		if err := handler(natMessage); err != nil {
			s.logger.Errorln(err)
		}
	})

	if err != nil {
		s.logger.Errorf(`❗️ Subscribe error on subject "%s": %s`, subject, err)
		return nil, err
	}

	return sub, nil
}

func (s *natClient) Publish(subject string, data interface{}) error {
	d, err := s.Marshal(data)
	if err != nil {
		return err
	}

	if err = s.client.Publish(subject, d); err != nil {
		s.logger.Errorf(`❗ Publish error on subject "%s": %s`, subject, err)
		return err
	}

	return nil
}

func (s *natClient) RequestReply(subject string, data []byte) (*Response, error) {
	msg, err := s.client.Request(subject, data, s.cf.RequestTimeout())
	if err != nil {
		s.logger.Errorf(`❗ Request error on subject "%s": %s`, subject, err)
		return nil, err
	}

	var res Response
	if err := json.Unmarshal(msg.Data, &res); err != nil {
		s.logger.Errorf(`❗ Unmarshal error: %s`, err)
		return nil, err
	}

	return &res, nil
}

func (s *natClient) Request(subject string, data interface{}) (*NatMessage, error) {
	d, err := s.Marshal(data)
	if err != nil {
		return nil, err
	}

	msg, err := s.client.Request(subject, d, s.cf.RequestTimeout())
	if err != nil {
		s.logger.Errorf(`❗ Request error on subject "%s": %s`, subject, err)
		return nil, err
	}

	natMsg := &NatMessage{
		msg:    msg,
		logger: s.logger,
	}

	return natMsg, nil
}

func (s *natClient) RequestRaw(subject string, data []byte) (*NatMessage, error) {
	msg, err := s.client.Request(subject, data, s.cf.RequestTimeout())
	if err != nil {
		s.logger.Errorf(`❗ Request error on subject "%s": %s`, subject, err)
		return nil, err
	}

	natMsg := &NatMessage{
		msg:    msg,
		logger: s.logger,
	}

	return natMsg, nil
}

func (s *natClient) Marshal(data interface{}) ([]byte, error) {
	d, err := json.Marshal(data)
	if err != nil {
		s.logger.Errorf(`❗ Marshal error: %s`, err)
		return nil, err
	}

	return d, nil
}
