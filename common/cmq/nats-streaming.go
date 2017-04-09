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
 *     Initial: 2017/04/09        HeCJ
 */

package cmq

import (
    "github.com/nats-io/go-nats-streaming"
    "github.com/nats-io/go-nats-streaming/pb"
    "time"
)

type NatssCMQ struct {
    conn 		stan.Conn
}

type NatssPublisher struct {
    cmq       *NatssCMQ
    subject   *string
}

type NatssSubcriber struct {
    subject *string
    durable *string
    qgroup  *string
    cmq     *NatssCMQ
    sub     stan.Subscription
    handle  SubHandle
}

type NatssHistory struct {
    subject *string
    cmq     *NatssCMQ
    sub     stan.Subscription
    handle  SubHandle
}

func NewNatssCMQ(urls, clusterID, clientID *string) (*NatssCMQ, error) {
    sc, err := stan.Connect(*clusterID, *clientID, stan.NatsURL(*urls))

    if err != nil {
        return nil, err
    }
    return &NatssCMQ{
        conn:   sc,
    }, nil
}

func (ncmq *NatssCMQ)NewPublisher(sub interface{}) *NatssPublisher {
    subject, _ := sub.(*string)

    return &NatssPublisher{
        cmq: ncmq,
        subject: subject,
    }
}


func (ncmq *NatssCMQ)NewSubscriber(sub, qg, dur interface{}) *NatssSubcriber {
    subject, _ := sub.(*string)
    qgroup, _ := qg.(*string)
    durable, _ := dur.(*string)

    return &NatssSubcriber{
        subject: subject,
        durable: durable,
        qgroup: qgroup,
        cmq: ncmq,
    }
}


func (ncmq *NatssCMQ)NewHistory(sub interface{}) *NatssHistory {
    subject, _ := sub.(*string)

    return &NatssHistory{
        subject: subject,
        cmq: ncmq,
    }
}

func (nsp *NatssPublisher)Publish(msg interface{}) error {
    m, _ := msg.([]byte)
    err := nsp.cmq.conn.Publish(*nsp.subject, m)

    if err != nil {
        return err
    }
    return nil
}

func (nss *NatssSubcriber)SetSubHandle(h SubHandle) {

    nss.handle = h
}

func (nss *NatssSubcriber)Start() error {
    var (
        err  error
    )

    startOpt := stan.StartAt(pb.StartPosition_NewOnly)

    messageHandle := func(msg *stan.Msg){
        nss.handle(msg)
    }

    nss.sub, err = nss.cmq.conn.QueueSubscribe(*nss.subject, *nss.qgroup, messageHandle, startOpt, stan.DurableName(*nss.durable))

    if err != nil {
        return err
    }
    return nil
}


func (nss *NatssSubcriber)Stop() error {
    if nss.sub != nil {
        nss.sub.Unsubscribe()
    }
    return nss.cmq.conn.Close()
}

func (nsh *NatssHistory)SetSubHandle(h SubHandle) {
    nsh.handle = h
}

func (nsh *NatssHistory)Start(startDelta interface{}) error {
    var (
        err  error
    )

    delta := startDelta.(time.Duration)
    startOpt := stan.StartAtTimeDelta(delta)

    messageHandle := func(msg *stan.Msg){
        nsh.handle(msg)
    }

    nsh.sub, err = nsh.cmq.conn.Subscribe(*nsh.subject, messageHandle, startOpt)

    if err != nil {
        return err
    }
    return nil
}


func (nsh *NatssHistory)Stop() error {
    if nsh.sub != nil {
        nsh.sub.Unsubscribe()
    }
    return nsh.cmq.conn.Close()
}
