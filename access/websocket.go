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
 *	   AddFunction: 2017/04/06    Yusan Kurban
 */

package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	ws "hypercube/common/server/websocket"
	"hypercube/proto/general"
	"encoding/json"
	"hypercube/proto/api"
	"time"
	"strings"
	"github.com/labstack/echo"
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

	server.GET("/join", serveWebSocket)

	logger.Debug("Configuration finished, starting servers...")

	if err = initWebSocketServer(); err != nil {
		panic(err)
	}
}

func initWebSocketServer() error {
	var (
		server      *ws.WebSocketServer
		err         error
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

func sendAccessInfo()  {
	var (
		r           api.Reply
		serverinfo  api.Access
	)

	for _, address := range configuration.Addrs {
		addr := strings.Split(address, ":")[0]
		serverinfo.ServerIp = &addr
		serverinfo.Subject = &configuration.Subject

		info, _ := json.Marshal(serverinfo)

		err := logicRequester.Request(&api.Request{Type: api.ApiTypeAccessInfo, Content: info}, &r, time.Duration(100) * time.Millisecond)

		if err != nil {
			logger.Error("send access info error:", err, " , address:", address)

			break
		}

		logger.Debug("send access info to logic:", serverinfo, " received reply:", r.Code)
	}
}

func serveWebSocket(c echo.Context) error {
	logger.Debug("New connection")

	if c.Request().Method != "GET" {
		return c.JSON(http.StatusMethodNotAllowed, "Method Not Allowed")
	}

	connection, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		logger.Error("websocket upgrade error:", err)
		return err
	}
	defer connection.Close()

	webSocketConnectionHandler(connection)

	return c.JSON(http.StatusInternalServerError,"serveWebSocket connection error")
}

type handlerFunc func(p interface{},req interface{}) interface{}

func webSocketConnectionHandler(conn *websocket.Conn) {
	var (
		err        error
		p          *general.Proto = &general.Proto{}
		ver        *general.Keepalive = &general.Keepalive{}
		mes        *general.Message = &general.Message{}
		user       *general.UserAccess = &general.UserAccess{}
		v          interface{}
		handler    handlerFunc
		id 	   general.UserKey
		ok 	   bool
	)

	for {
		if err = p.ReadWebSocket(conn); err != nil {
			id, ok = OnLineManagement.GetIDByConnection(conn)
			if ok {
				err = OnLineManagement.OnDisconnect(id)
				if err != nil {
					logger.Error("User Logout failed:", err)
				}
			}
			logger.Error("conn read error:", err)
			break
		}

		logger.Debug("Websocket received message type:", p.Type)

		switch p.Type {
		case general.GeneralTypeKeepAlive:
			v = ver
			handler = keepAliveRequestHandler
		case general.GeneralTypeTextMsg:
			v = mes
			handler = userMessageHandler
		case general.GeneralTypeLogin:
			v = user
		case general.GeneralTypeLogout:
			v = user
		}

		if v != nil {
			err = json.Unmarshal(p.Body, v)

			if err != nil {
				logger.Error("Receive unknown message:", err)
				continue
			} else {
				switch p.Type {
				case general.GeneralTypeLogin:
					user = v.(*general.UserAccess)
					err = OnLineManagement.OnConnect(user.UserID, conn)
					if err != nil {
						logger.Error("User Login failed:", err)
					}
				case general.GeneralTypeLogout:
					user = v.(*general.UserAccess)
					err = OnLineManagement.OnDisconnect(user.UserID)
					if err != nil {
						logger.Error("User Logout failed:", err)
					}
				default:
					handler(p,v)
				}
			}
		}
	}
}
