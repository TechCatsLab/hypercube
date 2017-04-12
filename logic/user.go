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
 *     Initial: 2017/04/04        Feng Yifei
 */

package main

import (
	"hypercube/proto/api"
	"hypercube/proto/general"
)

func userLoginRequestHandler(req interface{}) interface{} {
	var (
		login    *api.UserLogin

	)

	login = req.(*api.UserLogin)
	logger.Debug("userLoginRequestHandler userID:", login.UserID)

	err := OnLineUserMag.Add(login)

	if err != nil {

	}

	return nil
}

func userLogoutRequestHandler(req interface{}) interface{} {
	var (
		reply        *api.Reply
		logout       *api.UserLogout
	)

	logout = req.(*api.UserLogout)
	logger.Debug("userLogoutRequestHandler userID:", logout.UserID)

	err := OnLineUserMag.Remove(logout.UserID)

	if err != nil {
		return err
	}

	reply = &api.Reply{
		Code: api.ErrSucceed,
	}

	return reply
}

func userToUserMsgHandler(req interface{}) interface{} {
	var (
		msg          *general.Message
	)

	msg = req.(*general.Message)
	logger.Debug("userToUserMsgHandler:", msg)

	_, err := OnLineUserMag.Query(msg.To)

	if err != nil {
		return err
	}

	return err
}
