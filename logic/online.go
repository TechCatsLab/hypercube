/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
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
	"sync"

	"github.com/TechCatsLab/hypercube/libs/log"
	"github.com/TechCatsLab/hypercube/libs/message"
)

var (
	onlineUserManager *UserManager

	ParamErr = errors.New("The Parametric error!")
)

func initOnline() {
	onlineUserManager = &UserManager{
		users: sync.Map{},
	}
}

type UserManager struct {
	mux   sync.Mutex
	users sync.Map
}

func (this *UserManager) Add(user message.UserEntry) error {
	if user.ServerIP.ServerIP != "" && user.UserID.UserID != "" {
		this.mux.Lock()
		defer this.mux.Unlock()

		if _, exists := this.users.Load(user.UserID); exists {
			log.Logger.Warn("user already login!")
		}

		this.users.Store(user.UserID, &user.ServerIP)
		return nil
	}

	return ParamErr
}

func (this *UserManager) Remove(user message.UserEntry) error {
	if user.ServerIP.ServerIP != "" && user.UserID.UserID != "" {
		this.mux.Lock()
		defer this.mux.Unlock()

		if _, exists := this.users.Load(user.UserID); !exists {
			log.Logger.Warn("user hasn't login!")
		}

		this.users.Delete(user.UserID)
		return nil
	}

	return ParamErr
}

func (this *UserManager) Query(user message.User) (*message.Access, bool) {
	this.mux.Lock()
	defer this.mux.Unlock()

	serverIP, ok := this.users.Load(user)
	if !ok {
		return nil, false
	}

	return serverIP.(*message.Access), true
}

func (this *UserManager) PrintDebugInfo() {
	log.Logger.Debug(fmt.Sprintf("Online user manager: %v", this.users))
}
