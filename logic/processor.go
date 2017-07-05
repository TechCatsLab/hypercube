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
 *     Modify:  2017/06/07        Yang Chenglong       添加AccessHeart处理
 */

package main

import (
	"hypercube/libs/log"
	"hypercube/message/types"
	"hypercube/message/general"
	"encoding/json"
)

type handlerFunc func(interface{}) interface{}

func requestProcessor(req []byte) interface{} {
	var (
		err          error
		request      general.Request
		login        general.UserLogin
		logout       general.UserLogout
		msg          general.Message
		access       general.Access
		v            interface{}
		handler      handlerFunc
		heart        general.Proto
	)

	err = json.Unmarshal(req, &request)

	if err != nil {
		log.GlobalLogger.Error("Logic request processor error:", err)
		return nil
	}

	log.GlobalLogger.Debug("Logic RPC received message type:", request.Type)

	switch request.Type {
	case types.ApiTypeUserOnConnect:
	case types.ApiTypeUserOnDisConnect:
	case types.ApiTypeUserLogin:
		v = &login
		handler = userLoginRequestHandler
	case types.ApiTypeUserLogout:
		v = &logout
		handler = userLogoutRequestHandler
	case types.GeneralTypeTextMsg:
		v = &msg
		handler = MessageHandler
	case types.ApiTypeAccessInfo:
		v = &access
		handler = AccessConnectHandler
	case types.ApiTypeAccessHeart:
		v = &heart
		handler = AccessHeartHandler
	}

	if v != nil {
		err = json.Unmarshal(request.Content, v)
		if err != nil {
			log.GlobalLogger.Error("Logic request processor content error:", err)
			return nil
		}

		return handler(v)
	}

	return nil
}
