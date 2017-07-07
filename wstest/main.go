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
 * Created by HeChengJun on 12/04/2017.
 * Modify by SunAnxiang on 2017/07/06.
 */

package main

import (
	"os"
	"log"
	"sync"
	"time"
	"net/url"
	"math/rand"
	"os/signal"
	"encoding/json"

	"github.com/gorilla/websocket"

	"hypercube/libs/message"
)

const (
	userCount		= 10
	debugMsg		= false
	Duration		= 600
	Version			= 1
)

var addrs 	string   = "10.0.0.116:7000"
var userIDs []uint64 = []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for i := 0; i < userCount; i++ {
		newRoutine(userIDs[i])
	}

	select {
	case <-interrupt:
		return
	}
}

func newRoutine(from uint64) {
	go testRoutine(addrs, from)
}

func randUserID() uint64 {
	return userIDs[rand.Uint32()%userCount]
}

func dial(addr string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/join"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	return c, err
}

type UserAccess struct {
	UserID	uint64
}

func loginPackage(from uint64) []byte {
	messages := UserAccess{
		UserID: from,
	}
	byteMessage, _ := json.Marshal(messages)

	msg := message.Message{
		Version: Version,
		Type: message.MessageTypeLogin,
		Content: byteMessage,
	}
	byteMsg, _ := json.Marshal(msg)

	if debugMsg {
		log.Println("login: ", string(byteMsg))
	}

	return byteMsg
}

type Message struct {
	From message.User
	To   message.User
	Content string
	SendOrder int64
}

type Count struct {
	counter int64
	lock sync.Mutex
}

var counter = Count{
	counter: 0,
	lock: sync.Mutex{},
}

func testPackage(from, to uint64, t time.Time) []byte {
	counter.lock.Lock()
	counter.counter ++
	messages := Message{
		From:    message.User{UserID:"fdsjk"},
		To:      message.User{UserID:"fdsjk"},
		Content: t.String(),
		SendOrder: counter.counter,
	}
	counter.lock.Unlock()
	byteMessage, _ := json.Marshal(messages)

	msg := message.Message{
		Version:  Version,
		Type: message.MessageTypePlainText,
		Content: byteMessage,
	}
	byteMsg, _ := json.Marshal(msg)

	if debugMsg {
		log.Println("utu: ", string(byteMsg))
	}

	return byteMsg
}

func writeRoutine(c *websocket.Conn, addr string, from uint64) {
	var msgCount int32 = 0

	// 写入计时
	ticker := time.NewTicker(time.Microsecond * time.Duration(Duration))
	defer ticker.Stop()

	// 退出计时
    exitTimer := time.NewTimer(time.Second * time.Duration(Duration))
	defer exitTimer.Stop()

	for {
		select {
		case t := <-ticker.C:
			// 发送
			to := randUserID()
			messages := testPackage(from, to, t)

			err := c.WriteMessage(websocket.TextMessage, messages)
			if err != nil {
				log.Println("write:", err)
				goto exit
			}
			msgCount ++
		case <-exitTimer.C:
			log.Println("exitTimer : go routine exit, from = ", from)
			goto exit
		}
	}
exit:
	log.Printf("send %d messages, addr %s, from %d \n", msgCount, addr, from)
}

func testRoutine(addr string, from uint64) {
	log.Println("new routine, addr = ",addr ,"userID = ", from)

	// 拨号
	c, err := dial(addr)
	if err != nil {
		log.Println("dial:", err)
		return
	}
	defer c.Close()

	// 发送登录数据包
	messages := loginPackage(from)
	err = c.WriteMessage(websocket.TextMessage, messages)
	if err != nil {
		log.Println("write:", err)
		return
	}

	// 写
	go writeRoutine(c, addr, from)

	// 读
	exitTimer := time.NewTimer(time.Second * time.Duration(Duration + 10))
	defer exitTimer.Stop()

	var msgCount int32 = 0
	for {
		select {
		case <- exitTimer.C :
			goto exit
		default:
		}

		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			goto exit
		}
		log.Println("info:", string(msg))

		msgCount++

        if debugMsg {
            log.Printf("to: %d, count: %d, recv: %s \n", from, msgCount, messages)
        }
	}
exit:
    log.Printf("addr = %s recv %d messages, from = %d", addr, msgCount, from)
}
