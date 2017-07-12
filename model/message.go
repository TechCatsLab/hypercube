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
 *     Initial: 2017/07/05        Feng Yifei
 */

package message

import (
	"time"

	// TODO: glide install this package
	"github.com/jinzhu/gorm"

	"hypercube/orm"
)

type Message struct {
	Messageid		int64		`jsql:"auto_increment;primary_key;"`
	Source    	    string		`gorm:"not null"`
	Target		    string
	Type            uint16
	IsSend          bool		`gorm:"column:issend"`
	Content         string
	Created         time.Time
}

func (msg *Message) GetOffLineMessage(conn orm.Connection, userid string)([]Message, error) {
	var all []Message

	db := conn.(*gorm.DB).Exec("SET DATABASE = core")

	err := db.Where("target=? AND issend=?", userid, false).Find(&all).Error

	if err != nil {
		return all, err
	}

	return all, nil
}

func (msg *Message) ModifyMessageStatus(conn orm.Connection, msgid int64) error {
	var message Message

	db := conn.(*gorm.DB).Exec("SET DATABASE = core")

	updater := map[string]interface{}{"issend": true}
	err := db.Model(&message).Where("messageid=? AND issend=?", msgid, true).Update(updater).Limit(1).Error

	if err != nil {
		return err
	}

	return nil
}
