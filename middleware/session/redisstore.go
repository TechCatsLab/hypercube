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
 *     Initial: 2017/04/08        Feng Yifei
 */

package session

import (
	"bytes"
	"encoding/base32"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// cookies/redis key 默认失效时间
var defaultExpire = 86400		// 24 小时

type Serializer interface {
	Deserialize(d []byte, ss *sessions.Session) error
	Serialize(ss *sessions.Session) ([]byte, error)
}

type JSONSerializer struct{}

func (s JSONSerializer) Serialize(ss *sessions.Session) ([]byte, error) {
	m := make(map[string]interface{}, len(ss.Values))

	for k, v := range ss.Values {
		ks, ok := k.(string)
		if !ok {
			err := fmt.Errorf("Non-string key value, cannot serialize session to JSON: %v", k)
			return nil, err
		}
		m[ks] = v
	}

	return json.Marshal(m)
}

func (s JSONSerializer) Deserialize(d []byte, ss *sessions.Session) error {
	m := make(map[string]interface{})

	err := json.Unmarshal(d, &m)
	if err != nil {
		return err
	}

	for k, v := range m {
		ss.Values[k] = v
	}

	return nil
}

type GobSerializer struct{}

func (s GobSerializer) Serialize(ss *sessions.Session) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(ss.Values)
	if err == nil {
		return buf.Bytes(), nil
	}
	return nil, err
}

func (s GobSerializer) Deserialize(d []byte, ss *sessions.Session) error {
	dec := gob.NewDecoder(bytes.NewBuffer(d))
	return dec.Decode(&ss.Values)
}

type RedisStore struct {
	Pool          *redis.Pool
	Codecs        []securecookie.Codec
	Options       *sessions.Options
	DefaultMaxAge int
	maxLength     int
	keyPrefix     string
	serializer    Serializer
}

// 0： 无限制
func (s *RedisStore) SetMaxLength(l int) {
	if l >= 0 {
		s.maxLength = l
	}
}

// 设置 key 前缀
func (s *RedisStore) SetKeyPrefix(p string) {
	s.keyPrefix = p
}

func (s *RedisStore) SetSerializer(ss Serializer) {
	s.serializer = ss
}

// 0: 无限制
func (s *RedisStore) SetMaxAge(v int) {
	var (
		c *securecookie.SecureCookie
		ok bool
	)

	s.Options.MaxAge = v

	for i := range s.Codecs {
		if c, ok = s.Codecs[i].(*securecookie.SecureCookie); ok {
			c.MaxAge(v)
		}
	}
}

func dial(network, address, password string) (redis.Conn, error) {
	c, err := redis.Dial(network, address)

	if err != nil {
		return nil, err
	}

	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, err
}

func NewRediStore(size int, network, address, password string, keyPairs ...[]byte) (*RedisStore, error) {
	return NewRediStoreWithPool(&redis.Pool{
		MaxIdle:     size,
		IdleTimeout: 300 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			return dial(network, address, password)
		},
	}, keyPairs...)
}

func dialWithDB(network, address, password, DB string) (redis.Conn, error) {
	c, err := dial(network, address, password)

	if err != nil {
		return nil, err
	}

	if _, err := c.Do("SELECT", DB); err != nil {
		c.Close()
		return nil, err
	}
	return c, err
}

func NewRediStoreWithDB(size int, network, address, password, DB string, keyPairs ...[]byte) (*RedisStore, error) {
	return NewRediStoreWithPool(&redis.Pool{
		MaxIdle:     size,
		IdleTimeout: 300 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			return dialWithDB(network, address, password, DB)
		},
	}, keyPairs...)
}

func NewRediStoreWithPool(pool *redis.Pool, keyPairs ...[]byte) (*RedisStore, error) {
	rs := &RedisStore{
		Pool:   pool,
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: defaultExpire,
		},
		DefaultMaxAge: 20 * 60,			// 20 分钟
		maxLength:     4096,
		keyPrefix:     "session_",
		serializer:    GobSerializer{},
	}
	_, err := rs.ping()
	return rs, err
}

func (s *RedisStore) Close() error {
	return s.Pool.Close()
}

func (s *RedisStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

func (s *RedisStore) New(r *http.Request, name string) (*sessions.Session, error) {
	var (
		err error
		ok  bool
	)
	session := sessions.NewSession(s, name)

	options := *s.Options
	session.Options = &options
	session.IsNew = true

	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.Codecs...)
		if err == nil {
			ok, err = s.load(session)
			session.IsNew = !(err == nil && ok)
		}
	}
	return session, err
}

func (s *RedisStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	if session.Options.MaxAge < 0 {
		if err := s.delete(session); err != nil {
			return err
		}

		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
	} else {
		if session.ID == "" {
			session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		}

		if err := s.save(session); err != nil {
			return err
		}

		encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.Codecs...)

		if err != nil {
			return err
		}

		http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	}
	return nil
}

func (s *RedisStore) Delete(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	conn := s.Pool.Get()
	defer conn.Close()

	if _, err := conn.Do("DEL", s.keyPrefix + session.ID); err != nil {
		return err
	}

	options := *session.Options
	options.MaxAge = -1
	http.SetCookie(w, sessions.NewCookie(session.Name(), "", &options))

	for k := range session.Values {
		delete(session.Values, k)
	}

	return nil
}

func (s *RedisStore) ping() (bool, error) {
	conn := s.Pool.Get()
	defer conn.Close()

	data, err := conn.Do("PING")
	if err != nil || data == nil {
		return false, err
	}
	return (data == "PONG"), nil
}

func (s *RedisStore) save(session *sessions.Session) error {
	b, err := s.serializer.Serialize(session)

	if err != nil {
		return err
	}

	if s.maxLength != 0 && len(b) > s.maxLength {
		return errors.New("SessionStore: the value to store is too big")
	}

	conn := s.Pool.Get()
	defer conn.Close()

	if err = conn.Err(); err != nil {
		return err
	}

	age := session.Options.MaxAge
	if age == 0 {
		age = s.DefaultMaxAge
	}

	_, err = conn.Do("SETEX", s.keyPrefix+session.ID, age, b)
	return err
}

func (s *RedisStore) load(session *sessions.Session) (bool, error) {
	conn := s.Pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		return false, err
	}

	data, err := conn.Do("GET", s.keyPrefix + session.ID)

	if err != nil {
		return false, err
	}

	if data == nil {
		return false, nil
	}

	b, err := redis.Bytes(data, err)
	if err != nil {
		return false, err
	}

	return true, s.serializer.Deserialize(b, session)
}

func (s *RedisStore) delete(session *sessions.Session) error {
	conn := s.Pool.Get()
	defer conn.Close()

	if _, err := conn.Do("DEL", s.keyPrefix + session.ID); err != nil {
		return err
	}
	return nil
}
