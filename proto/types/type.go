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
 *     Initial: 2017/04/04        Feng Yifei
 *     Modify:  2017/06/07        Yang Chenglong     添加AccessHeart类型
 *     Modify: 2017/06/08         Yang Chenglong     整理type常量
 */

package types

const (
	ApiTypeUserLogin                = 0x1000
	ApiTypeUserLogout               = 0x1001
	ApiTypeUserOnConnect            = 0x1002
	ApiTypeUserOnDisConnect         = 0x1003

	ApiTypeAccessInfo               = 0x1004
	ApiTypeAccessHeart              = 0x1008

	GeneralTypeKeepAlive  			= uint32(0x2000)            // heartbeat
	GeneralTypeVersionConflict      = uint32(0x2001)            // proto version conflict
	GeneralTypeLogin           		= uint32(0x2002)            // login access type
	GeneralTypeLogout          		= uint32(0x2003)            // logout access type

	GeneralTypeTextMsg         		= uint32(0x2004)            // user to user message

	PushTypePushMsg            		= uint32(0x4000)            // logic server push message to all people

)
