/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2017/07/06        Feng Yifei
 */

package session

import (
	"encoding/json"

	"github.com/gorilla/websocket"

	"hypercube/libs/message"
	"hypercube/access/sender"
	"hypercube/libs/metrics/prometheus"
	"hypercube/libs/log"
)

// Session represents a client connection.
type Session struct {
	user     *message.User
	sender   sender.Sender
	mq       *MessageQueue
	ws       *websocket.Conn
	shutdown chan struct{}
}

// NewSession creates a session.
func NewSession(ws *websocket.Conn, user *message.User, sender sender.Sender, buffSize int) *Session {
	session := &Session{
		mq:             NewMessageQueue(buffSize),
		ws:             ws,
		user:           user,
		sender:         sender,
		shutdown:       make(chan struct{}),
	}

	return session
}

// PushMessage push message to queue.
func (s *Session)PushMessage(msg *message.Message) {
	s.mq.PushMessage(msg)
}

// StartMessageLoop start to handle message loop.
func (s *Session) StartMessageLoop() {
	go func() {
		for {
			select {
			case msg := <-s.mq.FetchMessage():
				s.HandleMessage(msg)
			case <-s.shutdown:
				return
			}
		}
	}()
}

// HandleMessage handle message method.
func (s *Session) HandleMessage(msg *message.Message) {
	var (
		user *message.User
		plainMsg message.PlainText
		pushMsg  message.PushPlainText
	)

	switch msg.Type {
	case message.MessageTypePlainText:
		json.Unmarshal(msg.Content, &plainMsg)
		user = &plainMsg.To
	case message.MessageTypePushPlainText:
		json.Unmarshal(msg.Content, &pushMsg)
		user = &pushMsg.To
	default:
		log.Logger.Warn("Unknown Type!")
		return
	}

	if user.UserID == s.user.UserID {
		err := s.ws.WriteJSON(*msg)
		if err != nil {
			log.Logger.Error("WriteMessage Error: %v", err)
			s.PushMessage(msg)
		}
		prometheus.SendMessageCounter.Add(1)
	} else {
		s.sender.Send(user, msg)
	}
}

// Stop stop handle message loop.
func (s *Session) Stop() {
	close(s.shutdown)
}
