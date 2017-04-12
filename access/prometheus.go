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
 */

package main

import (
	"net/http"
	"log"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	onlineUserDurations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:       "onlineUser_durations_seconds",
			Help:       "OnlineUser latency distributions.",
		},
		[]string{"service"},
	)

	sendMessageDurations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:       "sendMessage_durations_seconds",
			Help:       "SendMessage latency distributions.",
		},
		[]string{"service"},
	)

	resiveMessageDurations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:       "resiveMessage_durations_seconds",
			Help:       "ResiveMessage latency distributions.",
		},
		[]string{"service"},
	)
)

func initPrometheus() {
	prometheus.MustRegister(onlineUserDurations)
	prometheus.MustRegister(sendMessageDurations)
	prometheus.MustRegister(resiveMessageDurations)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(configuration.PrometheusPort, nil))
	}()
}
