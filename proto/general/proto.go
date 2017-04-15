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
 *     Initial: 2017/04/02        Liu Jiachang
 */

package general

import (
	"fmt"
	"encoding/json"
	"github.com/gorilla/websocket"
)

var (
	CurVer     = uint16(0x0001) // 从左向右 第一位 主版本号 0, 第二位 副版本号 0, 第三四位 次版本号
	emptyProto = Proto{}
)

type Proto struct {
	Ver     uint16              `json:"v"`
	Type    uint32              `json:"t"`
	Body    json.RawMessage     `json:"b"`
}

func (p *Proto) Reset() {
	*p = emptyProto
}

func (p *Proto) String() string {
	return fmt.Sprintf("\n-------- proto --------\nver: %d\ntype: %d\nbody: %v\n--------------------", p.Ver, p.Type, p.Body)
}

func (p *Proto) VerCheck() (*Proto, uint32) {
	if p.Ver > CurVer {
		p.Type = AccTypeVerConf
		return  p, ErrVerConf
	}

	if p.Ver < CurVer {
		p.Type = AccTypeSpiteAtt
		return  p, ErrVerConf
	}

	return p, ErrSucceed
}

func (p *Proto) ReadWebSocket(wr *websocket.Conn) error {
	return wr.ReadJSON(p)
}
