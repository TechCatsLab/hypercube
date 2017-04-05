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
 *     Initial: 2017/04/01        HeCJ
 */

package mq

import (
    "errors"
    "time"
    "github.com/nats-io/go-nats-streaming"
    "github.com/nats-io/go-nats-streaming/pb"
)

var (
    ErrMessageEmpty = errors.New("message queue is empty ")
)

var (
    READ_TIMEOUT = time.Second * 5
)

type NatsStreaming struct {
    subject     *string
    durable		*string
    conn 		*stan.Conn
}

func ConnectToServer(urls, clusterID, clientID, subject, durable  string) (*NatsStreaming, error) {
    sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(urls))

    if err != nil {
        return nil, err
    }
    return &NatsStreaming{
        subject: &subject,
        durable: &durable,
        conn: &sc,
    }, nil
}

func (ns *NatsStreaming)WriteMessage(msg interface{}) error  {
    err := (*ns.conn).Publish(*ns.subject, []byte(msg))

    if err != nil {
        return err
    }
    return nil
}

func (ns *NatsStreaming)ReadMessage() (msg string,  err error) {

    startOpt := stan.StartAt(pb.StartPosition_NewOnly)
    readMessage := make(chan string, 1)
    var (
        sub stan.Subscription
    )

    messageHandle := func(msg *stan.Msg){
        readMessage <- string(msg.Data)
    }

    sub, err = (*ns.conn).Subscribe(*ns.subject, messageHandle, startOpt, stan.DurableName(*ns.durable))
    if err != nil {
        (*ns.conn).Close()
        return "", err
    }

    select {
    case <- time.After(READ_TIMEOUT):
        sub.Unsubscribe()
        return "", ErrMessageEmpty
    case msg = <-readMessage:
        sub.Unsubscribe()
        return msg, nil
    }

    return "", ErrMessageEmpty
}
