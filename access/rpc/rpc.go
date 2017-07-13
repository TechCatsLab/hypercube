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

package rpc

import (
	"hypercube/libs/message"
	"hypercube/libs/rpc"
	"hypercube/access/config"
	"hypercube/access/sender"
)

var RpcClient *rpc.Client

func InitRPC() {
	op := rpc.Options{
		Proto:  "tcp",
		Addr:   config.GNodeConfig.LogicAddrs,
	}
	RpcClient = rpc.Dial(op)
}

// AccessRPC provides push functions.
type AccessRPC struct {
	node sender.Sender
}

// Ping is general rpc keepalive interface.
func (access *AccessRPC) Ping(req *rpc.ReqKeepAlive, resp *rpc.RespKeepAlive) error {
	return nil
}

// Push a message to a specific user.
func Push(user *message.User) error {
	return nil
}

// Send send a message to a specific user.
func (access *AccessRPC) Send(args *message.Args, reply *bool) error {
	access.node.Send(&args.User, &args.Message)
	*reply = true

	return nil
}
