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

package conn

import (
	"encoding/json"

	"hypercube/access/session"
	"hypercube/libs/log"
	msg "hypercube/libs/message"
	"hypercube/libs/metrics/prometheus"
)

// Client is a client connection.
type Client struct {
	user    *msg.User
	hub     *ClientHub
	session *session.Session
}

// NewClient creates a client.
func NewClient(user *msg.User, hub *ClientHub, session *session.Session) *Client {
	client := &Client{
		user:    user,
		hub:     hub,
		session: session,
	}

	hub.Add(user, client)

	return client
}

// UID returns the user identify for this connection
func (client *Client) UID() string {
	return client.user.UserID
}

// Handle incoming messages
func (client *Client) Handle(message *msg.Message) error {
	var err error
	switch message.Type {
	case msg.MessageTypePushPlainText, msg.MessageTypePlainText:
		client.Send(message)
	case msg.MessageTypeLogout:
		err = client.HandleLogoutMessage(message)
	default:
		log.Logger.Debug("No message type match!")
	}

	if err != nil {
		log.Logger.Error("Handle Message Error: ", err)
		return err
	}

	return nil
}

func (client *Client) HandleLogoutMessage(message *msg.Message) error {
	var mess msg.User

	err := json.Unmarshal(message.Content, &mess)
	if err != nil {
		return err
	}

	client.hub.Remove(&mess, client)
	prometheus.OnlineUserCounter.Add(-1)

	return nil
}

func (client *Client) StartHandleMessage() {
	client.session.StartMessageLoop()
}

func (client *Client) StopHandleMessage() {
	client.session.Stop()
}

// Send messages from peers or push server
func (client *Client) Send(msg *msg.Message) {
	client.session.PushMessage(msg)
}

// Close finish the client message loop.
func (client *Client) Close() {
}
