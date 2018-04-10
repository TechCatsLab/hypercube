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
 *     Initial: 2017/07/06        Feng Yifei
 *     Modify : 2017/07/08        Liu Jiachang		ModifyFunction
 *     Modify : 2017/07/28        Yang Chenglong
 */

package endpoint

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/TechCatsLab/hypercube/access/endpoint/handler"
	"github.com/TechCatsLab/hypercube/access/rpc"
	"github.com/TechCatsLab/hypercube/libs/conn"
	"github.com/TechCatsLab/hypercube/libs/log"
	"github.com/TechCatsLab/hypercube/libs/message"
	"github.com/TechCatsLab/hypercube/libs/metrics/prometheus"
	"github.com/TechCatsLab/hypercube/libs/session"
)

// HTTPServer represents the http server accepts the client websocket connections.
type HTTPServer struct {
	node   *Endpoint
	server *echo.Echo
}

// NewHTTPServer create a http server.
func NewHTTPServer(node *Endpoint) *HTTPServer {
	server := &HTTPServer{
		node:   node,
		server: echo.New(),
	}

	rpc.RPCServer.Node = server.node

	server.server.Use(middleware.Logger())
	server.server.Use(middleware.Recover())

	config := middleware.JWTConfig{
		Claims:      jwt.MapClaims{},
		TokenLookup: "query:" + echo.HeaderAuthorization,
		SigningKey:  []byte(node.Conf.SecretKey),
	}
	server.server.Use(middleware.JWTWithConfig(config))

	server.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: node.Conf.CorsHosts,
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	server.server.GET("/join", server.serve())

	return server
}

func (server *HTTPServer) serve() echo.HandlerFunc {
	var (
		upgrader *websocket.Upgrader
	)

	upgrader = &websocket.Upgrader{
		ReadBufferSize:  server.node.Conf.WSReadBufferSize,
		WriteBufferSize: server.node.Conf.WSWriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
		if err != nil {
			log.Logger.Error("Upgrade error!", err)
			return err
		}

		user, err := handler.GetUser(c)
		if err != nil {
			return err
		}

		client, err := server.NewClient(user, server.node.clientHub(), session.NewSession(ws, user, server.node, server.node.hub.Mq()))
		if err != nil {
			log.Logger.Error("HTTPServer NewClient Error: %v", err)

			return err
		}

		err = server.ReadMessage(ws, client, user)
		if err != nil {
			log.Logger.Error("server ReadMessage Error: %v", err)
			return err
		}

		return nil
	}
}

func (server *HTTPServer) NewClient(user *message.User, hub *conn.ClientHub, session *session.Session) (*conn.Client, error) {
	var (
		err       error
		reply     int
		userEntry message.UserEntry
	)

	client := conn.NewClient(user, server.node.clientHub(), session)
	client.StartHandleMessage()

	userEntry = message.UserEntry{
		UserID:   *user,
		ServerIP: message.Access{ServerIP: server.node.Conf.Address},
	}

	rpcClient, err := rpc.RpcClients.Get(server.node.Conf.LogicAddrs)
	if err != nil || rpcClient == nil {
		client.Close()
		log.Logger.Error("Get rpcClients Error: %v", err)
		return nil, err
	}

	err = rpcClient.Call("LogicRPC.LoginHandler", userEntry, &reply)
	if err != nil {
		client.Close()
		log.Logger.Error("LogicRPC.LoginHandler Error: %v", err)
		return nil, err
	}

	hub.Add(user, client)
	prometheus.OnlineUserCounter.Add(1)

	log.Logger.Info("Endpoint info: ", server.node.Snapshot())

	return client, nil
}

func (server *HTTPServer) ReadMessage(ws *websocket.Conn, client *conn.Client, user *message.User) error {
	var msg message.Message

	defer server.removeUser(user, client)

	for {
		if err := ws.ReadJSON(&msg); err != nil {
			log.Logger.Error("ReadMessage Error: %v", err)
			return err
		} else {
			prometheus.ReceiveMessageCounter.Add(1)
			err = client.Handle(&msg)
			if err != nil {
				log.Logger.Error("Handle Message Error: %v", err)
				return err
			}
		}
	}
}

func (server *HTTPServer) removeUser(user *message.User, client *conn.Client) {
	var (
		userEntry = message.UserEntry{
			UserID:   *user,
			ServerIP: message.Access{ServerIP: server.node.Conf.Addrs},
		}
		reply int
	)

	server.node.clientHub().Remove(user)
	client.Close()
	prometheus.OnlineUserCounter.Desc()
	log.Logger.Info("Endpoint info: ", server.node.Snapshot())

	RpcClient, _ := rpc.RpcClients.Get(server.node.Conf.LogicAddrs)
	err := RpcClient.Call("LogicRPC.LogoutHandle", userEntry, &reply)
	if err != nil {
		log.Logger.Error("LogicRPC.LogoutHandle Error: %v", err)
	}
}
