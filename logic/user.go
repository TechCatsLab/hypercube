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
 *     Initial: 2017/04/04        Feng Yifei
 *     Modify:  2017/06/07        Yang Chenglong  添加AccessHeartHandler
 */

package main

import (
	"encoding/json"

	"github.com/jinzhu/gorm"

	"hypercube/libs/log"
	"hypercube/libs/message"
	"hypercube/orm/cockroach"
	db "hypercube/model"
)

func (this *UserManager) LoginHandle(user message.UserEntry, reply *int) error {
	err := this.Add(user)
	if err != nil {
		log.Logger.Error("LoginHandle Add Error %+v: .", err)
		return err
	}

	conn, err := cockroach.DbConnPool.GetConnection()
	if err != nil {
		log.Logger.Error("Get cockroach connect error:", err)
		return err
	}
	defer cockroach.DbConnPool.ReleaseConnection(conn)

	mes, err := db.MessageService.GetOffLineMessage(conn, user.UserID.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Debug("")
			goto Mess
		}

		log.Logger.Error("GetOffLineMessage Error %+v", err)
		return err
	}
Mess:
	for _, msg := range mes {
		switch msg.Type {
		case message.MessageTypePlainText:
			content := message.PlainText{
				From:    message.User{UserID:msg.Source},
				To:      message.User{UserID:msg.Target},
				Content: msg.Content,
			}

			text, err := json.Marshal(content)
			if err != nil {
				log.Logger.Error("OffLineMessage Marshal Error %+v", err)
				return err
			}

			mesg := &message.Message{
				Type:       msg.Type,
				Version:    msg.Version,
				Content:    text,
			}

			TransmitMsg(mesg)
		}
	}

	return nil
}

