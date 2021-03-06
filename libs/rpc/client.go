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
 *     Initial: 2017/07/05        Feng Yifei
 *     Modify : 2017/07/28        Yang Chenglong
 */

package rpc

import (
	"errors"
	"net"
	"net/rpc"
	"time"
)

const (
	dialTimeout  = 10 * time.Second
	callTimeout  = 3 * time.Second
	pingDuration = 1 * time.Second
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
	c.err = c.dial()

	return c
}

func (c *Client) dial() error {
	var (
		err  error
		conn net.Conn
	)

	conn, err = net.DialTimeout(c.options.Proto, c.options.Addr, dialTimeout)
	if err == nil {
		c.Client = rpc.NewClient(conn)
	}

	return err
}

// Call - Invoke the remote method by name, wait for the reply.
func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	var (
		err error
	)

	if c.Client == nil {
		return ErrRPCNotAvailable
	}

	select {
	case call := <-c.Client.Go(serviceMethod, args, reply, make(chan *rpc.Call, 1)).Done:
		err = call.Error
	case <-time.After(callTimeout):
		err = ErrRPCTimeout
	}

	return err
}

func (c *Client) Error() error {
	return c.err
}

// Close - Shutdown the rpc client.
func (c *Client) Close() {
	c.quit <- struct{}{}
}

// Ping - Ping the rpc server or reconnect when has an error.
func (c *Client) Ping(ping string) {
	var (
		arg   = ReqKeepAlive{}
		reply = RespKeepAlive{}
		err   error
	)

	for {
		select {
		case <-c.quit:
			goto closed
		default:
		}

		if c.Client != nil && c.err == nil {
			if err = c.Call(ping, &arg, &reply); err != nil {
				c.err = err
				if err != rpc.ErrShutdown {
					c.Client.Close()
				}
			}
		} else {
			if err = c.dial(); err == nil {
				c.err = nil
			}
		}
		time.Sleep(pingDuration)
	}

closed:
	if c.Client != nil {
		c.Client.Close()
	}
}
