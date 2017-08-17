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
 *     Modify : 2017/07/28        Yang Chenglong
 */

package rpc

import (
	"net"
	"net/rpc"

	"github.com/fengyfei/hypercube/access/config"
	"github.com/fengyfei/hypercube/libs/log"
)

// InitServer initialize the RPC server.
func InitServer() {
	rpc.Register(new(AccessRPC))
	rpc.HandleHTTP()

	go rpcListen()
}

func rpcListen() {
	l, err := net.Listen("tcp", config.GNodeConfig.Address)
	if err != nil {
		log.Logger.Error("net.Listen(\"%s\", \"%s\") error(%v)"+"tcp"+config.GNodeConfig.Address, err)
		panic(err)
	}

	defer func() {
		log.Logger.Info("listen rpc: \"%s\" close", config.GNodeConfig.Address)
		if err := l.Close(); err != nil {
			log.Logger.Error("listener.Close() error(%v)", err)
		}
	}()
	rpc.Accept(l)
}
