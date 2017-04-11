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
	"errors"
)

const (
	AddType = uint8(0)
	RmType  = uint8(1)
	QrType  = uint8(2)
)

var (
	OnLineUserMag *OnLineUserMagServer

	ParamErr     = errors.New("Your input parametric error!")
	UserNotExist = errors.New("User not exist!")
)

func init() {
	OnLineUserMag = &OnLineUserMagServer{
		users:      make(map[uint64]api.UsrInfo),
		reqchan:    make(chan usrbasic),
		replychan:  make(chan chanreply),
	}

	OnLineUserMag.loop()
}

type OnLineUserMagServer struct {
	users       map[uint64]api.UsrInfo
	reqchan     chan usrbasic
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

func (this *OnLineUserMagServer) Add(user api.UsrInfo) error {
	if user.ServerIP != "" && user.UserID != 0 {
		usrb := usrbasic{user.UserID, user.ServerIP, AddType}

		this.reqchan <- usrb
		repl := <-this.replychan

		return repl.Err
	}

	return ParamErr
}

func (this *OnLineUserMagServer) Remove(uid uint64) error {
	if uid != 0 {
		user := usrbasic{uid, "", RmType}

		this.reqchan <- user
		repl := <-this.replychan

		return repl.Err
	}

	return ParamErr
}

func (this *OnLineUserMagServer) Query(uid uint64) (string, error) {
	if uid != 0 {
		user := usrbasic{uid, "", QrType}

		this.reqchan <- user
		repl := <-this.replychan

		return repl.ServerIP, repl.Err
	}

	return "", ParamErr
}

func (this *OnLineUserMagServer)loop() {
	go func() {
		for {
			repl := chanreply{
				ServerIP:    "",
				Err:         nil,
			}
			user := <-this.reqchan
			switch {
			case user.Type == AddType:
				usrlogic := api.UsrInfo{user.UserID, user.ServerIP}
				this.users[user.UserID] = usrlogic

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
				if serverip, ok := this.users[user.UserID]; ok {
					repl.ServerIP = serverip

					this.replychan <- repl
				} else {
					repl.Err = ParamErr
					this.replychan <- repl
				}
			}
		}
	}()
}
