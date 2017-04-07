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
	"hypercube/proto/push"
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

	webSocketUserHandler(connection) // Todo: 与下面函数合并
	webSocketConnectionHandler(connection)
}

type handlerFunc func(p interface{},req interface{}) interface{}

func webSocketConnectionHandler(conn *websocket.Conn) {
	var (
		err        error
		p          *general.Proto = &general.Proto{}
		ver        *general.Keepalive = &general.Keepalive{}
		mes        *general.Message = &general.Message{}
		pushmany   *push.MPushMsg = &push.MPushMsg{}
		pushone    *push.PushMsg = &push.PushMsg{}
		v          interface{}
		handler    handlerFunc
	)

	for {
		if err = p.ReadWebSocket(conn); err != nil {
			logger.Error("conn read error:", err)
			break
		}

		switch p.Type {
		case general.TpHeartbeat:
			v = ver
			handler = keepAliveRequestHandler
		case general.TpUTUMsg:
			v = mes
			handler = userToUserRequestHandler
		case general.TpRoomMsg: // Todo: Room 相关代码清理掉
			v = pushmany
			handler = pushManyRequestHandler
		case general.TpPushMsg: // Todo: Push 消息处理不在这里
			v = pushone
			handler = pushOneRequestHandler
		}

		if v != nil {
			err = json.Unmarshal(p.Body, v)
			if err != nil {
				// Todo: 记录错误
				continue
			} else {
				handler(p,v)
			}
		}
	}
}

type user struct {
	Userid 		string
}

func webSocketUserHandler (conn *websocket.Conn) {
	var (
		p 		*general.Proto = &general.Proto{}
		err 	error
		user 	user
	)

	err = p.ReadWebSocket(conn)
	if err != nil {
		logger.Error(err)
	}

	err = json.Unmarshal(p.Body, user)
	if err != nil {
		logger.Error(err)
	}

	OnLineUser.OnConnect(user.Userid, conn)
}