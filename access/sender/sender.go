package sender

import (
	"hypercube/libs/message"
)

type Sender interface {
	Send(user *message.User, msg *message.Message)
}
