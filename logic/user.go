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
 *     Modify:  2017/06/07        Yang Chenglong  添加AccessHeartHandler
 */

package main

import (
	"hypercube/libs/log"
	"hypercube/proto/types"
	"hypercube/proto/general"
)

func userLoginRequestHandler(req interface{}) interface{} {
	var (
		login       *general.UserLogin
		reply       *general.Reply
	)

	login = req.(*general.UserLogin)

	log.GlobalLogger.Debug("userLoginRequestHandler userID:", login.UserID)

	err := OnLineUserMag.Add(login)

	if err != nil {
		log.GlobalLogger.Error("userLoginRequestHandler-->login err:", err)

		reply = &general.Reply{
			Code: types.ErrLogin,
		}
	} else {
		reply = &general.Reply{
			Code: types.ErrSucceed,
		}
	}

	err = userSendMessageHandler(login.UserID)
	if err != nil {
		log.GlobalLogger.Error("userLoginRequestHandler-->userSendMessageHandler:", err)
	}

	return reply
}

func userLogoutRequestHandler(req interface{}) interface{} {
	var (
		reply        *general.Reply
		logout       *general.UserLogout
	)

	logout = req.(*general.UserLogout)
	log.GlobalLogger.Debug("userLogoutRequestHandler userID:", logout.UserID)

	err := OnLineUserMag.Remove(logout.UserID)

	if err != nil {
		log.GlobalLogger.Error("userLogoutRequestHandler--> err logout:", err)

		reply = &general.Reply{
			Code: types.ErrLogout,
		}
	} else {
		reply = &general.Reply{
			Code: types.ErrSucceed,
		}
	}

	return reply
}

func MessageHandler(req interface{}) interface{} {
	var (
		msg          *general.Message
		reply        *general.Reply
	)

	msg = req.(*general.Message)
	log.GlobalLogger.Debug("userToUserMsgHandler:", *msg)

	accessip, err := OnLineUserMag.Query(msg.To)

	if err != nil {
		addHistMessage(msg.To, msg)

		log.GlobalLogger.Error("MessageHandler-->OnLineUserMag.Query():", err)

		reply = &general.Reply{
			Code: types.ErrUserQuery,
		}

		return reply
	}

	if req, ok := OnLineUserMag.access[accessip]; ok {
		err = req.SendMessage(msg)

		if err != nil {
			log.GlobalLogger.Error("MessageHandler-->SendMessage:", err)

			reply = &general.Reply{
				Code: types.ErrSendToAccess,
			}
		} else {
			log.GlobalLogger.Debug("SendMessage Succeed")
			reply = &general.Reply{
				Code: types.ErrSucceed,
			}
		}
	} else {
		reply = &general.Reply{
			Code: types.ErrFindAccess,
		}
	}

	return reply
}

func AccessConnectHandler(req interface{}) interface{} {
	var (
		access          *general.Access
		reply           *general.Reply
	)

	access = req.(*general.Access)
	log.GlobalLogger.Debug("userToUserMsgHandler:", access)

	err := OnLineUserMag.AddAccess(access)

	if err != nil {
		log.GlobalLogger.Error("AccessConnectHandler-->AddAccess:", err)

		reply = &general.Reply{
			Code: types.ErrAddAccess,
		}
	}

	return reply
}

func AccessHeartHandler(req interface{}) interface{} {
	var (
		heart          *general.Proto
		reply          *general.Reply
	)

	heart = req.(*general.Proto)

	if _, err := heart.VerCheck(); err != types.ErrSucceed {
		log.GlobalLogger.Error("keepAliveRequestHandler vercheck:", err)
		return err
	}

	log.GlobalLogger.Debug("a heartbeat from access which versions is:", heart.Ver)

	return reply
}
