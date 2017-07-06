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
 *     Initial: 2017/04/11        Feng Yifei
 *      Modify: 2017/06/04        Yang Chenglong		ModifyFunction
 */

package endpoint

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"hypercube/access/endpoint/handler"

	"github.com/dgrijalva/jwt-go"
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

	server.server.Use(middleware.Logger())
	server.server.Use(middleware.Recover())
	server.server.Use(handler.LoginMiddleWare)
	server.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: node.Conf.CorsHosts,
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	config := middleware.JWTConfig{
		Claims:     jwt.MapClaims{},
		SigningKey: []byte(node.Conf.SecretKey),
	}
	server.server.Use(middleware.JWTWithConfig(config))

	server.server.GET("/join", server.serve())

	return server
}

func (server *HTTPServer) serve() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
