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
 *     Modify: 2017/06/05         Yang Chenglong   添加login中间件
 */

package handler

import (
	"net/http"

	"hypercube/libs/message"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

const (
	TokenIndex = 0
	BearerSize = 7
)

func GetUser(user *jwt.Token) (interface{}, error) {
	var u message.User

	claims := user.Claims.(jwt.MapClaims)
	u.Token = claims["token"].(string)
	u.UserID = claims["uid"].(string)

	return u, nil
}

func LoginMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header["Authorization"]

		t := token[TokenIndex][BearerSize:]

		u, err := jwt.Parse(t, GetUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Token Parse Err")
		}

		c.JSON(http.StatusOK, u.Claims)

		return next(c)
	}
}
