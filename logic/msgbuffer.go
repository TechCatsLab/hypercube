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

/**
 * Created by HeChengJun on 13/04/2017.
 */

package main

var(
    msgbuf map[uint64][]interface{}
)

const (
    maxBufferSize = 100
)

func init()  {
    msgbuf = make(map[uint64][]interface{})
}

func addHistMessage(userID uint64, msg interface{})  {
    if msgbuf[userID] == nil {
        msgbuf[userID] = make([]interface{}, maxBufferSize)
    }
    if len(msgbuf[userID]) < maxBufferSize {
        msgbuf[userID] = append(msgbuf[userID], msg)
    }
}

func getHistMessages(userID uint64) []interface{} {
    return msgbuf[userID]
}

func clearHistMessages(userID uint64) {
    msgbuf[userID] = []interface{}{}
}
