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
 *     Initial: 2017/03/28        Feng Yifei
 */

package websocket

import (
	"net/http"
	"net"
	"hypercube/common/log"
)

const tcpType = "tcp4"

var logger *log.S8ELogger

func init() {
	logger = log.S8ECreateLogger(
		&log.S8ELogTag{
			log.LogTagService: "access layer",
			log.LogTagType: "websocket",
		},
		log.S8ELogLevelDefault)
}

type WebSocketServer struct {
	address *net.TCPAddr
	mux *http.ServeMux
	server *http.Server
}

func CreateWebSocketServer(address string, mux *http.ServeMux) (*WebSocketServer, error) {
	var (
		addr      *net.TCPAddr
		err       error
	)

	if addr, err = net.ResolveTCPAddr(tcpType, address); err != nil {
		logger.Error("WebSocket Server address format error:", err)
		return nil, err
	}

	ws := &WebSocketServer{
		address: addr,
		mux: mux,
		server: nil,
	}

	return ws, nil
}

func (this *WebSocketServer) Run() error {
	var (
		listener *net.TCPListener
		err      error
	)

	if listener, err = net.ListenTCP(tcpType, this.address); err != nil {
		logger.Error("Listen with error:", err)
		return err
	}

	this.server = &http.Server{Handler: this.mux}

	go func() {
		logger.Debug("WebSocket Server starting on:", this.address)

		if err = this.server.Serve(listener); err != nil {
			logger.Error("WebSocket Server start on:", this.address, " with error ", err)
			panic(err)
		}
	}()

	return nil
}

func (this * WebSocketServer) Shutdown() error {
	logger.Warn("WebSocket Server " ,this.address, " shutdown:")
	return nil
}
