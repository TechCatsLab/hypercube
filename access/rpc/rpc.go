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
 *     Initial: 2017/07/05        Feng Yifei
 *     Modify : 2017/07/28        Yang Chenglong
 */

package rpc

import (
	"github.com/TechCatsLab/hypercube/access/config"
	"github.com/TechCatsLab/hypercube/access/sender"
	"github.com/TechCatsLab/hypercube/libs/message"
	"github.com/TechCatsLab/hypercube/libs/rpc"
)

var RpcClients *rpc.Clients
var RPCServer *AccessRPC

const RPCNumber = 10

func InitRPC() {
	RPCServer = new(AccessRPC)
	op := rpc.Options{
		Proto: "tcp",
		Addr:  config.GNodeConfig.LogicAddrs,
	}
	options := make([]rpc.Options, RPCNumber)

	for i := 0; i < RPCNumber; i++ {
		options = append(options, op)
	}

	RpcClients = rpc.Dials(options)
}

// AccessRPC provides push functions.
type AccessRPC struct {
	Node sender.Sender
}

// Ping is general rpc keepalive interface.
func (access *AccessRPC) Ping(req *rpc.ReqKeepAlive, resp *rpc.RespKeepAlive) error {
	return nil
}

// Send send a message to a specific user.
func (access *AccessRPC) Send(args *message.Args, reply *bool) error {
	RPCServer.Node.Send(&args.User, &args.Message)
	*reply = true

	return nil
}
