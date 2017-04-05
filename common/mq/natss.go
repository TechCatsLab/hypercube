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
    "encoding/json"
    "github.com/nats-io/go-nats-streaming"
    "github.com/nats-io/go-nats-streaming/pb"
)

var (
    ErrMessageEmpty = errors.New("message queue is empty ")
)

type NatsStreaming struct {
    subject     *string
    durable		*string
    conn 		*stan.Conn
    timeout     uint32
    readMessage chan interface{}
}

type Proto struct {
    Ver     uint32
    Body    []byte
}

func ConnectToServer(urls, clusterID, clientID, subject, durable  *string, timeout uint32) (*NatsStreaming, error) {
    sc, err := stan.Connect(*clusterID, *clientID, stan.NatsURL(*urls))

    if err != nil {
        return nil, err
    }
    return &NatsStreaming{
        subject: subject,
        durable: durable,
        conn: &sc,
        readMessage: make(chan interface{}, 1),
        timeout: timeout,
    }, nil
}

func (ns *NatsStreaming)WriteMessage(msg interface{}) error  {

    m, err := json.Marshal(msg)
    if err != nil {
        return err
    }

    err = (*ns.conn).Publish(*ns.subject, m)

    if err != nil {
        return err
    }
    return nil
}

func (ns *NatsStreaming)ReadMessage(startPosition uint64, count uint32) (msg interface{},  err error) {

    startOpt := stan.StartAt(pb.StartPosition_NewOnly)
    if startPosition > 0 {
        startOpt = stan.StartAtSequence(startPosition)
    }

    var (
        sub stan.Subscription
        msgs []*stan.Msg
    )

    messageHandle := func(msg *stan.Msg){
        msgs = append(msgs, msg)
        if uint32(len(msgs)) >= count {
            ns.readMessage <- msgs
        }
    }

    sub, err = (*ns.conn).Subscribe(*ns.subject, messageHandle, startOpt, stan.DurableName(*ns.durable))
    if err != nil {
        (*ns.conn).Close()
        return nil, err
    }

    select {
    case <- time.After(time.Millisecond * time.Duration(ns.timeout)):
        sub.Unsubscribe()
        if len(msgs) > 0 {
            return msgs, nil
        }
        return nil, ErrMessageEmpty
    case msg = <-ns.readMessage:
        sub.Unsubscribe()
        return msg, nil
    }

    return nil, ErrMessageEmpty
}
