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
 *     Initial: 2017/07/28        Yang Chenglong
 */

package main

import (
	"net"
	"net/rpc"

	"github.com/TechCatsLab/hypercube/libs/log"
	"github.com/TechCatsLab/hypercube/libs/message"
	rp "github.com/TechCatsLab/hypercube/libs/rpc"
)

func initServer() {
	rpc.Register(new(LogicRPC))
	rpc.HandleHTTP()

	go rpcListen()
}

func rpcListen() {
	l, err := net.Listen("tcp", configuration.Addr)
	if err != nil {
		log.Logger.Error("net.Listen(\"%s\", \"%s\") error(%v)"+"tcp"+configuration.Addr, err)
		panic(err)
	}

	defer func() {
		log.Logger.Info("listen rpc: \"%s\" close", configuration.Addr)
		if err := l.Close(); err != nil {
			log.Logger.Error("listener.Close() error(%v)", err)
		}
	}()
	rpc.Accept(l)
}

func Send(user message.User, msg message.Message, op rp.Options) error {
	var (
		args message.Args
		ok   bool
	)

	args = message.Args{
		User:    user,
		Message: msg,
	}

	rpcClient, err := rpcClients.Get(op.Addr)
	if err != nil || rpcClient == nil {
		log.Logger.Error("Clients.Get returned error: %v.", err)

		return err
	}

	err = rpcClient.Call("AccessRPC.Send", &args, &ok)

	if err != nil || !ok {
		log.Logger.Error("Logic Call Send failed: %v", err)
		rpcClient.Close()

		return err
	}

	return nil
}
