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
 *	    Modify: 2017/04/15 		  Yusan Kurban
 *			#66 修改常量名和结构名
 *
 */

package main

import (
	"hypercube/proto/api"
	"hypercube/common/mq"
	"errors"
	"fmt"
)

const (
	invalideUserID = 0

	addUser    = uint8(0)
	removeUser = uint8(1)
	queryUser  = uint8(2)
)

var (
	OnLineUserMag *OnlineUserManager

	ParamErr     = errors.New("Your input parametric error!")
)

func init() {
	OnLineUserMag = &OnlineUserManager{
		users:  make(map[uint64]*api.UsrInfo),
		access: make(map[string]mq.Requester),
		req:    make(chan interface{}),
		rep:    make(chan reply),
	}

	OnLineUserMag.loop()
}

type OnlineUserManager struct {
	users  map[uint64]*api.UsrInfo
	access map[string]mq.Requester
	req    chan interface{}
	rep    chan reply
}

type userEntry struct {
	UserID       uint64
	ServerIP     string
	Type         uint8
}

type reply struct {
	ServerIP     string
	Err          error
}

func (this *OnlineUserManager) Add(user *api.UserLogin) error {
	if user.ServerIP != "" && user.UserID != invalideUserID {
		usrb := userEntry{user.UserID, user.ServerIP, addUser}

		this.req <- &usrb
		repl := <-this.rep

		return repl.Err
	}

	return ParamErr
}

func (this *OnlineUserManager) Remove(uid uint64) error {
	if uid != invalideUserID {
		user := userEntry{uid, "", removeUser}

		this.req <- &user
		repl := <-this.rep

		return repl.Err
	}

	return ParamErr
}

func (this *OnlineUserManager) Query(uid uint64) (string, error) {
	if uid != invalideUserID {
		user := userEntry{uid, "", queryUser}

		this.req <- &user
		repl := <-this.rep

		return repl.ServerIP, repl.Err
	}

	return "", ParamErr
}

func (this *OnlineUserManager) AddAccess(access *api.Access) error {
	if *access.ServerIp != "" && *access.Subject != "" {
		this.req <- access
		repl := <- this.rep

		return repl.Err
	}

	return ParamErr
}

func (this *OnlineUserManager)loop() {
	go func() {
		for {
			repl := reply{
				ServerIP:    "",
				Err:         nil,
			}
			myinterface := <-this.req

			if user, ok := myinterface.(*userEntry); ok {
				switch {
				case user.Type == addUser:
					usrlogic := api.UsrInfo{UserID: user.UserID, ServerIP: user.ServerIP}
					this.users[user.UserID] = &usrlogic

					this.rep <- repl
				case user.Type == removeUser:
					if _, ok := this.users[user.UserID]; ok {
						delete(this.users, user.UserID)

						this.rep <- repl
					} else {
						repl.Err = ParamErr
						this.rep <- repl
					}
				case user.Type == queryUser:
					if userb, ok := this.users[user.UserID]; ok {
						repl.ServerIP = userb.ServerIP

						this.rep <- repl
					} else {
						repl.Err = ParamErr
						this.rep <- repl
					}
				}
			}

			if access, ok := myinterface.(*api.Access); ok {
				req := createAccessRPC(access.Subject)
				if req != nil {
					this.access[*access.ServerIp] = req

					this.rep <- repl
				} else {
					repl.Err = ParamErr
					this.rep <- repl
				}
			}
		}
	}()
}

func (this *OnlineUserManager) PrintDebugInfo()  {
	logger.Debug(fmt.Sprintf("Online user manager:(%+v, %+v)", this.users, this.access))
}
