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
 *     Initial: 2017/03/29        Feng Yifei
 */

package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	ws "hypercube/common/server/websocket"
)

var (
	upgrader              *websocket.Upgrader
	mux                   *http.ServeMux
	webSocketServers      []*ws.WebSocketServer
)

func initWebsocket()  {
	var err error

	upgrader = &websocket.Upgrader{
		ReadBufferSize:     configuration.WSReadBufferSize,
		WriteBufferSize:    configuration.WSWriteBufferSize,
		CheckOrigin:        func(r *http.Request) bool {
			return true
		},
	}

	mux = http.NewServeMux()
	mux.HandleFunc("/join", serveWebSocket)

	logger.Debug("Configuration finished, starting servers...")

	if err = initWebSocketServer(); err != nil {
		panic(err)
	}
}

func initWebSocketServer() error {
	var (
		server    *ws.WebSocketServer
		err       error
	)

	webSocketServers = make([]*ws.WebSocketServer, len(configuration.Addrs))

	for index, address := range configuration.Addrs {
		server, err = ws.CreateWebSocketServer(address, mux)

		if err != nil {
			logger.Error("start websocket server error:", err, " , address:", address)

			break
		}

		logger.Debug("start websocket server succeed at address:", address)

		webSocketServers[index] = server
		server.Run()
	}

	return err
}

func serveWebSocket(w http.ResponseWriter, req *http.Request) {
	logger.Debug("New connection")

	if req.Method != "GET" {
		// TODO: 考虑加入黑名单逻辑
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	connection, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		logger.Error("websocket upgrade error:", err)
		return
	}
	defer connection.Close()

	webSocketConnectionHandler(connection)
}

func webSocketConnectionHandler(conn *websocket.Conn) {
	for {
	}
}
