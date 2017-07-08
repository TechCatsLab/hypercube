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

package rpcc

import (
	"errors"
)

var (
	// ErrRPCNoClientAvailable - No client left in collection
	ErrRPCNoClientAvailable = errors.New("No rpc client available now")
)

// Clients - rpc client collections
type Clients struct {
	// rpc.Client is thread-safe, so Locker isn't required here
	clients []*Client
}

// Dials to rpc servers
func Dials(ops []Options) *Clients {
	clients := new(Clients)

	for _, op := range ops {
		clients.clients = append(clients.clients, Dial(op))
	}

	return clients
}

// get a usable client
func (c *Clients) get() (*Client, error) {
	for _, cli := range c.clients {
		if cli != nil && cli.Client != nil && cli.Error() == nil {
			return cli, nil
		}
	}
	return nil, ErrRPCNoClientAvailable
}

// Available checks if exists a available client.
func (c *Clients) Available() error {
	_, err := c.get()

	return err
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (c *Clients) Call(serviceMethod string, args interface{}, reply interface{}) error {
	var (
		err error
		cli *Client
	)

	if cli, err = c.get(); err == nil {
		err = cli.Call(serviceMethod, args, reply)
	}

	return err
}

// Ping the rpc connect and reconnect when has an error.
func (c *Clients) Ping() {
	for _, cli := range c.clients {
		go cli.Ping()
	}
}
