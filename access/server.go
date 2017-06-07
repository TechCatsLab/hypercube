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
 *     ModifyFunction: 2017/06/04        Yang Chenglong
 */

package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	e "hypercube/common/error"
	"hypercube/access/echo/handler"
	"github.com/dgrijalva/jwt-go"
)

var (
	server *echo.Echo
)

func initEchoServer() {
	server = echo.New()

	server.HTTPErrorHandler = e.HTTPErrorHandler

	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(handler.LoginMiddleWare)
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:    configuration.CorsHosts,
		AllowMethods:    []string{echo.GET, echo.POST},
	}))

	config := middleware.JWTConfig{
		Claims:         jwt.MapClaims{},
		SigningKey:     []byte(configuration.SecretKey),
	}

	server.Use(middleware.JWTWithConfig(config))
}