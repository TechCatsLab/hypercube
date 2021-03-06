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
 *     Initial: 2017/07/11        Sun Anxiang
 *     Modify : 2017/07/28        Yang Chenglong
 */

package main

import (
	"encoding/json"

	"gopkg.in/mgo.v2"

	"github.com/TechCatsLab/hypercube/libs/connector"
	"github.com/TechCatsLab/hypercube/libs/connector/mongo"
	"github.com/TechCatsLab/hypercube/libs/log"
	"github.com/TechCatsLab/hypercube/libs/message"
	rp "github.com/TechCatsLab/hypercube/libs/rpc"
)

type OfflineMessage chan message.UserEntry

var (
	Queue    chan *message.Message
	Shutdown chan struct{}
	offline  OfflineMessage
	conn     connector.Connector = &mongo.MgoConnector{}
)

func initQueue() {
	Queue = make(chan *message.Message, 100)
	Shutdown = make(chan struct{})
	offline = make(chan message.UserEntry)

	conn.Initialize()
	QueueStart()
}

func QueueStart() {
	go func() {
		for {
			select {
			case msg := <-Queue:
				HandleMessage(msg)
			case user := <-offline:
				OfflineMessageHandler(user)
			case <-Shutdown:
				close(Queue)
				close(offline)
				return
			}
		}
	}()
}

func OfflineMessageHandler(user message.UserEntry) error {
	messages, err := conn.Get(user.UserID.UserID, message.MessageUnsent)
	if err != nil {
		if err == mgo.ErrNotFound {
			log.Logger.Debug("User doesn't have offline messages!")
			return nil
		}

		log.Logger.Error("GetOffLineMessage Error %v", err)
		return err
	}

	for _, msg := range messages {
		switch msg.Type {
		case message.MessageTypePlainText, message.MessageTypeEmotion:
			content := message.PlainText{
				From:    message.User{UserID: msg.From},
				To:      message.User{UserID: msg.To},
				Content: msg.Content,
			}

			text, err := json.Marshal(content)
			if err != nil {
				log.Logger.Error("OfflineMessage Marshal Error %v", err)
				return err
			}

			m := &message.Message{
				Type:    msg.Type,
				Content: text,
			}

			id := msg.Messageid
			status := TransmitMsg(m)
			if status == message.MessageSent {
				err = conn.Update(id.Hex(), message.MessageUnsent)
				if err != nil {
					log.Logger.Error("ModifyMessageStatus error:", err)
					ShutDown()
				}
			}
		}
	}

	return nil
}

func HandleMessage(msg *message.Message) {
	switch msg.Type {
	case message.MessageTypePlainText, message.MessageTypeEmotion:
		status := TransmitMsg(msg)

		HandlePlainText(msg, status)
	default:
		log.Logger.Debug("Not recognized message type!")
	}
}

func TransmitMsg(msg *message.Message) int {
	var plainText message.PlainText

	err := json.Unmarshal(msg.Content, &plainText)
	if err != nil {
		log.Logger.Error("TransmitMsg Unmarshal Error %v", err)

		return message.MessageUnsent
	}

	serverIP, isOnline := onlineUserManager.Query(plainText.To)
	if isOnline {
		op := rp.Options{
			Proto: "tcp",
			Addr:  serverIP.ServerIP,
		}

		err := Send(plainText.To, *msg, op)
		if err != nil {
			log.Logger.Error("TransmitMsg Send Error %v", err)

			return message.MessageUnsent
		}

		return message.MessageSent
	}

	return message.MessageUnsent
}

func HandlePlainText(msg *message.Message, status int) {
	err := conn.Put(msg, status)

	if err != nil {
		log.Logger.Error("Insert into message error:", err)

		// TODO: 消息已经发送成功，这里只是存入数据库时失败。
		// TODO：解决思路：先把未存入数据库的消息放入缓存，再统一存入数据库。
		return
	}
}

func ShutDown() {
	Shutdown <- struct{}{}
}
