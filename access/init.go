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
 *     Initial: 2017/04/05            HeCJ
 *     Modify : 2017/06/24            Yang Chenglong
 */

package main

import (
	"github.com/TechCatsLab/hypercube/access/config"
	"github.com/TechCatsLab/hypercube/access/endpoint"
	"github.com/TechCatsLab/hypercube/access/rpc"
)

var (
	configuration = config.Load()
	ep            *endpoint.Endpoint
)

func init() {
	initSignal()
	rpc.InitRPC()
	HttpPprof()
	initPrometheus()
}

func run() {
	// Start a access endpoint.
	ep = endpoint.NewEndpoint(configuration)
	go rpc.RpcClients.Ping("LogicRPC.Ping")

	go rpc.InitServer()

	ep.Run()

	sigHandler.Wait()
}
