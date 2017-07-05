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
 *     Initial: 2017/04/05        Yang Chenglong
 *     Modify:  2017/06/07        Yang Chenglong   添加向logic发送心跳
 */

package main

import (
	"encoding/json"
	"time"
	"hypercube/libs/log"
	"hypercube/proto/general"
	"hypercube/proto/types"
)

func keepAliveRequestHandler(p interface{},req interface{}) interface{} {
	var (
		heart    *general.Keepalive
		pro      *general.Proto
	)

	pro = p.(*general.Proto)

	if _, err := pro.VerCheck(); err != types.ErrSucceed {
		log.GlobalLogger.Error("keepAliveRequestHandler vercheck:", err)
		return err
	}

	heart = req.(*general.Keepalive)

	log.GlobalLogger.Debug("a heartbeat from UserId:", heart.Uid)

	return nil
}

func sendAccessHeart() {
	var (
		r 		general.Reply
		ver     general.AccessHeart
	)

	ver.Ver = general.CurVer
	v, _ := json.Marshal(&ver)


	heart   := &general.Request{
		Type:      types.ApiTypeAccessHeart,
		Content:   v,
	}

	for {
		err := logicRequester.Request(heart, &r, time.Duration(100) * time.Millisecond)

		if err != nil {
			log.GlobalLogger.Error("sendAccessHeart error:", err)
		}

		if r.Code != types.ErrSucceed {
			log.GlobalLogger.Error("request error code:", r.Code)
		}

		time.Sleep(time.Duration(configuration.AccessHeartRate) * time.Second)
	}
}
