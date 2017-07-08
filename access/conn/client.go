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
	"hypercube/access/session"
	mes "hypercube/libs/message"
	"encoding/json"
	"hypercube/libs/metrics/prometheus"
	"hypercube/libs/log"
)

// Client is a client connection.
type Client struct {
	user    *mes.User
	hub     *ClientHub
	session *session.Session
}

var HubService = &ClientHub{}

// NewClient creates a client.
func NewClient(user *mes.User, hub *ClientHub, session *session.Session) *Client {
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
func (client *Client) Handle(message *mes.Message) error {
	var err error
	switch message.Type {
	case mes.MessageTypePlainText:
		err = client.HandleUserMessage(message)
	case mes.MessageTypePushPlainText:
		err = client.HandlePushMessage(message)
	case mes.MessageTypeLogout:
		err = client.HandleLogoutMessage(message)
	default:

	}

	if err != nil {
		log.Logger.Error("Handle Message Error: ", err)

		return err
	}
	return nil
}

func (client *Client) HandleUserMessage(message *mes.Message) error {
	var (
		mess mes.PlainText
		v    interface{}
	)

	switch message.Type {
	case mes.MessageTypePlainText:
		v = mess

	}

	err := json.Unmarshal(message.Content, &v)
	if err != nil {
		return err
	}

	client.hub.Send(*mess.To, message)

	return nil
}

func (client *Client) HandlePushMessage(pmessage *mes.Message) error {
	var pmess mes.PushPlainText

	err := json.Unmarshal(pmessage.Content, &pmess)
	if err != nil {
		return err
	}

	switch pmess.Type {
	case mes.PushToAll:

		HubService.PushMessageToAll(pmessage)
		return nil

	default:
		for _, user := range pmess.To {
			to, ok := HubService.Get(user.UserID)
			if !ok {
				continue
			}

			to.session.PushMessage(pmessage)
		}
		return nil
	}
}

func (client *Client) HandleLogoutMessage(message *mes.Message) error {
	var mess mes.User

	err := json.Unmarshal(message.Content, &mess)
	if err != nil {
		return err
	}

	client.hub.Remove(&mess, client)
	prometheus.OnlineUserCounter.Add(-1)

	return nil
}

// Send messages from peers or push server
func (client *Client) Send(msg *mes.Message) error {
	var mess mes.PlainText

	err := json.Unmarshal(msg.Content, &mess)
	if err != nil {
		log.Logger.Error("Client Send Unmarshal Message Error: %v", err)

		return err
	}

	client.session.PushMessage(mess)

	return nil
}

// Close finish the client message loop.
func (client *Client) Close() {
}
