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
 *     Initial: 2017/07/05        Feng Yifei
 */

package rpcc

import (
	"errors"
	"net"
	"net/rpc"
	"rpcx/log"
	"time"
)

var (
	// ErrRPCNotAvailable - rpc service not available
	ErrRPCNotAvailable = errors.New("rpc service not available")

	// ErrRPCTimeout - rpc server timer out
	ErrRPCTimeout = errors.New("rpc dial timeout")
)

// Options - rpc client options.
type Options struct {
	Proto string
	Addr  string
}

// Client - rpc client
type Client struct {
	*rpc.Client
	options Options
	quit    chan struct{}
	err     error
}

// Dial - Dial to rpc server.
func Dial(options Options) *Client {
	c := new(Client)
	c.options = options
	c.dial()
	return c
}

func (c *Client) dial() (err error) {
	var conn net.Conn
	conn, err = net.DialTimeout(c.options.Proto, c.options.Addr, dialTimeout)
	if err != nil {
		log.Error("net.Dial(%s, %s), error(%v)", c.options.Proto, c.options.Addr, err)
	} else {
		c.Client = rpc.NewClient(conn)
	}
	return
}

// Call - Invoke the remote method by name, wait for the reply.
func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) (err error) {
	if c.Client == nil {
		err = ErrRpc
		return
	}
	select {
	case call := <-c.Client.Go(serviceMethod, args, reply, make(chan *rpc.Call, 1)).Done:
		err = call.Error
	case <-time.After(callTimeout):
		err = ErrRpcTimeout
	}
	return
}

func (c *Client) Error() error {
	return c.err
}

// Close - Shutdown the rpc client.
func (c *Client) Close() {
	c.quit <- struct{}{}
}

// Ping - Ping the rpc server or reconnect when has an error.
func (c *Client) Ping(serviceMethod string) {

}
