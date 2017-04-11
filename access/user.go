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
	"sync"
	"github.com/gorilla/websocket"
	. "hypercube/proto/api"
	"time"
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

func (this *OnLineManager) SendToLogic(userID uint64) error {
	var (
		userlog UserLogin
		r		bool
		err		error
	)

	userlog = UserLogin{
		UserID: 	userID,
		ServerIP:   configuration.Addrs[0],
	}
	err = requester.Request(userlog, r, time.Duration(5)*time.Second)

	return err
}