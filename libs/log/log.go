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
 *     Initial: 2017/03/28        Feng Yifei
 */

package log

import (
	"go.uber.org/zap"
	"fmt"
	"log"
)

type RecordLog struct {}

var (
	Logger *RecordLog
	zapLog *zap.Logger
)

func init() {
	Logger = &RecordLog{}
	zapLog, _ = zap.NewDevelopment()
}

func (l *RecordLog) Error(desc string, err error) {
	zapLog.Error(desc, zap.Error(err))
}

func (l *RecordLog) Debug(format string, a ...interface{}) {
	info := fmt.Sprintf(format, a)
	zapLog.Debug(info, zap.Skip())
}

func (l *RecordLog) Fatal(v ...interface{}){
	log.Fatal(v)
}

func (l *RecordLog) Info(format string, a ...interface{}) {
	info := fmt.Sprintf(format, a)
	zapLog.Info(info, zap.Skip())
}

func (l *RecordLog) Warn(format string, a ...interface{}) {
	warn := fmt.Sprintf(format, a)
	zapLog.Warn(warn, zap.Skip())
}
