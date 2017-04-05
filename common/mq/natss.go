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
    "bytes"
    "encoding/binary"
    "github.com/nats-io/go-nats-streaming"
    "github.com/nats-io/go-nats-streaming/pb"
    "hypercube/proto/general"
)

var (
    ErrMessageEmpty = errors.New("message queue is empty ")
)

var (
    READ_TIMEOUT = time.Second * 5
)

const (
    clusterID = "test-cluster"
    subject = "chatmessage"
    durable = "chatmessage"
)

type NatsStreaming struct {
    subject     *string
    durable		*string
    conn 		*stan.Conn
}

func ConnectToServer(urls, clientID  *string) (*NatsStreaming, error) {
    sc, err := stan.Connect(clusterID, *clientID, stan.NatsURL(*urls))

    if err != nil {
        return nil, err
    }
    return &NatsStreaming{conn: &sc}, nil
}

func (ns *NatsStreaming)WriteMessage(msg interface{}) error  {

    m, _ := msg.(*general.Proto)
    mBytes,_ := protpToBytes(m)
    err := (*ns.conn).Publish(subject, mBytes)

    if err != nil {
        return err
    }
    return nil
}

func (ns *NatsStreaming)ReadMessage() (msg interface{},  err error) {

    startOpt := stan.StartAt(pb.StartPosition_NewOnly)
    readMessage := make(chan *general.Proto, 1)
    var (
        sub stan.Subscription
    )

    messageHandle := func(msg *stan.Msg){
        m, _ := bytesToProto(msg.Data)
        readMessage <- m
    }

    sub, err = (*ns.conn).Subscribe(subject, messageHandle, startOpt, stan.DurableName(durable))
    if err != nil {
        (*ns.conn).Close()
        return nil, err
    }

    select {
    case <- time.After(READ_TIMEOUT):
        sub.Unsubscribe()
        return nil, ErrMessageEmpty
    case msg = <-readMessage:
        sub.Unsubscribe()
        return msg, nil
    }

    return nil, ErrMessageEmpty
}

func bytesToProto(buffer []byte) (*general.Proto, error)  {
    var proto general.Proto

    buf := bytes.NewReader(buffer)
    err := binary.Read(buf, binary.BigEndian, &proto)

    if err != nil {
        return nil, err
    }
    return &proto, nil
}

func protpToBytes(proto *general.Proto) ([]byte, error)  {
    buf := new(bytes.Buffer)

    err := binary.Write(buf, binary.BigEndian, *proto)

    if err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}
