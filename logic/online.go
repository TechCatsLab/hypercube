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
 *     Initial: 2017/04/04        Liu Jiachang
 */

package main

import (
	"hypercube/proto/api"
	"hypercube/common/mq"
	"errors"
)

const (
	invalideUserID = 0

	AddType = uint8(0)
	RmType  = uint8(1)
	QrType  = uint8(2)
)

var (
	OnLineUserMag *OnlineUserManager

	ParamErr     = errors.New("Your input parametric error!")
	UserNotExist = errors.New("User not exist!")
)

func init() {
	OnLineUserMag = &OnlineUserManager{
		users:      make(map[uint64]*api.UsrInfo),
		access:     make(map[string]mq.Requester),
		reqchan:    make(chan interface{}),
		replychan:  make(chan chanreply),
	}

	OnLineUserMag.loop()
}

type OnlineUserManager struct {
	users       map[uint64]*api.UsrInfo
	access      map[string]mq.Requester
	reqchan     chan interface{}
	replychan   chan chanreply
}

type usrbasic struct {
	UserID       uint64
	ServerIP     string
	Type         uint8
}

type chanreply struct {
	ServerIP     string
	Err          error
}

func (this *OnlineUserManager) Add(user *api.UserLogin) error {
	if user.ServerIP != "" && user.UserID == invalideUserID {
		usrb := usrbasic{user.UserID, user.ServerIP, AddType}

		this.reqchan <- &usrb
		repl := <-this.replychan

		return repl.Err
	}

	return ParamErr
}

func (this *OnlineUserManager) Remove(uid uint64) error {
	if uid != 0 {
		user := usrbasic{uid, "", RmType}

		this.reqchan <- &user
		repl := <-this.replychan

		return repl.Err
	}

	return ParamErr
}

func (this *OnlineUserManager) Query(uid uint64) (string, error) {
	if uid != 0 {
		user := usrbasic{uid, "", QrType}

		this.reqchan <- &user
		repl := <-this.replychan

		return repl.ServerIP, repl.Err
	}

	return "", ParamErr
}

func (this *OnlineUserManager) AddAccess(access *api.Access) error {
	if *access.ServerIp != "" && *access.Subject != "" {
		this.reqchan <- access
		repl := <- this.replychan

		return repl.Err
	}

	return ParamErr
}

func (this *OnlineUserManager)loop() {
	go func() {
		for {
			repl := chanreply{
				ServerIP:    "",
				Err:         nil,
			}
			myinterface := <-this.reqchan

			if user, ok := myinterface.(*usrbasic); ok {
				switch {
				case user.Type == AddType:
					usrlogic := api.UsrInfo{UserID: user.UserID, ServerIP: user.ServerIP}
					this.users[user.UserID] = &usrlogic

					this.replychan <- repl
				case user.Type == RmType:
					if _, ok := this.users[user.UserID]; ok {
						delete(this.users, user.UserID)

						this.replychan <- repl
					} else {
						repl.Err = ParamErr
						this.replychan <- repl
					}
				case user.Type == QrType:
					if userb, ok := this.users[user.UserID]; ok {
						repl.ServerIP = userb.ServerIP

						this.replychan <- repl
					} else {
						repl.Err = ParamErr
						this.replychan <- repl
					}
				}
			}

			if access, ok := myinterface.(*api.Access); ok {
				req := createAccessRPC(access.Subject)
				if req != nil {
					this.access[*access.ServerIp] = req

					this.replychan <- repl
				} else {
					repl.Err = ParamErr
					this.replychan <- repl
				}
			}
		}
	}()
}

func (this *OnlineUserManager) PrintDebugInfo()  {
	logger.Debug("Online user manager:", this.users, this.access)
}
