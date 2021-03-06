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
 *     Initial: 2017/07/11        Jia Chenhui
 *     Modify : 2017/07/28        Yang Chenglong
 */

package main

import (
	"github.com/TechCatsLab/hypercube/libs/log"
	"github.com/TechCatsLab/hypercube/libs/message"
	"github.com/TechCatsLab/hypercube/libs/rpc"
)

const RPCNumber = 10

type LogicRPC struct {
}

var (
	options    []rpc.Options
	rpcClients *rpc.Clients
)

func initRPC() {
	for _, addr := range configuration.AccessAddrs {
		op := rpc.Options{
			Proto: "tcp",
			Addr:  addr,
		}

		for i := 0; i < RPCNumber; i++ {
			options = append(options, op)
		}
	}

	rpcClients = rpc.Dials(options)
}

func (lr *LogicRPC) LoginHandler(user message.UserEntry, reply *int) error {
	err := onlineUserManager.Add(user)
	if err != nil {
		log.Logger.Error("LoginHandle Add Error %+v: ", err)
		*reply = message.ReplyFailed
		return err
	}

	offline <- user
	*reply = message.ReplySucceed
	return nil
}

func (lr *LogicRPC) LogoutHandle(user message.UserEntry, reply *int) error {
	err := onlineUserManager.Remove(user)
	if err != nil {
		log.Logger.Error("LogoutHandle Error %v", err)
		*reply = message.ReplyFailed
		return err
	}

	*reply = message.ReplySucceed
	return nil
}

func (lr *LogicRPC) Add(msg *message.Message, reply *bool) error {
	Queue <- msg
	*reply = true

	return nil
}

func (lr *LogicRPC) Ping(req *rpc.ReqKeepAlive, resp *rpc.RespKeepAlive) error {
	return nil
}
