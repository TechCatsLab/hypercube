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
 *     Initial: 2017/07/12        Yang Chenglong
 */

package main

import (
	"net/rpc"
	"errors"
	"net"

	"hypercube/libs/log"
	rp "hypercube/libs/rpc"
	"hypercube/libs/message"
)

func initServer()  {
	rpc.Register(new(LogicRPC))
	rpc.HandleHTTP()

	go rpcListen()
}

func rpcListen()  {
	l, err := net.Listen("tcp", configuration.Addrs)
	if err != nil {
		log.Logger.Error("net.Listen(\"%s\", \"%s\") error(%v)" + "tcp" + configuration.Addrs, err)
		panic(err)
	}

	defer func() {
		log.Logger.Info("listen rpc: \"%s\" close", configuration.Addrs)
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

	client, err := clients.Get(op.Addr)
	if err != nil {
		log.Logger.Error("Clients.Get returned error: %v.", err)

		return err
	}

	err = client.Call("AccessRPC.Send", &args, &ok)
	if err != nil {
		log.Logger.Error("RPC Call returned error: %v", err)
		client.Close()

		return err
	}

	if !ok {
		client.Close()

		return errors.New("logic send message failed")
	}

	client.Close()

	return nil
}
