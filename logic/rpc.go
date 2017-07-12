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
 *     Initial: 2017/07/11        Jia Chenhui
 */

package main

import (
	"errors"

	server "hypercube/access/rpc"
	"hypercube/libs/log"
	"hypercube/libs/message"
	"hypercube/libs/rpc"
)

var (
	options []rpc.Options
	clients *rpc.Clients
)

func init() {
	for _, addr := range configuration.AccessAddrs {
		op := rpc.Options{
			Proto: "tcp",
			Addr:  addr,
		}

		options = append(options, op)
	}

	clients = rpc.Dials(options)
}

// Send calls the function of the access layer remotely, send a message to a specific user.
func Send(user message.User, msg message.Message, op rpc.Options) error {
	var (
		args server.Args
		ok   bool
	)

	args = server.Args{
		User:    user,
		Message: msg,
	}

	client, err := clients.Get(op.Addr)
	if err != nil {
		log.Logger.Error("Clients.Get returned error: %+v.", err)

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
