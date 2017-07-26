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
 *     Initial: 2017/07/05        Feng Yifei
 *     Modify:  2017/07/07        Liu  Jiachang
 */

package message

const (
	// MessageTypeLogin - Login
	MessageTypeLogin = 0x0001

	// MessageTypeLogout - Logout
	MessageTypeLogout = 0x0002

	// MessageTypeKeepAlive - Client to server keepalive message
	MessageTypeKeepAlive = 0x0003

	// MessageTypePlainText - Plain text message
	MessageTypePlainText = 0x0100

	// MessageTypeEmotion - Emotion text message
	MessageTypeEmotion = 0x0101

	// MessageTypeGroupPlainText - Group plain text message
	MessageTypeGroupPlainText = 0x0200

	// MessageTypePushPlainText - Push plain text message
	MessageTypePushPlainText = 0x0300

	// PushToAll - Push plain text to all
	PushToAll = 0x0301

	// PushToPart - Push plain text to selected users
	PushToPart = 0x0302
)
