/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
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
 *     Modify : 2017/07/28        Yang Chenglong
 */

package conn

import (
	"encoding/json"

	"github.com/TechCatsLab/hypercube/access/config"
	"github.com/TechCatsLab/hypercube/access/rpc"
	"github.com/TechCatsLab/hypercube/libs/log"
	msg "github.com/TechCatsLab/hypercube/libs/message"
	"github.com/TechCatsLab/hypercube/libs/metrics/prometheus"
	"github.com/TechCatsLab/hypercube/libs/session"
	"strings"
)

// Client is a client connection.
type Client struct {
	hub     *ClientHub
	user    *msg.User
	session *session.Session
}

// NewClient creates a client.
func NewClient(user *msg.User, hub *ClientHub, session *session.Session) *Client {
	client := &Client{
		hub:     hub,
		user:    user,
		session: session,
	}

	return client
}

// UID returns the user identify for this connection
func (client *Client) UID() string {
	return client.user.UserID
}

// Handle incoming messages
func (client *Client) Handle(message *msg.Message) error {
	var (
		ok  bool
		err error
	)

	switch message.Type {
	case msg.MessageTypePushPlainText, msg.MessageTypePlainText, msg.MessageTypeEmotion:
		rpcClient, err := rpc.RpcClients.Get(config.GNodeConfig.LogicAddrs)
		if err != nil {
			log.Logger.Error("Handle Get RpcClients Error: %v", err)
			return err
		}

		err = rpcClient.Call("LogicRPC.Add", message, &ok)
		if err != nil {
			log.Logger.Error("Call LogicRPC Add Error: %v", err)
			return err
		}
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
	var user msg.User
	var reply int

	err := json.Unmarshal(message.Content, &user)
	if err != nil {
		log.Logger.Error("HandleLogoutMessage Unmarshal Error: %v", err)
		return err
	}

	u := msg.User{
		UserID: strings.Trim(user.UserID, "\""),
	}

	userEntry := msg.UserEntry{
		UserID:   u,
		ServerIP: msg.Access{ServerIP: config.GNodeConfig.Addrs},
	}

	rpcClient, err := rpc.RpcClients.Get(config.GNodeConfig.LogicAddrs)
	if err != nil {
		log.Logger.Error("HandleLogoutMessage Get RpcClients Error: %v", err)
		return err

	}

	err = rpcClient.Call("LogicRPC.LogoutHandle", userEntry, &reply)
	if err != nil {
		log.Logger.Error("LogicRPC.LogoutHandle Error: %v", err)
		return err
	}

	client.hub.Remove(&u)
	client.Close()
	prometheus.OnlineUserCounter.Desc()

	return nil
}

func (client *Client) StartHandleMessage() {
	client.session.StartMessageLoop()
}

// Send messages from peers or push server
func (client *Client) Send(msg *msg.Message) {
	client.session.PushMessage(msg)
}

// Close finish the client message loop.
func (client *Client) Close() {
	client.session.Stop()
}
