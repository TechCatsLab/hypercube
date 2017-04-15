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
 *     Initial: 2017/04/12        Feng Yifei
 */

package main

import (
	"encoding/json"
	"errors"

	"hypercube/common/workq"
	"hypercube/proto/api"
	"hypercube/proto/general"
)

const (
	workersCount = 128
)

var (
	pushWorkQueue *workq.Dispatcher
)

func userMessageHandler(p interface{}, req interface{}) interface{} {
	var (
		message *pushMessageJob
	)

	message = &pushMessageJob{
		message: req.(*general.Message),
	}

	return appendPushMessage(message)
}

type pushMessageJob struct {
	message *general.Message
}

func (this *pushMessageJob) sendToLogic() error {
	msg, err := json.Marshal(this.message)

	if err != nil {
		return err
	}

	return logicRequester.SendMessage(&api.Request{
		Type:    general.AccTypeUTUMsg,
		Content: msg,
	})
}

func (this *pushMessageJob) Do() error {
	conn, ok := OnLineUser.IsUserOnline(this.message.To)

	if !ok {
		if !this.message.Pushed {
			logger.Debug("User not on this server, sending to logic:", this.message.From, "->", this.message.To)
			return this.sendToLogic()
		}
		return errors.New("user is offline")
	}

	logger.Debug("Sending:", this.message.From, "->", this.message.To)

	conn.Mutex.Lock()
	defer conn.Mutex.Unlock()

	return conn.Conn.WriteJSON(this.message)
}

func initPushMessageQueue() {
	pushWorkQueue = workq.NewDispatcher(workersCount)
	pushWorkQueue.Run()
	logger.Debug("message queue is running")
}

func appendPushMessage(msg *pushMessageJob) error {
	logger.Debug("push to job queue...", msg.message.From, "->", msg.message.To)

	pushWorkQueue.PushToJobQ(msg)

	return nil
}
