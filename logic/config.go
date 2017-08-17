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
 *     Initial: 2017/06/08        Feng Yifei
 *     Modify : 2017/07/28        Yang Chenglong
 */

package main

import (
	"github.com/spf13/viper"
	"github.com/fengyfei/hypercube/libs/log"
)

// 配置文件结构
type LogicLayerConfig struct {
	Addrs          string
	AccessAddrs    []string
	PprofAddrs     string
	PrometheusPort string
}

var configuration *LogicLayerConfig

// 初始化配置
func readConfiguration() {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		log.Logger.Error("Read configuration file with error:", err)
		panic(err)
	}

	configuration = &LogicLayerConfig{
		Addrs:          viper.GetString("serverAddrs"),
		AccessAddrs:    viper.GetStringSlice("accessAddrs"),
		PprofAddrs:     viper.GetString("monitor.pprofAddrs"),
		PrometheusPort: viper.GetString("monitor.prometheusPort"),
	}
}
