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
 *     Initial: 2017/04/09        He CJ
 */

package main

import (
    "hypercube/common/cmq"
)

var (
    natssCMQ *cmq.NatssCMQ
    publisher *cmq.NatssPublisher
    subscriber *cmq.NatssSubcriber
    history *cmq.NatssHistory
)

func initCMQ() error {
    var (
        err error
    )
    natssCMQ, err = cmq.NewNatssCMQ(
        &configuration.NatssUrl,
        &configuration.NatssClusterID,
        &configuration.NatssClientID )

    if err != nil {
        logger.Error("Logic request processor error:", err)
        return err
    }

    publisher = natssCMQ.NewPublisher(&configuration.NatssSubject)

    subscriber = natssCMQ.NewSubscriber(
        &configuration.NatssSubject,
        &configuration.NatssQGroup,
        &configuration.NatssDurable )

    history = natssCMQ.NewHistory(&configuration.NatssSubject)

    return nil
}
