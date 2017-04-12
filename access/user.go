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
 *     Initial: 2017/04/05        Yusan Kurban
 */

package main

import (
	"strings"
	"sync"
	"github.com/gorilla/websocket"
	"hypercube/proto/api"
	"time"
	"encoding/json"
)

var OnLineUser *OnLineManager

type OnLineManager struct {
	locker 		sync.RWMutex
	onLineMap	map[uint64]*websocket.Conn
}

func init() {
	OnLineUser = &OnLineManager{
		onLineMap: map[uint64]*websocket.Conn{},
	}
}

func (this *OnLineManager) OnConnect(userID uint64, conn *websocket.Conn) {
	this.locker.Lock()
	defer this.locker.Unlock()

	this.onLineMap[userID] = conn
}

func (this *OnLineManager) OnDisconnect(userID uint64) {
	this.locker.Lock()
	defer this.locker.Unlock()

	if _, ok := this.onLineMap[userID]; ok {
		delete(this.onLineMap, userID)
	}
}

func (this *OnLineManager) IsUserOnline(userID uint64) (*websocket.Conn, bool) {
	this.locker.RLock()
	defer this.locker.Unlock()

	user, ok := this.onLineMap[userID]
	return user, ok
}

func (this *OnLineManager) OnUnusualDisConnect(conn *websocket.Conn) (uint64, bool) {
	this.locker.RLock()
	defer this.locker.Unlock()

	for k, v := range this.onLineMap {
		if v == conn {
			return k, true
		}
	}

	return 0, false
}

func (this *OnLineManager) UserLoginHandler(userID uint64) error {
	var (
		userlog api.UserLogin
		proto 	*api.Request
		conv 	[]byte
		r		api.Reply
		err		error
	)

	userlog = api.UserLogin{
		UserID: 	userID,
		ServerIP:   strings.Split(configuration.Addrs[0], ":")[0],
	}

	conv, err = json.Marshal(&userlog)
	if err != nil {
		logger.Error(err)
	}

	proto = &api.Request{
		Type:    api.ApiTypeUserLogin,
		Content: conv,
	}

	err = logicRequester.Request(proto, &r, time.Duration(100) * time.Millisecond)

	if err != nil {
		logger.Error(err)
	}

	if r.Code != api.ErrSucceed {
		logger.Error("request error code:", r.Code)
	}

	return err
}

func (this *OnLineManager) UserLogoutHandler(userID uint64) error {
	var (
		userlog api.UserLogout
		proto 	*api.Request
		conv    []byte
		r		api.Reply
		err		error
	)

	userlog = api.UserLogout{
		UserID: 	userID,
	}

	conv, err = json.Marshal(&userlog)
	if err != nil {
		logger.Error(err)
	}

	proto = &api.Request{
		Type:    api.ApiTypeUserLogin,
		Content: conv,
	}

	err = logicRequester.Request(proto, &r, time.Duration(100) * time.Millisecond)
	if err != nil {
		logger.Error(err)
	}

	return err
}

func (this *OnLineManager) PrintDebugInfo() {
	this.locker.RLock()
	defer this.locker.RUnlock()

	logger.Debug("Online user manager:", this.onLineMap)
}
