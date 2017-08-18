/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
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
 *     Initial: 2017/08/18        Yang Chenglong
 */

package mongo

import (
	"encoding/json"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/hypercube/libs/log"
	"github.com/fengyfei/hypercube/libs/message"
	"github.com/fengyfei/nuts/mgo/refresh"
)

type MgoConnector struct {
	session *mgo.Session
}

type Message struct {
	Messageid bson.ObjectId `bson:"_id"   json:"messageid"`
	From      string
	To        string
	Type      uint16
	Status    int
	Content   string
	Created   time.Time
}

const (
	dbName         = "chat"
	collectionName = "message"
	mgoUrl         = "10.0.0.251:27067"
)

var (
	RefChat *mgo.Collection
)

func (mongo *MgoConnector) Initialize() error {
	var err error

	mongo.session, err = mgo.Dial(mgoUrl)

	if err != nil {
		panic(err)
	}
	mongo.session.SetMode(mgo.Monotonic, true)

	log.Logger.Debug("MongoDB has Connected %v", mgoUrl)

	RefChat = mongo.session.DB(dbName).C(collectionName)
	nameIndex := mgo.Index{
		Key:        []string{"From"},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	if err := RefChat.EnsureIndex(nameIndex); err != nil {
		panic(err)
	}

	return nil
}

func (mongo *MgoConnector) Put(msg *message.Message, status int) error {
	var content message.PlainText
	var dbmsg Message

	err := json.Unmarshal(msg.Content, &content)
	if err != nil {
		return err
	}

	dbmsg = Message{
		Messageid: bson.NewObjectId(),
		From:      content.From.UserID,
		To:        content.To.UserID,
		Type:      msg.Type,
		Status:    status,
		Content:   content.Content,
		Created:   time.Now(),
	}

	return refresh.Insert(mongo.session, RefChat, &dbmsg)
}

func (mongo *MgoConnector) Get(id string, status int) ([]Message, error) {
	var all []Message

	filter := bson.M{"To": id, "Status": status}
	err := refresh.GetMany(mongo.session, RefChat, filter, &all)

	return all, err
}

func (mongo *MgoConnector) Update(id string, status int) error {
	filter := bson.M{"_id": bson.ObjectIdHex(id), "Status": status}
	updater := bson.M{"$set": bson.M{"Status": message.MessageSent}}

	return refresh.Update(mongo.session, RefChat, filter, updater)
}
