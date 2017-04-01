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
 *     Initial: 2017/03/30        Feng Yifei
 *     AddFunc: 2017/04/01        Yusank Kurban
 */

package timer

import (
	"time"
    "sync"
    "fmt"

    "hypercube/common/log"
)

const (
    timeFormat       = "2006-01-02 15:04:05"
    infiniteDuration = time.Duration(1<<63 - 1)
)

var common *log.S8ELogger = log.S8ECreateLogger(
    &log.S8ELogTag{
        log.LogTagService: "common",
        log.LogTagType: "timer",
    },
    log.S8ELogLevelDefault)

type TimerFunc func()

type Block interface {
	Expired() bool
	ExpiredString() string
}

type BlockBasic struct {
    Key         string
    expire      time.Time
    fn          TimerFunc
    index       int
    next        *BlockBasic
}


func (bb *BlockBasic) Expired() bool {
    sub := bb.expire.Sub(time.Now())

    if sub > 0 {
        return false
    }

    return true
}

func (bb *BlockBasic) ExpiredString() string {
    return bb.expire.Format(timeFormat)
}

type Timer interface {
	Create(capacity int)
	Add(expire time.Duration, fn TimerFunc) *Block
	Set(block *Block, expire time.Duration)
	Remove(block *Block)
}

type TimeBasic struct {
    lock        sync.Mutex
    free        *BlockBasic
    timers      []*BlockBasic
    signal      *time.Timer
    capacity    int
}

func NewTimer(capacity int) (tb *TimeBasic) {
    tb = new(TimeBasic)
    tb.init(capacity)
    return tb
}

func (tb *TimeBasic) Create(capacity int) {
    tb.init(capacity)
}

// 初始化 timer
func (tb *TimeBasic) init(capacity int){
    tb.signal = time.NewTimer(infiniteDuration)
    tb.timers = make([]*BlockBasic, 0, capacity)
    tb.capacity = capacity
    tb.grow()
    go tb.start()
}

func (tb *TimeBasic) grow() {
    var (
        i   int
        bb  *BlockBasic
        bbs = make([]BlockBasic, tb.capacity)
    )

    tb.free = &(bbs[0])
    bb = tb.free

    for i = 1; i < tb.capacity; i++ {
        bb.next = &(bbs[i])
        bb = bb.next
    }
    bb.next = nil

    return
}

func (tb *TimeBasic) get() (bb *BlockBasic) {
    bb = tb.free
    if bb == nil {
        tb.grow()
        bb = tb.free
    }

    tb.free = bb.next
    return
}

func (tb *TimeBasic) put(bb *BlockBasic) {
    bb.fn = nil
    bb.next = tb.free
    tb.free = bb
}

func (tb *TimeBasic) Add(expire  time.Duration, fn TimerFunc) (bb *BlockBasic) {
    tb.lock.Lock()
    bb = tb.get()
    bb.expire = time.Now().Add(expire)
    bb.fn = fn
    tb.add(bb)
    tb.lock.Unlock()
    return
}

func (tb *TimeBasic) Remove(bb *BlockBasic) {
    tb.lock.Lock()
    tb.remove(bb)
    tb.put(bb)
    tb.lock.Unlock()
    return
}

// 更新 timer 数据
func (tb *TimeBasic) Set(bb *BlockBasic, expire time.Duration) {
    tb.lock.Lock()
    tb.remove(bb)
    bb.expire = time.Now().Add(expire)
    tb.add(bb)
    tb.lock.Unlock()
    return
}

func (tb *TimeBasic) add(bb *BlockBasic) {
    var d time.Duration

    bb.index = len(tb.timers)
    tb.timers = append(tb.timers, bb)
    tb.up(bb.index)

    if bb.index == 0 {
        d = bb.expire.Sub(time.Now())
        tb.signal.Reset(d)
        if false {
            format := fmt.Sprintf("timer: add reset delay %d ms" , int64(d) / int64(time.Millisecond))
            common.Debug(format)
        }
    }
    if false {
        format := fmt.Sprintf("timer: push item key: %s, expire: %s, index: %d", bb.Key, bb.ExpiredString(), bb.index)
        common.Debug(format)
    }

}

func (tb *TimeBasic) remove(bb *BlockBasic) {
    var (
        i     = bb.index
        last  = len(tb.timers) - 1
    )

    if i < 0 || i > last || tb.timers[i] != bb {
        // 已经被移除，通常是因为到期 (expire)
        if false {
            format := fmt.Sprintf("timer remove i: %d, last: %d, %p", i, last, bb)
            common.Debug(format)
        }

        return
    }

    if i != last {
        tb.swap(i, last)
        tb.down(i, last)
        tb.up(i)
    }
    // 被移除的元素是最后一个节点的元素
    tb.timers[last].index = -1
    tb.timers = tb.timers[:last]
    if false {
        format := fmt.Sprintf("timer: remove item key: %s, expire: %s, index: %d", bb.Key, bb.ExpiredString(), bb.index)
        common.Debug(format)
    }
    return
}

func (tb *TimeBasic) start() {
    for {
        tb.expire()
        <-tb.signal.C
    }
}

func (tb *TimeBasic) expire() {
    var (
        fn TimerFunc
        bb *BlockBasic
        d  time.Duration
    )

    tb.lock.Lock()

    for {
       if len(tb.timers) == 0 {
           d = infiniteDuration
           if false {
               common.Debug("timer: no other instance")
           }
           break
       }

        bb = tb.timers[0]

        if !bb.Expired() {
            break
        }
        fn = bb.fn
        tb.remove(bb)
        tb.lock.Unlock()

        if fn == nil {

            common.Warn("expire timer no fn")
        } else {
            if false {
                format := fmt.Sprintf("timer key: %s, expire: %s, index: %d expired, call fn", bb.Key, bb.ExpiredString(), bb.index)
                common.Debug(format)
            }
            fn()
        }
        tb.lock.Lock()
    }
    tb.signal.Reset(d)

    if false {
        format := fmt.Sprintf("timer: expire reset delay %d ms", int64(d) / int64(time.Millisecond))
        common.Debug(format)
    }

    tb.lock.Unlock()

    return
}

func (tb *TimeBasic) up(j int) {
    for {
        i := (j - 1) / 2

        if j <= 0 || !tb.less(j, i) {
            break
        }
        tb.swap(i, j)
        j = i
    }
}

func (tb *TimeBasic) down(i, n int) {
            for {
                jl := 2 * i + 1
                if jl >= n || jl < 0 {
            break
        }
        j := jl
        jr := jl + 1
        if jr < n && !tb.less(jl, jr) {
            j = jr  // =2*i + 2
        }
        if !tb.less(j, i) {
            break
        }
        tb.swap(i, j)
        i = j
    }
}

func (tb *TimeBasic) less(i, j int) bool {
    return tb.timers[i].expire.Before(tb.timers[j].expire)
}

func (tb *TimeBasic) swap(i, j int) {
    tb.timers[i], tb.timers[j] = tb.timers[j], tb.timers[i]
    tb.timers[i].index = i
    tb.timers[j].index = j
}