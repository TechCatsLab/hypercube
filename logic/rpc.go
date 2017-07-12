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

	"hypercube/access/config"
	server "hypercube/access/rpc"
	"hypercube/libs/log"
	"hypercube/libs/message"
	"hypercube/libs/rpc"
)

var (
	configuration = config.Load()
)

// Send calls the function of the access layer remotely, send a message to a specific user.
func Send(user message.User, msg message.Message) error {
	var (
		args server.Args
		ok   bool
	)
	op := rpc.Options{
		Proto: "tcp",
		Addr:  configuration.Addrs,
	}

	client := rpc.Dial(op)
	defer client.Close()

	err := client.Call("AccessRPC.Send", &args, &ok)
	if err != nil {
		log.Logger.Error("RPC Call returned error: %v", err)
		return err
	}

	if !ok {
		return errors.New("logic send message failed")
	}

	return nil
}
