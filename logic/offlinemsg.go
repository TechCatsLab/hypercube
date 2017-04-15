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

/**
 * Created by HeChengJun on 13/04/2017.
 */

package main

import (
    "hypercube/common/workq"
    "hypercube/proto/general"
    "hypercube/common/container"
)
var(
    msgbuf map[uint64]*container.Ring
)

const (
    maxMsgSize = 100
)

func init()  {
    msgbuf = make(map[uint64]*container.Ring)
    initSendMessageQueue()
}

func addHistMessage(userID uint64, msg interface{})  {
    if msgbuf[userID] == nil {
        msgbuf[userID] = container.NewRing(maxMsgSize)
    }

    mesg := msg.(*general.Message)

    if msgbuf[userID].Full() == false {
        msgbuf[userID].Push(mesg)
    }
}

func getHistMessages(userID uint64) *container.Ring {
    return msgbuf[userID]
}

func clearHistMessages(userID uint64) {
    msgbuf[userID].MPop(msgbuf[userID].Len())
}

const (
    workersCount = 128
)

var (
    sendWorkQueue *workq.Dispatcher
)

func userSendMessageHandler(userID uint64) error {
    var (
        message *pushMessageJob
    )

    messages := getHistMessages(userID).GetAll()

    for _, mes := range messages {
        message = &pushMessageJob{
            message: mes.(*general.Message),
        }

        if err := appendPushMessage(message); err != nil {
            return err
        }
    }

    clearHistMessages(userID)

    return nil
}

type pushMessageJob struct {
    message *general.Message
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

func appendPushMessage(msg *pushMessageJob) error {
    logger.Debug("push to job queue...", msg.message.From, "->", msg.message.To)

    sendWorkQueue.PushToJobQ(msg)

    return nil
}
