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
	. "hypercube/proto/general"
)

var OnLineUser *OnLineManager

type OnLineManager struct {
	locker 		*sync.Mutex
	onLineMap	map[string]UserConnect
}

func init() {
	OnLineUser = &OnLineManager{
		locker: &sync.Mutex{},
		onLineMap: map[string]UserConnect{},
	}
}

func (this *OnLineManager) OnConnect(userID string, user UserConnect) {
	this.locker.Lock()
	defer this.locker.Unlock()

	this.onLineMap[userID] = user
}

func (this *OnLineManager) OnDisconnect(userID string) {
	this.locker.Lock()
	defer this.locker.Unlock()

	if _, ok := this.onLineMap[userID]; ok {
		delete(this.onLineMap, userID)
	}
}

func (this *OnLineManager) IsUserOnline(userID string) (UserConnect, bool) {
	this.locker.Lock()
	defer this.locker.Unlock()

	user, ok := this.onLineMap[userID]
	return user, ok
}
