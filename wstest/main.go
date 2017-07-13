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
 *     Initial: 2017/04/12        He ChengJun
 *      Modify: 2017/07/09        Sun Anxiang
 */

package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"time"

	ui "github.com/gizak/termui"
	"github.com/gorilla/websocket"

	"hypercube/libs/message"
)

// websocket————————————————————————————————————————————————————————————————————————————
const (
	userCount = 10
	debugMsg  = false
	Duration  = 600
	Version   = 1
)

type Message struct {
	From    message.User
	To      message.User
	Content string
	Time    string
}

var MsgSendChan = make(chan Message, 10)
var MsgRcvChan = make(chan Message, 10)
var msgSend, msgRsv Message

var addrs string = "127.0.0.1:7000"
var userIDs []string = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go loop()

	for i := 0; i < 1; i++ {
		newRoutine(userIDs[i])
	}

	cmdClear := exec.Command("clear")
	cmdClear.Run()

	draw()

	select {
	case <-interrupt:
		return
	}
}

func newRoutine(from string) {
	go testRoutine(addrs, from)
}

func randUserID() string {
	return userIDs[rand.Uint32()%userCount]
}

func dial(addr string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/join"}
	log.Printf("connecting to %s", u.String())

	connectHeader := make(http.Header)
	connectHeader.Set("Authorization", "Bearer "+"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiIxIn0.csvtKWbINxS2NOWbqq9uS1Q_GUVWa8I7oTZj2rxIdFU")

	c, _, err := websocket.DefaultDialer.Dial(u.String(), connectHeader)

	return c, err
}

type UserAccess struct {
	UserID string
}

func loginPackage(from string) []byte {
	messages := UserAccess{
		UserID: from,
	}
	byteMessage, _ := json.Marshal(messages)

	msg := message.Message{
		Version: Version,
		Type:    message.MessageTypeLogin,
		Content: byteMessage,
	}
	byteMsg, _ := json.Marshal(msg)

	if debugMsg {
		log.Println("login: ", string(byteMsg))
	}

	return byteMsg
}

func testPackage(from, to string, t time.Time) []byte {
	messages := Message{
		From:    message.User{UserID: from},
		To:      message.User{UserID: to},
		Time:    t.String(),
		Content: "test",
	}
	byteMessage, _ := json.Marshal(messages)

	msg := message.Message{
		Version: Version,
		Type:    message.MessageTypePlainText,
		Content: byteMessage,
	}
	byteMsg, _ := json.Marshal(msg)

	if debugMsg {
		log.Println("utu: ", string(byteMsg))
	}

	return byteMsg
}

func writeRoutine(c *websocket.Conn, addr string, from string) {
	var msgCount int32 = 0

	// 写入计时
	ticker := time.NewTicker(time.Second)
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

			pkg, err := UnmarshalPkg(messages)
			if err != nil {
				return
			}
			MsgSendChan <- pkg

			msgCount++
		case <-exitTimer.C:
			log.Println("exitTimer : go routine exit, from = ", from)
			goto exit
		}
	}
exit:
	log.Printf("send %d messages, addr %s, from %d \n", msgCount, addr, from)
}

func testRoutine(addr string, from string) {
	log.Println("new routine, addr = ", addr, "userID = ", from)

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
	exitTimer := time.NewTimer(time.Second * time.Duration(Duration+10))
	defer exitTimer.Stop()

	var msgCount int32 = 0
	for {
		select {
		case <-exitTimer.C:
			goto exit
		default:
		}

		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			goto exit
		}

		pkg, err := UnmarshalPkg(msg)
		if err != nil {
			return
		}
		MsgRcvChan <- pkg

		msgCount++

		if debugMsg {
			log.Printf("to: %d, count: %d, recv: %s \n", from, msgCount, messages)
		}
	}
exit:
	log.Printf("addr = %s recv %d messages, from = %d", addr, msgCount, from)
}

func UnmarshalPkg(pkg []byte) (Message, error) {
	var pack message.Message
	var kage Message

	err := json.Unmarshal(pkg, &pack)
	if err == nil {
		err = json.Unmarshal(pack.Content, &kage)

		if err == nil {
			return kage, nil
		}
	}

	return Message{}, err
}

// ui——————————————————————————————————————————————————————————————————————————————————————————————
const colSpacing = 10

var sendCount int64 = 0
var rsvCount int64 = 0

// per-column width. 0 == auto width
var colWidths = []int{
	10, // send or receive
	5,  // from
	5,  // to
	0,  // time
	0,  // content
	10, // Count
}

type CompactHeader struct {
	Header  *ui.Par
	From    *ui.Par
	To      *ui.Par
	Time    *ui.Par
	Content *ui.Par
	Count   *ui.Par
	pars    []*ui.Par
	X, Y    int
	Width   int
	Height  int
}

func NewCompactHeader() *CompactHeader {
	row := &CompactHeader{
		Header:  NewHeaderPar(""),
		From:    NewSimplePar("From"),
		To:      NewSimplePar("To"),
		Time:    NewSimplePar("Time"),
		Content: NewSimplePar("Content"),
		Count:   NewSimplePar("Count"),
		X:       1,
		Height:  3,
	}

	return row
}

func NewCompactSimple(s string) *CompactHeader {
	row := &CompactHeader{
		Header:  NewHeaderPar(s),
		From:    NewSimplePar("-"),
		To:      NewSimplePar("-"),
		Time:    NewSimplePar("-"),
		Content: NewSimplePar("-"),
		Count:   NewSimplePar("—"),
		X:       1,
		Height:  3,
	}

	return row
}

func NewHeaderPar(s string) *ui.Par {
	p := ui.NewPar(s)
	p.Height = 2
	p.Border = false

	return p
}

func NewSimplePar(s string) *ui.Par {
	p := ui.NewPar(s)
	p.Height = 2
	p.Border = false

	return p
}

func (ch *CompactHeader) SetPars() {
	ch.pars = append(ch.pars, ch.Header, ch.From, ch.To, ch.Time, ch.Content, ch.Count)
}

func (ch *CompactHeader) SetWidth(w int) {
	x := ch.X
	autoWidth := calcWidth(w)

	for n, col := range ch.pars {
		// set column to static width
		if colWidths[n] != 0 {
			col.SetX(x)
			col.SetWidth(colWidths[n])
			x += colWidths[n]
			x += colSpacing
			continue
		}
		col.SetX(x)
		col.SetWidth(autoWidth)
		x += autoWidth + colSpacing
	}
	ch.Width = w
}

func (ch *CompactHeader) SetY(y int) {
	for _, p := range ch.pars {
		p.SetY(y)
	}
	ch.Y = y
}

func (ch *CompactHeader) SetData(message Message) {
	ch.From.Text = message.From.UserID
	ch.To.Text = message.To.UserID
	ch.Time.Text = message.Time
	ch.Content.Text = message.Content
}

// Calculate per-column width, given total width
func calcWidth(width int) int {
	spacing := colSpacing * len(colWidths)

	var staticCols int
	for _, w := range colWidths {
		width -= w
		if w == 0 {
			staticCols += 1
		}
	}

	return (width - spacing) / staticCols
}

// Get package infomation
func loop() {
	for {
		select {
		case spkg := <-MsgSendChan:
			msgSend = spkg
			sendCount++
		case rpkg := <-MsgRcvChan:
			msgRsv = rpkg
			rsvCount++
		}
	}
}

func draw() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	header := NewCompactHeader()
	header.SetPars()
	header.SetWidth(200)
	for _, x := range header.pars {
		x.TextFgColor = ui.ColorBlack
		ui.Render(x)
	}

	sender := NewCompactSimple("send:")
	sender.SetPars()
	sender.SetWidth(200)
	sender.SetY(1)
	sender.SetData(msgSend)
	for _, y := range sender.pars {
		y.TextFgColor = ui.ColorBlue
		ui.Render(y)
	}

	receiver := NewCompactSimple("receive:")
	receiver.SetPars()
	receiver.SetWidth(200)
	receiver.SetY(2)
	receiver.SetData(msgRsv)
	for _, z := range receiver.pars {
		z.TextFgColor = ui.ColorGreen
		ui.Render(z)
	}

	reDraw := func() {
		sc := strconv.FormatInt(sendCount, 10)
		sender.Count.Text = sc
		rc := strconv.FormatInt(rsvCount, 10)
		receiver.Count.Text = rc

		sender.SetData(msgSend)
		receiver.SetData(msgRsv)

		for _, y := range sender.pars {
			ui.Render(y)
		}

		for _, z := range receiver.pars {
			ui.Render(z)
		}
	}

	ui.Merge("/timer/10ms", ui.NewTimerCh(time.Millisecond*10))
	ui.Handle("/timer/10ms", func(e ui.Event) {
		reDraw()
	})

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Loop()
}
