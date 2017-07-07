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
	"hypercube/libs/message"
	"hypercube/libs/log"
	"sync"
)

// ClientHub represents a collection of client sessions.
type ClientHub struct {
	clients map[string]*Client
	mux     sync.Mutex
}

// NewClientHub creates a client hub.
func NewClientHub() *ClientHub {
	return &ClientHub{
		clients: map[string]*Client{},
		mux:     sync.Mutex{},
	}
}

// Add a client connection
func (hub *ClientHub) Add(user *message.User, client *Client) {
	hub.mux.Lock()
	defer hub.mux.Unlock()
	if _, exists := hub.clients[user.UserID]; exists {
		// Warning
		log.zapLog.Debug("user already login")
	}
	hub.clients[user.UserID] = client
}

// Remove a client connection
func (hub *ClientHub) Remove(user *message.User, client *Client) {
	hub.mux.Lock()
	defer hub.mux.Unlock()
	if _, exists := hub.clients[user.UserID]; !exists {
		// Warning
		log.zapLog.Debug("user hasn't login")
	}
	delete(hub.clients, user.UserID)
}

func (hub *ClientHub) Get(user string) (*Client, bool) {
	hub.mux.Lock()
	defer hub.mux.Unlock()

	client, ok := hub.clients[user];

	return client, ok
}

func (hub *ClientHub) PushMessageToAll(message *message.Message) {
	hub.mux.Lock()
	defer hub.mux.Unlock()

	for _, c := range hub.clients {
		c.session.Mq.PushMessage(message)
	}
}
