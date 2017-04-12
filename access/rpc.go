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
 *     Initial: 2017/04/04        Feng Yifei
 */

package main

import (
	"os"
	"hypercube/common/mq"
)

var (
	natsMQ        *mq.NatsJsonMQ
	logicRequester     *mq.NatsRequester
)

// 创建到 Logic 模块儿 RPC 调用
func initLogicRPC() {
	var (
		err        error
	)

	natsMQ, err = mq.NewNatsMQ(&configuration.NatsUrl)
	if err != nil {
		logger.Error("Initialize RPC with error:", err)
		goto exit
	}

	logicRequester, err = natsMQ.CreateRequester(&configuration.ApiChannel)
	if err != nil {
		logger.Error("Initialize RPC channel with error:", err)
		goto exit
	}

	logger.Debug("RPC messaging channel connected...")
	return

exit:
	if natsMQ != nil {
		natsMQ.Close()
	}
	os.Exit(1)
}
