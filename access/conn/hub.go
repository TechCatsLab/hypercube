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
)

// ClientHub represents a collection of client sessions.
type ClientHub struct {
	clients map[string]*Client
}

// NewClientHub creates a client hub.
func NewClientHub() *ClientHub {
	return &ClientHub{}
}

// Add a client connection
func (hub *ClientHub) Add(user *message.User, client *Client) {
	if _, exists := hub.clients[user.UserID]; exists {
		// Warning
	}
	hub.clients[user.UserID] = client
}

// Remove a client connection
func (hub *ClientHub) Remove(user *message.User, client *Client) {
	if _, exists := hub.clients[user.UserID]; !exists {
		// Warning
	}
	delete(hub.clients, user.UserID)
}
