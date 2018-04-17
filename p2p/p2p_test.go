package p2p

import (
	"fmt"
	"testing"
	"github.com/magiconair/properties/assert"
	"github.com/iost-official/prototype/core"
)

func TestNetwork(t *testing.T) {
	nn, err := NewNaiveNetwork(3)
	if err != nil {
		t.Errorf("NewNaiveNetwork encounter err %+v", err)
		return
	}
	lis1, err := nn.Listen(11037)
	if err != nil {
		fmt.Println(err)
	}
	lis2, err := nn.Listen(11038)
	if err != nil {
		fmt.Println(err)
	}
	req := core.Request{
		Time:    1,
		From:    "test1",
		To:      "test2",
		ReqType: 1,
		Body:    []byte{1, 1, 2},
	}
	if err := nn.Broadcast(req); err != nil {
		t.Log("send request encounter err: %+v\n", err)
	}

	message := <-lis1
	assert.Equal(t, message, req)

	message = <-lis2
	assert.Equal(t, message, req)
}
