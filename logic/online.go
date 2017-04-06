package main

import (
	"hypercube/proto/api"
	"errors"
)

var (
	OnLineUserMag *OnLineUserMagServer

	ParamErr     = errors.New("Your input parametric error!")
	UserNotExist = errors.New("User not exist!")
)

func init() {
	OnLineUserMag = &OnLineUserMagServer{
		users:      make(map[uint64]string),
		addchan:    make(chan api.UserLogin),
		rmchan:     make(chan uint64),
		qrchan:     make(chan uint64),
		replychan:  make(chan chanreply),
	}

	OnLineUserMag.loop()
}

type OnLineUserMagServer struct {
	users       map[uint64]string
	addchan     chan api.UserLogin
	rmchan      chan uint64
	qrchan      chan uint64
	replychan   chan chanreply
}

type chanreply struct {
	ServerIP     string
	Err          error
}

func (this *OnLineUserMagServer) Add(user api.UserLogin) error {
	if user.ServerIP != "" && user.UserID != 0 {
		this.addchan <- user
		repl := <-this.replychan

		return repl.Err
	}

	return ParamErr
}

func (this *OnLineUserMagServer) Remove(uid uint64) error {
	if uid != 0 {
		this.rmchan <- uid
		repl := <-this.replychan

		return repl.Err
	}

	return ParamErr
}

func (this *OnLineUserMagServer) Query(uid uint64) (string, error) {
	if uid != 0 {
		this.qrchan <- uid
		repl := <-this.replychan

		return repl.ServerIP, repl.Err
	}

	return "", ParamErr
}

func (this *OnLineUserMagServer)loop() {
	for {
		repl := chanreply{
			ServerIP:    "",
			Err:         nil,
		}

		select {
		case user := <-this.addchan:
			this.users[user.UserID] = user.ServerIP

			this.replychan <- repl
		case uid := <-this.rmchan:
			if _, ok := this.users[uid]; ok {
				delete(this.users, uid)

				this.replychan <- repl
			} else {
				repl.Err = ParamErr
				this.replychan <- repl
			}
		case uid := <-this.qrchan:
			if serverip, ok := this.users[uid]; ok {
				repl.ServerIP = serverip

				this.replychan <- repl
			} else {
				repl.Err = ParamErr
				this.replychan <- repl
			}
		}

	}
}
