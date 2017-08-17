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

package endpoint

import (
	"context"
	"fmt"

	"github.com/fengyfei/hypercube/access/config"
	"github.com/fengyfei/hypercube/access/conn"
	"github.com/fengyfei/hypercube/libs/log"
	"github.com/fengyfei/hypercube/libs/message"
)

// Endpoint represents a access server.
type Endpoint struct {
	Conf     *config.NodeConfig
	ws       *HTTPServer
	hub      *conn.ClientHub
	shutdown chan struct{}
}

// NewEndpoint create a new access point.
func NewEndpoint(conf *config.NodeConfig) *Endpoint {
	var (
		ep *Endpoint
	)

	ep = &Endpoint{
		Conf:     conf,
		hub:      conn.NewClientHub(),
		shutdown: make(chan struct{}),
	}
	ep.ws = NewHTTPServer(ep)

	return ep
}

func (ep *Endpoint) clientHub() *conn.ClientHub {
	return ep.hub
}

func (ep *Endpoint) Send(user *message.User, msg *message.Message) {
	ep.hub.Send(user, msg)
}

// Run starts the access server.
func (ep *Endpoint) Run() error {
	log.Logger.Info("run %v", ep.ws.server.Start(ep.Conf.Addrs))

	select {
	case <-ep.shutdown:
		return ep.ws.server.Server.Shutdown(context.Background())
	}
}

// Shutdown stops the access server.
func (ep *Endpoint) Shutdown() {
	close(ep.shutdown)
}

// Snapshot view the struct information of the program runtime.
func (ep *Endpoint) Snapshot() string {
	return fmt.Sprintf("%+v", ep.clientHub())
}
