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
 *     Initial: 2017/07/10        Yang Chenglong
 */

package main

import (
	"errors"
	"fmt"
	"net/rpc"
	"sync"

	"hypercube/libs/log"
	"hypercube/libs/message"
)

const (
	Success int = 1
	Failed  int = 0
)

var (
	OnLineUserMag *OnlineUserManager

	ParamErr = errors.New("The Parametric error!")
)

func init() {
	OnLineUserMag = &OnlineUserManager{
		users: make(map[message.User]*message.Access),
	}

	rpc.Register(OnLineUserMag)
	rpc.HandleHTTP()
}

type OnlineUserManager struct {
	mux   sync.Mutex
	users map[message.User]*message.Access
}

func (this *OnlineUserManager) Add(user message.UserEntry, reply *int) error {
	if user.ServerIP.ServerIp != "" && user.UserID.UserID != "" {
		this.mux.Lock()
		defer this.mux.Unlock()

		if _, exists := this.users[user.UserID]; exists {
			log.Logger.Warn("user already login!")
		}

		this.users[user.UserID] = &user.ServerIP
		*reply = Success
	}
	*reply = Failed

	return ParamErr
}

func (this *OnlineUserManager) Remove(user message.UserEntry, reply *int) error {
	if user.ServerIP.ServerIp != "" && user.UserID.UserID != "" {
		this.mux.Lock()
		defer this.mux.Unlock()

		if _, exists := this.users[user.UserID]; !exists {
			log.Logger.Warn("user hasn't login!")
		}

		delete(this.users, user.UserID)
		*reply = Success
	}
	*reply = Failed

	return ParamErr
}

func (this *OnlineUserManager) Query(user message.User) (string, bool) {
	this.mux.Lock()
	defer this.mux.Unlock()

	serverIP, ok := this.users[user]

	return serverIP.ServerIp, ok
}

func (this *OnlineUserManager) PrintDebugInfo() {
	log.Logger.Debug(fmt.Sprintf("Online user manager: %+v", this.users))
}
