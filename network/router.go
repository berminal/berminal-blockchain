package network

import (
	"fmt"

	"github.com/iost-official/prototype/core/message"
)

//ReqType Marked request types using by protocol

//go:generate mockgen -destination network/mocks/mock_router.go -package protocol_mock github.com/iost-official/prototype/network Router
type ReqType int32

const (
	ReqPublishTx ReqType = iota
	ReqNewBlock
)

//Router Forwarding specific request to other components and sending messages for them
type Router interface {
	Init(base Network, port uint16) error
	FilteredChan(filter Filter) (chan message.Message, error)
	Run()
	Stop()
	Send(req message.Message)
	Broadcast(req message.Message)
	Download(req message.Message) chan []byte
}

func RouterFactory(target string) (Router, error) {
	switch target {
	case "base":
		return &RouterImpl{}, nil
	}
	return nil, fmt.Errorf("target Router not found")
}

type RouterImpl struct {
	base Network

	chIn  <-chan message.Message
	chOut chan<- message.Message

	filterList  []Filter
	filterMap   map[int]chan message.Message
	knownMember []string
	ExitSignal  chan bool
}

func (r *RouterImpl) Init(base Network, port uint16) error {
	var err error
	r.base = base
	r.filterList = make([]Filter, 0)
	r.filterMap = make(map[int]chan message.Message)
	r.knownMember = make([]string, 0)
	r.ExitSignal = make(chan bool)
	r.chIn, err = r.base.Listen(port)
	if err != nil {
		return err
	}

	return nil
}

//FilteredChan Get filtered request channel
func (r *RouterImpl) FilteredChan(filter Filter) (chan message.Message, error) {
	chReq := make(chan message.Message)

	r.filterList = append(r.filterList, filter)
	r.filterMap[len(r.filterList)-1] = chReq

	return chReq, nil
}

func (r *RouterImpl) receiveLoop() {
	for true {
		select {
		case <-r.ExitSignal:
			return
		case req := <-r.chIn:
			for i, f := range r.filterList {
				if f.check(req) {
					r.filterMap[i] <- req
				}
			}
		}
	}
}

func (r *RouterImpl) Run() {
	go r.receiveLoop()
}

func (r *RouterImpl) Stop() {
	r.ExitSignal <- true
}

func (r *RouterImpl) Send(req message.Message) {
	r.base.Send(req)
}

// Broadcast to all known members
func (r *RouterImpl) Broadcast(req message.Message) {
	for _, to := range r.knownMember {
		req.To = to

		go func() {
			r.Send(req)
		}()
	}
}

func (r *RouterImpl) Download(req message.Message) chan []byte {
	return nil // TODO 实现
}

//Filter The filter used by Router
// Rulers :
//     1. if both white list and black list are nil, this filter is all-pass
//     2. if one of those is not nil, filter as it is
//     3. if both of those list are not nil, filter as white list
type Filter struct {
	WhiteList  []message.Message
	BlackList  []message.Message
	RejectType []ReqType
	AcceptType []ReqType
}

func (f *Filter) check(req message.Message) bool {
	var memberCheck, typeCheck byte
	if f.WhiteList == nil && f.BlackList == nil {
		memberCheck = byte(0)
	} else if f.WhiteList != nil {
		memberCheck = byte(1)
	} else {
		memberCheck = byte(2)
	}
	if f.AcceptType == nil && f.RejectType == nil {
		typeCheck = byte(0)
	} else if f.AcceptType != nil {
		typeCheck = byte(1)
	} else {
		typeCheck = byte(2)
	}

	var m, t bool

	switch memberCheck {
	case 0:
		m = true
	case 1:
		m = memberContain(req.From, f.WhiteList)
	case 2:
		m = !memberContain(req.From, f.BlackList)
	}

	switch typeCheck {
	case 0:
		t = true
	case 1:
		t = reqTypeContain(req.ReqType, f.AcceptType)
	case 2:
		t = !reqTypeContain(req.ReqType, f.RejectType)
	}

	return m && t
}

func memberContain(a string, c []message.Message) bool {
	for _, m := range c {
		if m.From == a {
			return true
		}
	}
	return false
}

func reqTypeContain(a int32, c []ReqType) bool {
	for _, t := range c {
		if int32(t) == a {
			return true
		}
	}
	return false

}
