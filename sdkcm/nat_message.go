package sdkcm

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type NatMessage struct {
	msg        *nats.Msg
	logger     *logrus.Logger
	statusCode int
}

func (s *NatMessage) Message() *nats.Msg {
	return s.msg
}

func (s *NatMessage) Logger() *logrus.Logger {
	return s.logger
}

func (s *NatMessage) StatusCode() int {
	return s.statusCode
}

func (s *NatMessage) SetStatusCode(code int) {
	s.statusCode = code
}

func (s *NatMessage) Respond(data interface{}) error {
	d, err := json.Marshal(data)
	if err != nil {
		// Mỗi lần request hay subscribe, publish của nat, nó sẽ tạo ra 1 object gọi là Msg
		// trong message này có hàm respond để trả kết quả về lại cho ai request, mình marshal rồi đẩy vào là xong
		s.msg.Respond([]byte(fmt.Sprintf("can not marshal object: %v", data)))
		return err
	}

	return s.msg.Respond(d)
}

func (s *NatMessage) BadRequestRespond(input interface{}, rootCause error, messages ...string) error {
	return s.Respond(NewBadRequestResponse(input, rootCause, messages...))
}

func (s *NatMessage) NotFoundRespond(input interface{}, rootCause error, messages ...string) error {
	return s.Respond(NewNotFoundResponse(input, rootCause, messages...))
}

func (s *NatMessage) ConflictRespond(input interface{}, rootCause error, messages ...string) error {
	return s.Respond(NewConflictResponse(input, rootCause, messages...))
}

func (s *NatMessage) UnauthorizedRespond(input interface{}, rootCause error, messages ...string) error {
	return s.Respond(NewUnauthorizedResponse(input, rootCause, messages...))
}

func (s *NatMessage) InternalServerErrorRespond(input interface{}, rootCause error, messages ...string) error {
	return s.Respond(NewInternalServerErrorResponse(input, rootCause, messages...))
}

func (s *NatMessage) SuccessRespond(input, data interface{}, paging *Paging) error {
	return s.Respond(SuccessResponse(input, data, paging))
}

func (s *NatMessage) Data() []byte {
	return s.msg.Data
}

func (s *NatMessage) Marshal(data interface{}) ([]byte, error) {
	d, err := json.Marshal(data)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}

	return d, nil
}

func (s *NatMessage) Unmarshal(data []byte, out interface{}) error {
	if err := json.Unmarshal(data, out); err != nil {
		s.logger.Errorln(err)
		return err
	}

	return nil
}

func (s *NatMessage) UnmarshalData(out interface{}) error {
	return s.Unmarshal(s.msg.Data, out)
}
