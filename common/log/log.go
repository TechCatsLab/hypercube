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

/*
 * Example:
 *    logger := log.S8ECreateLogger(&log.S8ELogTag{"key": "value"}, log.S8ELogLevelInfo)
 *    logger.Debug("....")
 */
import (
	"github.com/Sirupsen/logrus"
)

type S8ELogLevel logrus.Level

const (
	S8ELogLevelDebug = S8ELogLevel(logrus.DebugLevel)
	S8ELogLevelInfo = S8ELogLevel(logrus.InfoLevel)
	S8ELogLevelWarn = S8ELogLevel(logrus.WarnLevel)
	S8ELogLevelError = S8ELogLevel(logrus.ErrorLevel)
)

// 日志通用接口
type S8ELoggerInf interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type S8ELogger struct {
	entry *logrus.Entry
}

type S8ELogTag logrus.Fields

// 创建日志管理工具
func S8ECreateLogger(tags *S8ELogTag, level S8ELogLevel) *S8ELogger {
	entry := logrus.WithFields(logrus.Fields(*tags))
	entry.Logger.Level = logrus.Level(level)

	return &S8ELogger{entry: entry}
}

func (this *S8ELogger) Debug(args ...interface{}) {
	this.entry.Debug(args)
}

func (this *S8ELogger) Info(args ...interface{}) {
	this.entry.Info(args)
}

func (this *S8ELogger) Warn(args ...interface{}) {
	this.entry.Warn(args)
}

func (this *S8ELogger) Error(args ...interface{}) {
	this.entry.Error(args)
}
