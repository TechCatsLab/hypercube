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
 *     Initial: 2017/03/29        Feng Yifei
 */

package main

import (
	"hypercube/common/interrupt"
	"os"
	"syscall"
)

var (
	sigHandler *interrupt.Handler
)

func initSignal() {
	sigHandler = interrupt.New(finalHandler, func() {})
	logger.Debug("Interrupt handler initialized")
}

func finalHandler(sig os.Signal) {
	switch sig {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
		logger.Info("Signal quit/term/int captured")
		return

	case syscall.SIGHUP:
		logger.Info("Signal hup captured")
		return

	case syscall.SIGALRM:
		logger.Info("Signal alrm captured")
		
		OnLineManagement.PrintDebugInfo()
		return
	}
}
