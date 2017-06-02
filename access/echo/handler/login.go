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
 */

package handler

import (
	"net/http"
	"github.com/labstack/echo"
	"hypercube/middleware/session"
	"hypercube/proto/general"
)

type User struct {
	Name        string  `json:"name"    form:"name"   query:"name"`
	MDUserID    string  `json:"mgo_id"  form:"mgo_id" query:"mgo_id"`
}

func Login(c echo.Context) error {
	var (
		userid general.UserKey = general.UserKey{}
		user User                  = User{}
	)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": 1})
	}

	session := session.GetSession(c)

	if session == nil {
		//logger.Error("Server session error")
		return c.JSON(http.StatusOK, map[string]interface{}{"status": 0})
	}

	userID := session.Get("user")

	if userID == nil {
		//logger.Debug("Not login, try login")
	}

	session.Set("user", user)
	session.Save()

	userid.MySQLUserID = 0
	userid.MDUserID = session.Get("user").(User).MDUserID

	return c.JSON(http.StatusOK, map[string]interface{}{"status": 0})
}

func Logout(c echo.Context) error {
	session := session.GetSession(c)

	userID := session.Get("id")

	if userID == nil {
		//logger.Debug("Not login, exit")
		return c.JSON(http.StatusOK, map[string]interface{}{"status": 1})
	}

	session.Clear()
	session.Save()

	return c.JSON(http.StatusOK, map[string]interface{}{"status": 0})
}
