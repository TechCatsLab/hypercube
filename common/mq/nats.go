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
 *     Initial: 2017/04/01        Feng Yifei
 */

package mq

import (
	"strings"
	"github.com/nats-io/go-nats"
	"errors"
)

var (
	ErrMessageQueueFull = errors.New("Message queue full ")
)

type NatsJsonMQ struct {
	conn *nats.EncodedConn
}

type NatsPublisher struct {
	subject   *string
	sender    chan interface{}
}

func newNatsPublisher(subject *string) *NatsPublisher {
	sender := make(chan interface{}, 1)

	return &NatsPublisher{
		subject: subject,
		sender: sender,
	}
}

func (this *NatsPublisher) SendMessage(v interface{}) error {
	this.sender <- v
	return nil
}

func (this *NatsPublisher) PostMessage(v interface{}) error {
	select {
	case this.sender <-v:
		return nil

	default:
		return ErrMessageQueueFull
	}
}

func NewNatsMQ(urls *string) (*NatsJsonMQ, error) {
	var (
		err      error
		opts     nats.Options
		conn     *nats.Conn
		encoded  *nats.EncodedConn
		q        *NatsJsonMQ
	)

	opts = nats.DefaultOptions
	opts.Servers = strings.Split(*urls, ",")

	conn, err = opts.Connect()
	if err != nil {
		return nil, err
	}

	encoded, err = nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		conn.Close()
		return nil, err
	}

	q = &NatsJsonMQ{
		conn: encoded,
	}

	return q, nil
}

func (this *NatsJsonMQ) CreateProducer(v interface{}) (interface {}, error) {
	subject, _ := v.(*string)
	producer := newNatsPublisher(subject)

	err := this.conn.BindSendChan(*subject, producer.sender)

	if err != nil {
		return nil, err
	}

	return producer, nil
}
