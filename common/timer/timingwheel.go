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
 *     Initial: 2017/04/01        Yusank Kurban
 */

package timer

import (
    "sync"
    "time"
)

type TimingWheel struct {
    sync.Mutex
    interval    time.Duration
    ticker      *time.Ticker
    quit        chan struct{}
    maxTimeout  time.Duration
    cs          []chan struct{}
    pos         int
}

func NewTimingWheel(interval time.Duration, buckets int) *TimingWheel {
    w := new(TimingWheel)

    w.interval = interval
    w.quit = make(chan struct{})
    w.pos = 0
    w.maxTimeout = time.Duration(interval * (time.Duration(buckets)))
    w.cs = make([]chan struct{}, buckets)

    for i := range w.cs {
        w.cs[i] = make(chan struct{})
    }

    w.ticker = time.NewTicker(interval)
    go w.run()

    return w
}

func (w *TimingWheel) Stop() {
    close(w.quit)
}

func (w *TimingWheel) After(timeout time.Duration) <-chan struct{} {
    if timeout >= w.maxTimeout {
        panic("timeout too much, over maxtimeout")
    }

    index := int(timeout / w.interval)
    if index > 0 {
        index--
    }

    w.Lock()

    index = (w.pos + index) % len(w.cs)
    b := w.cs[index]

    w.Unlock()

    return b
}

func (w *TimingWheel) run() {
    for {
        select {
        case <-w.ticker.C:
            w.onTicker()
        case <-w.quit:
            w.ticker.Stop()
            return
        }
    }
}

func (w *TimingWheel) onTicker() {
    w.Lock()

    lastC := w.cs[w.pos]
    w.cs[w.pos] = make(chan struct{})
    w.pos = (w.pos + 1) % len(w.cs)

    w.Unlock()

    close(lastC)
}
