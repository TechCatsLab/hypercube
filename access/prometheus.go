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
 *     Initial: 2017/03/31        Feng Yifei
 *     Modify: 2017/06/04         Yang Chenglong  修改counter创建
 */

package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/client_golang/prometheus"
	"hypercube/libs/log"
	pro "hypercube/libs/metrics/prometheus"
)

func initPrometheus() {
	err := prometheus.Register(pro.OnlineUserCounter)
	if err != nil {
		log.Logger.Error("onlineUser counter couldn't be registered AGAIN, no counting will happen:", err)
		return
	}

	err = prometheus.Register(pro.SendMessageCounter)
	if err != nil {
		log.Logger.Error("sendMessage counter couldn't be registered AGAIN, no counting will happen:", err)
		return
	}

	err = prometheus.Register(pro.ReceiveMessageCounter)
	if err != nil {
		log.Logger.Error("resiveMessage counter couldn't be registered AGAIN, no counting will happen:", err)
		return
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())

		log.Logger.Fatal(http.ListenAndServe(configuration.PrometheusPort, nil))
	}()
}
