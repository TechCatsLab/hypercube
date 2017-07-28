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
 *     Initial: 2017/07/11        Jia Chenhui
 */

package main

import (
	"hypercube/libs/log"
	"hypercube/libs/message"
	"hypercube/libs/rpc"
	"encoding/json"

	"github.com/jinzhu/gorm"

	"hypercube/orm/cockroach"
	db "hypercube/model"
)

type LogicRPC struct{

}

var (
	options []rpc.Options
	clients *rpc.Clients
	logic   *LogicRPC
)

func initRPC() {
	logic = new(LogicRPC)

	for _, addr := range configuration.AccessAddrs {
		op := rpc.Options{
			Proto: "tcp",
			Addr:  addr,
		}

		options = append(options, op)
	}

	clients = rpc.Dials(options)
}

func (this *LogicRPC) LoginHandler(user message.UserEntry, reply *int) error {
	err := onLineUserMag.Add(user)
	if err != nil {
		log.Logger.Error("LoginHandle Add Error %+v: ", err)
		*reply = message.ReplyFailed
		return err
	}

	offline <- user
	onLineUserMag.PrintDebugInfo()
	*reply = message.ReplySucceed
	return nil
}

func OfflineMessageHandler(user message.UserEntry) error {
	conn, err := cockroach.DbConnPool.GetConnection()
	if err != nil {
		log.Logger.Error("Get cockroach connect error:", err)
		return err
	}
	defer cockroach.DbConnPool.ReleaseConnection(conn)

	mes, err := db.MessageService.GetOffLineMessage(conn, user.UserID.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Debug("")
			goto Mess
		}

		log.Logger.Error("GetOffLineMessage Error %v", err)
		return err
	}
	Mess:
	for _, msg := range mes {
		switch msg.Type {
		case message.MessageTypePlainText:
			content := message.PlainText{
				From:    message.User{UserID:msg.Source},
				To:      message.User{UserID:msg.Target},
				Content: msg.Content,
			}

			text, err := json.Marshal(content)
			if err != nil {
				log.Logger.Error("OffLineMessage Marshal Error %v", err)
				return err
			}

			mesg := &message.Message{
				Type:       msg.Type,
				Version:    msg.Version,
				Content:    text,
			}

			TransmitMsg(mesg)
		}
	}

	return nil
}

func (this *LogicRPC) LogoutHandle(user message.UserEntry, reply *int) error {
	err := onLineUserMag.Remove(user)
	if err != nil {
		log.Logger.Error("LogoutHandle Error %v", err)
		*reply = message.ReplyFailed
		return err
	}

	onLineUserMag.PrintDebugInfo()
	*reply = message.ReplySucceed
	return nil
}


func (m *LogicRPC) Add(msg *message.Message, reply *bool) error {
	Queue <- msg
	*reply = true

	return nil
}

// Send calls the function of the access layer remotely, send a message to a specific user.
