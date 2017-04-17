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
	"time"
	"encoding/json"

	"github.com/gorilla/websocket"

	"hypercube/proto/api"
	"hypercube/proto/general"
)

var OnLineManagement *OnLineTable

type onlineEntry struct {
	mutex sync.Mutex
	conn  *websocket.Conn
}

type OnLineTable struct {
	locker 		sync.RWMutex
	onLineMap	map[general.UserKey]*onlineEntry
}

func init() {
	OnLineManagement = &OnLineTable{
		onLineMap: map[general.UserKey]*onlineEntry{},
	}
}

// 用户登录时，保存用户id 和 connection到在线用户表，并把用户id 传到逻辑层
func (this *OnLineTable) OnConnect(userID general.UserKey, conn *websocket.Conn) error {
	var err	error

	this.locker.Lock()

	this.onLineMap[userID] = &onlineEntry{conn: conn}

	this.locker.Unlock()

	err = OnLineManagement.loginReport(userID)
	if err != nil {
		logger.Error("User OnConnect error:",err)
	}

	return err
}

// 用户登出时，把用户 id 和 connection 从在线用户表中删除，并通知逻辑层
func (this *OnLineTable) OnDisconnect(userID general.UserKey) error {
	var err error

	this.locker.Lock()

	if _, ok := this.onLineMap[userID]; ok {
		delete(this.onLineMap, userID)
	}

	this.locker.Unlock()
	
	err = OnLineManagement.logoutReport(userID)
	if err != nil {
		logger.Error("User OnDisonnect error",err)
	}

	return err
}

// 通过用户 id 获取对应的 connection
func (this *OnLineTable) GetEntryByID(userID general.UserKey) (*onlineEntry, bool) {
	this.locker.RLock()
	defer this.locker.RUnlock()

	user, ok := this.onLineMap[userID]
	return user, ok
}

// 通过 connection 获取对应的用户id
func (this *OnLineTable) GetIDByConnection(conn *websocket.Conn) (general.UserKey, bool) {
	var user = general.UserKey{}
	this.locker.RLock()
	defer this.locker.RUnlock()

	for k, v := range this.onLineMap {
		if v.conn == conn {
			return k, true
		}
	}

	return user, false
}

// 用户登录时，通知逻辑层并把用户 id 传到逻辑层，并作出错误处理
func (this *OnLineTable) loginReport(userID general.UserKey) error {
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
		logger.Error("UserLoginHandler marshall error:", err)
	}

	proto = &api.Request{
		Type:    api.ApiTypeUserLogin,
		Content: conv,
	}

	err = logicRequester.Request(proto, &r, time.Duration(100) * time.Millisecond)

	if err != nil {
		logger.Error("UserLoginHandler request receive error:", err)
	}

	if r.Code != api.ErrSucceed {
		logger.Error("request error code:", r.Code)
	}

	return err
}

// 用户登出时，通知逻辑层并把用户id 传过去，并作出错误处理
func (this *OnLineTable) logoutReport(userID general.UserKey) error {
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
		logger.Error("UserLogoutHandler format request error:", err)
	}

	proto = &api.Request{
		Type:    api.ApiTypeUserLogin,
		Content: conv,
	}

	err = logicRequester.Request(proto, &r, time.Duration(100) * time.Millisecond)
	if err != nil {
		logger.Error("UserLogoutHandler send request error:", err)
	}

	return err
}

// 打印出在线用户表
func (this *OnLineTable) PrintDebugInfo() {
	this.locker.RLock()
	defer this.locker.RUnlock()

	logger.Debug("Online user manager:", this.onLineMap)
}
