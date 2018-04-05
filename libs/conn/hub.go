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
	"errors"
	"sync"

	"github.com/TechCatsLab/hypercube/libs/log"
	"github.com/TechCatsLab/hypercube/libs/message"
)

// ClientHub represents a collection of client sessions.
type ClientHub struct {
	mq      *MessageQueue
	clients sync.Map
}

// NewClientHub creates a client hub.
func NewClientHub(buffSize int) *ClientHub {
	return &ClientHub{
		clients: sync.Map{},
		mq:      NewMessageQueue(buffSize),
	}
}

func (hub *ClientHub) Mq() *MessageQueue {
	return hub.mq
}

// Add a client connection
func (hub *ClientHub) Add(user *message.User, client *Client) {
	if _, exists := hub.clients.Load(user.UserID); exists {
		log.Logger.Warn("user already login!")
	}

	hub.clients.Store(user.UserID, client)
}

// Remove a client connection
func (hub *ClientHub) Remove(user *message.User) {
	if _, exists := hub.clients.Load(user.UserID); !exists {
		log.Logger.Warn("user hasn't login!")
		return
	}

	hub.clients.Delete(user.UserID)
}

// Get a client by user
func (hub *ClientHub) Get(user string) (*Client, bool) {
	client, ok := hub.clients.Load(user)
	if !ok {
		return nil, false
	}

	return client.(*Client), true
}

func (hub *ClientHub) PushMessageToAll(msg *message.Message) {
	hub.clients.Range(func(key, value interface{}) bool {
		value.(*Client).session.PushMessage(msg)
		return true
	})
}

func (hub *ClientHub) Send(user *message.User, msg *message.Message) error {
	client, exist := hub.Get(user.UserID)
	if !exist {
		log.Logger.Debug("user hasn't login!")
		return errors.New("User Hasn't Login")
	}

	client.Send(msg)

	return nil
}

func (hub *ClientHub) GetAllUser() []*Client {
	var clients []*Client

	hub.clients.Range(func(key, value interface{}) bool {
		clients = append(clients, value.(*Client))
		return true
	})

	return clients
}
