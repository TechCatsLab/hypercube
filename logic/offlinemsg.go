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
 *     Initial: 2017/04/13        He ChengJun
 *      Modify: 2017/04/15        Yang ChengLong
 *          #60 使用环形队列存储消息
 *
 */

package main

import (
	"hypercube/common/workq"
	"hypercube/proto/general"
	"hypercube/common/container"
)
var(
    msgbuf map[general.UserKey]*container.Ring
)

const (
    maxMsgSize = 100
)

func init()  {
    msgbuf = make(map[general.UserKey]*container.Ring)
    initSendMessageQueue()
}

func addHistMessage(userID general.UserKey, msg interface{})  {
    if msgbuf[userID] == nil {
        msgbuf[userID] = container.NewRing(maxMsgSize)
    }

    mesg := msg.(*general.Message)

    if msgbuf[userID].Full() == false {
        msgbuf[userID].Push(mesg)
    }
}

func getHistMessages(userID general.UserKey) (interface{}, error) {
    return msgbuf[userID].Pop()
}

func clearHistMessages(userID general.UserKey) {
    msgbuf[userID].MPop(msgbuf[userID].Len())
}

const (
    workersCount = 128
)

var (
    sendWorkQueue *workq.Dispatcher
)

func userSendMessageHandler(userID general.UserKey) error {
    var (
	    message  *pushMessageJob
    )
	num := msgbuf[userID].Len()

    for i := 0; i < num; i++ {
	    mes,err := getHistMessages(userID)

	    if err != nil {
		    logger.Error("User get history err:", err)

		    return err
	    }

	    message = &pushMessageJob{
		    message: mes.(*general.Message),
	    }
	    message.Send(sendWorkQueue)
	    num = msgbuf[userID].Len()
    }
    return nil
}

type pushMessageJob struct {
    message *general.Message
}

func (this *pushMessageJob) Send(sendWorkQueue *workq.Dispatcher) {
	sendWorkQueue.PushToJobQ(this)
}

func (this *pushMessageJob) Do() error {
    serverip,err := OnLineUserMag.Query(this.message.To)
    if err != nil {
        return err
    }

    if req, ok := OnLineUserMag.access[serverip]; ok {
        req.SendMessage(this.message)
    }

    logger.Debug("Sending:", this.message.From, "->", this.message.To)

    return nil
}

func initSendMessageQueue() {
    sendWorkQueue = workq.NewDispatcher(workersCount)
    sendWorkQueue.Run()
    logger.Debug("message queue is running")
}
