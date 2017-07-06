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
	"hypercube/libs/message"
)

// Client is a client connection.
type Client struct {
	user    *message.User
	hub     *ClientHub
	session *session.Session
}

// NewClient creates a client.
func NewClient(user *message.User, hub *ClientHub, session *session.Session) *Client {
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
	return ""
}

// Handle incoming messages
func (client *Client) Handle(message *message.Message) error {
	return nil
}

// Send messages from peers or push server
func (client *Client) Send() error {
	return nil
}

// Close finish the client message loop.
func (client *Client) Close() {
}
