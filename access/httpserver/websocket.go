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
 *     AddEcho: 2017/06/04        Yang Chenglong
 *     Modify:  2017/06/07        Yang Chenglong    添加接收消息数量统计
 */

package httpserver

import (
	"encoding/json"
	"hypercube/libs/log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

func sendAccessInfo() {
	var (
		r          general.Reply
		serverinfo general.Access
	)

	addr := strings.Split(configuration.Addrs, ":")[0]

	serverinfo.ServerIp = &addr
	serverinfo.Subject = &configuration.Subject

	info, _ := json.Marshal(serverinfo)

	err := logicRequester.Request(&general.Request{Type: types.ApiTypeAccessInfo, Content: info}, &r, time.Duration(100)*time.Millisecond)

	if err != nil {
		log.GlobalLogger.Error("send access info error:", err, " , address:", configuration.Addrs)
	}

	log.GlobalLogger.Debug("send access info to logic:", serverinfo, " received reply:", r.Code)
}

func serveWebSocket(c echo.Context) error {
	var upgrader = &websocket.Upgrader{
		ReadBufferSize:  configuration.WSReadBufferSize,
		WriteBufferSize: configuration.WSWriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	log.GlobalLogger.Debug("New connection")

	if c.Request().Method != "GET" {
		return c.JSON(http.StatusMethodNotAllowed, "Method Not Allowed")
	}

	connection, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		log.GlobalLogger.Error("websocket upgrade error:", err)
		return err
	}
	defer connection.Close()

	webSocketConnectionHandler(connection)

	return c.JSON(http.StatusInternalServerError, "serveWebSocket connection error")
}

type handlerFunc func(p interface{}, req interface{}) interface{}

func webSocketConnectionHandler(conn *websocket.Conn) {
	var (
		err     error
		p       *general.Proto      = &general.Proto{}
		ver     *general.Keepalive  = &general.Keepalive{}
		mes     *general.Message    = &general.Message{}
		user    *general.UserAccess = &general.UserAccess{}
		v       interface{}
		handler handlerFunc
		id      general.UserKey
		ok      bool
	)

	for {
		if err = p.ReadWebSocket(conn); err != nil {
			id, ok = OnLineManagement.GetIDByConnection(conn)
			if ok {
				err = OnLineManagement.OnDisconnect(id)
				if err != nil {
					log.GlobalLogger.Error("User Logout failed:", err)
				}
			}
			log.GlobalLogger.Error("conn read error:", err)
			break
		}

		log.GlobalLogger.Debug("Websocket received message type:", p.Type)

		switch p.Type {
		case types.GeneralTypeKeepAlive:
			v = ver
			handler = keepAliveRequestHandler
		case types.GeneralTypeTextMsg:
			v = mes
			handler = userMessageHandler
			receiveMessageCounter.Add(1)
		case types.GeneralTypeLogin:
			v = user
		case types.GeneralTypeLogout:
			v = user
		}

		if v != nil {
			err = json.Unmarshal(p.Body, v)

			if err != nil {
				log.GlobalLogger.Error("Receive unknown message:", err)
				continue
			} else {
				switch p.Type {
				case types.GeneralTypeLogin:
					user = v.(*general.UserAccess)
					err = OnLineManagement.OnConnect(user.UserID, conn)
					if err != nil {
						log.GlobalLogger.Error("User Login failed:", err)
					}
				case types.GeneralTypeLogout:
					user = v.(*general.UserAccess)
					err = OnLineManagement.OnDisconnect(user.UserID)
					if err != nil {
						log.GlobalLogger.Error("User Logout failed:", err)
					}
				default:
					handler(p, v)
				}
			}
		}
	}
}
