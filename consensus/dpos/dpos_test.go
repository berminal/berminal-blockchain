package dpos

import (
	"testing"

	. "github.com/bouk/monkey"
	. "github.com/golang/mock/gomock"

	"bytes"
	"time"

	"github.com/iost-official/prototype/account"
	"github.com/iost-official/prototype/common"
	"github.com/iost-official/prototype/consensus/common"
	"github.com/iost-official/prototype/core/block"
	"github.com/iost-official/prototype/core/message"
	"github.com/iost-official/prototype/core/mocks"
	"github.com/iost-official/prototype/core/state"
	"github.com/iost-official/prototype/core/tx"
	"github.com/iost-official/prototype/network"
	"github.com/iost-official/prototype/network/mocks"
	"github.com/iost-official/prototype/verifier"
	"github.com/iost-official/prototype/vm"
	"github.com/iost-official/prototype/vm/lua"
	. "github.com/smartystreets/goconvey/convey"
	"fmt"
)

func TestNewDPoS(t *testing.T) {
	Convey("Test of NewDPos", t, func() {
		mockCtr := NewController(t)
		mockRouter := protocol_mock.NewMockRouter(mockCtr)
		mockBc := core_mock.NewMockChain(mockCtr)
		mockPool := core_mock.NewMockPool(mockCtr)
		mockPool.EXPECT().Copy().Return(mockPool).AnyTimes()
		mockPool.EXPECT().PutHM(Any(), Any(), Any()).AnyTimes().Return(nil)

		network.Route = mockRouter
		//获取router实例
		guard := Patch(network.RouterFactory, func(_ string) (network.Router, error) {
			return mockRouter, nil
		})
		defer guard.Unpatch()

		heightChan := make(chan message.Message, 1)
		blkChan := make(chan message.Message, 1)
		mockRouter.EXPECT().FilteredChan(Any()).Return(heightChan, nil)
		mockRouter.EXPECT().FilteredChan(Any()).Return(blkChan, nil)

		// 设置第一个通道txchan
		txChan := make(chan message.Message, 1)
		mockRouter.EXPECT().FilteredChan(Any()).Return(txChan, nil)

		// 设置第二个通道Blockchan
		blockChan := make(chan message.Message, 1)
		mockRouter.EXPECT().FilteredChan(Any()).Return(blockChan, nil)

		// 创世块的询问和插入
		var getNumber uint64
		var pushNumber int64
		mockBc.EXPECT().GetBlockByNumber(Any()).Do(func(number uint64) *block.Block {
			getNumber = number
			return nil
		})
		mockBc.EXPECT().Push(Any()).Do(func(block *block.Block) error {
			pushNumber = block.Head.Number
			return nil
		})

		p, _ := NewDPoS(account.Account{"id0", []byte("pubkey"), []byte("seckey")}, mockBc, mockPool, []string{})

		So(p.Account.GetId(), ShouldEqual, "id0")
		So(getNumber, ShouldEqual, 0)
		So(pushNumber, ShouldEqual, 0)
	})

}

func TestRunGenerateBlock(t *testing.T) {
	Convey("Test of Run (Generate Block)", t, func() {
		mockCtr := NewController(t)
		mockRouter := protocol_mock.NewMockRouter(mockCtr)
		mockBc := core_mock.NewMockChain(mockCtr)
		mockBc.EXPECT().HasTx(Any()).AnyTimes().Return(false, nil)
		mockPool := core_mock.NewMockPool(mockCtr)

		mockPool.EXPECT().Copy().Return(mockPool).AnyTimes()
		mockPool.EXPECT().PutHM(Any(), Any(), Any()).AnyTimes().Return(nil)

		mockBc.EXPECT().Iterator().AnyTimes().Return(nil)
		network.Route = mockRouter
		//获取router实例
		guard := Patch(network.RouterFactory, func(_ string) (network.Router, error) {
			return mockRouter, nil
		})
		defer guard.Unpatch()

		guard1 := Patch(consensus_common.VerifyTxSig, func(_ tx.Tx) bool {
			return true
		})
		defer guard1.Unpatch()

		guard2 := Patch(consensus_common.VerifyTx, func(_ *tx.Tx, _ *verifier.CacheVerifier) (state.Pool, bool) {
			return nil, true
		})
		defer guard2.Unpatch()

		guard3 := Patch(consensus_common.StdBlockVerifier, func(_ *block.Block, _ state.Pool) (state.Pool, error) {
			return nil, nil
		})
		defer guard3.Unpatch()

		heightChan := make(chan message.Message, 1)
		blkSyncChan := make(chan message.Message, 1)
		mockRouter.EXPECT().FilteredChan(Any()).Return(heightChan, nil)
		mockRouter.EXPECT().FilteredChan(Any()).Return(blkSyncChan, nil)

		//设置第一个通道txchan
		txChan := make(chan message.Message, 5)
		mockRouter.EXPECT().FilteredChan(Any()).Return(txChan, nil)

		//设置第二个通道Blockchan
		blockChan := make(chan message.Message, 1)
		mockRouter.EXPECT().FilteredChan(Any()).Return(blockChan, nil)

		mockBc.EXPECT().GetBlockByNumber(Eq(uint64(0))).Return(nil)
		var genesis *block.Block
		mockBc.EXPECT().Push(Any()).Do(func(block *block.Block) error {
			genesis = block
			return nil
		})
		seckey := common.Sha256([]byte("SeckeyId0"))
		pubkey := common.CalcPubkeyInSecp256k1(seckey)
		p, _ := NewDPoS(account.Account{"id0", pubkey, seckey}, mockBc, mockPool, []string{"id0", "id1", "id2"})

		main := lua.NewMethod("main", 0, 1)
		code := `function main()
						Put("hello", "world")
						return "success"
					end`
		lc := lua.NewContract(vm.ContractInfo{Prefix: "test", GasLimit: 100, Price: 1, Publisher: vm.IOSTAccount("ahaha")}, code, main)
		//构造测试数据
		newTx := tx.NewTx(0, &lc)
		txChan <- message.Message{
			Time:    20180426111111,
			From:    "0xaaaaaaaaaaaaa",
			To:      "0xbbbbbbbbbbbb",
			ReqType: 1,
			Body:    newTx.Encode()}

		newTx = tx.NewTx(1, &lc)
		txChan <- message.Message{
			Time:    20180426111111,
			From:    "0xaaaaaaaaaaaaa",
			To:      "0xbbbbbbbbbbbb",
			ReqType: 1,
			Body:    newTx.Encode()}

		newTx = tx.NewTx(2, &lc)
		txChan <- message.Message{
			Time:    20180426111111,
			From:    "0xaaaaaaaaaaaaa",
			To:      "0xbbbbbbbbbbbb",
			ReqType: 1,
			Body:    newTx.Encode()}

		mockBc.EXPECT().Top().Return(genesis).AnyTimes()

		var blk block.Block
		var reqType network.ReqType
		mockRouter.EXPECT().Broadcast(Any()).Do(func(req message.Message) {
			reqType = network.ReqType(req.ReqType)
			blk.Decode(req.Body)
		}).AnyTimes()

		p.Run()

		time.Sleep(time.Second * 2)
		So(reqType, ShouldEqual, network.ReqNewBlock)
		So(blk.Head.Number, ShouldEqual, 1)
		So(string(blk.Head.ParentHash), ShouldEqual, string(genesis.Head.Hash()))
		So(blk.Head.Witness, ShouldEqual, "id0")

		p.Stop()

	})

}

func TestRunReceiveBlock(t *testing.T) {
	Convey("Test of Run (Receive Block)", t, func() {
		mockCtr := NewController(t)
		mockRouter := protocol_mock.NewMockRouter(mockCtr)
		mockBc := core_mock.NewMockChain(mockCtr)
		mockPool := core_mock.NewMockPool(mockCtr)
		mockPool.EXPECT().Copy().Return(mockPool).AnyTimes()
		mockPool.EXPECT().PutHM(Any(), Any(), Any()).AnyTimes().Return(nil)

		network.Route = mockRouter
		//获取router实例
		guard := Patch(network.RouterFactory, func(_ string) (network.Router, error) {
			return mockRouter, nil
		})

		defer guard.Unpatch()

		heightChan := make(chan message.Message, 1)
		blkSyncChan := make(chan message.Message, 1)
		mockRouter.EXPECT().FilteredChan(Any()).Return(heightChan, nil)
		mockRouter.EXPECT().FilteredChan(Any()).Return(blkSyncChan, nil)

		//设置第一个通道txchan
		txChan := make(chan message.Message, 1)
		mockRouter.EXPECT().FilteredChan(Any()).Return(txChan, nil)

		//设置第二个通道Blockchan
		blkChan := make(chan message.Message, 1)
		mockRouter.EXPECT().FilteredChan(Any()).Return(blkChan, nil)

		mockBc.EXPECT().GetBlockByNumber(Eq(uint64(0))).Return(nil)
		var genesis *block.Block
		mockBc.EXPECT().Push(Any()).Do(func(block *block.Block) error {
			genesis = block
			return nil
		})
		seckey := common.Sha256([]byte("SeckeyId1"))
		pubkey := common.CalcPubkeyInSecp256k1(seckey)
		p, _ := NewDPoS(account.Account{"id1", pubkey, seckey}, mockBc, mockPool, []string{"id0", "id1", "id2"})

		main := lua.NewMethod("main", 0, 1)
		code := `function main()
						Put("hello", "world")
						return "success"
					end`
		lc := lua.NewContract(vm.ContractInfo{Prefix: "test", GasLimit: 100, Price: 1, Publisher: vm.IOSTAccount("ahaha")}, code, main)
		newTx := tx.NewTx(0, &lc)
		//构造测试数据
		txChan <- message.Message{
			Time:    20180426111111,
			From:    "0xaaaaaaaaaaaaa",
			To:      "0xbbbbbbbbbbbb",
			ReqType: 1,
			Body:    newTx.Encode()}

		mockBc.EXPECT().Top().Return(genesis).AnyTimes()
		mockBc.EXPECT().Length().Return(uint64(0)).AnyTimes()
		//mockPool.EXPECT().Copy().Return(nil).AnyTimes()

		blk, msg := generateTestBlockMsg("id0", "seckeyId0", 1, genesis.Head.Hash())
		blkChan <- msg

		p.Run()

		time.Sleep(time.Second * 1)
		So(blk.Head.Number, ShouldEqual, 1)
		So(string(blk.Head.ParentHash), ShouldEqual, string(genesis.Head.Hash()))
		So(blk.Head.Witness, ShouldEqual, "id0")

		p.Stop()

	})

}

func TestRunMultipleBlocks(t *testing.T) {
	Convey("Test of Run (Multiple Blocks)", t, func() {
		mockCtr := NewController(t)
		mockRouter := protocol_mock.NewMockRouter(mockCtr)
		mockBc := core_mock.NewMockChain(mockCtr)
		mockPool := core_mock.NewMockPool(mockCtr)
		mockPool.EXPECT().Copy().Return(mockPool).AnyTimes()
		mockPool.EXPECT().PutHM(Any(), Any(), Any()).AnyTimes().Return(nil)

		mockBc.EXPECT().Iterator().AnyTimes().Return(nil)

		network.Route = mockRouter
		//获取router实例
		guard := Patch(network.RouterFactory, func(_ string) (network.Router, error) {
			return mockRouter, nil
		})

		defer guard.Unpatch()

		heightChan := make(chan message.Message, 1)
		blkSyncChan := make(chan message.Message, 1)
		mockRouter.EXPECT().FilteredChan(Any()).Return(heightChan, nil)
		mockRouter.EXPECT().FilteredChan(Any()).Return(blkSyncChan, nil)

		//设置第一个通道txchan
		txChan := make(chan message.Message, 1)
		mockRouter.EXPECT().FilteredChan(Any()).Return(txChan, nil)

		//设置第二个通道Blockchan
		blkChan := make(chan message.Message, 1)
		mockRouter.EXPECT().FilteredChan(Any()).Return(blkChan, nil)

		mockBc.EXPECT().GetBlockByNumber(Eq(uint64(0))).Return(nil)
		var genesis *block.Block
		mockBc.EXPECT().Push(Any()).Do(func(block *block.Block) error {
			genesis = block
			return nil
		})
		seckey := common.Sha256([]byte("SeckeyId1"))
		pubkey := common.CalcPubkeyInSecp256k1(seckey)
		p, _ := NewDPoS(account.Account{"id1", pubkey, seckey}, mockBc, mockPool, []string{"id0", "id1", "id2"})

		main := lua.NewMethod("main", 0, 1)
		code := `function main()
						Put("hello", "world")
						return "success"
					end`
		lc := lua.NewContract(vm.ContractInfo{Prefix: "test", GasLimit: 100, Price: 1, Publisher: vm.IOSTAccount("ahaha")}, code, main)
		newTx := tx.NewTx(0, &lc)
		//构造测试数据
		txChan <- message.Message{
			Time:    20180426111111,
			From:    "0xaaaaaaaaaaaaa",
			To:      "0xbbbbbbbbbbbb",
			ReqType: 1,
			Body:    newTx.Encode(),
		}

		mockBc.EXPECT().Top().Return(genesis).AnyTimes()
		mockBc.EXPECT().Length().Return(uint64(0)).AnyTimes()
		mockPool.EXPECT().Copy().Return(nil).AnyTimes()
		mockPool.EXPECT().Flush().Return(nil).AnyTimes()
		mockRouter.EXPECT().CancelDownload(Any(), Any()).Return(nil).AnyTimes()

		Convey("correct blocks", func() {
			blk, msg := generateTestBlockMsg("id0", "SeckeyId0", 1, genesis.Head.Hash())
			blkChan <- msg

			var reqType network.ReqType
			var reqBlk block.Block
			mockRouter.EXPECT().Broadcast(Any()).Do(func(req message.Message) {
				reqType = network.ReqType(req.ReqType)
				reqBlk.Decode(req.Body)
			}).AnyTimes()

			var pushed int64
			mockBc.EXPECT().Push(Any()).Do(func(block *block.Block) error {
				pushed = block.Head.Number
				return nil
			})
			p.Run()

			time.Sleep(time.Second)
			// block 1 by id0

			time.Sleep(time.Second * consensus_common.SlotLength)
			// block 2 by id1, the node itself
			So(reqType, ShouldEqual, network.ReqNewBlock)
			So(reqBlk.Head.Number, ShouldEqual, 2)
			So(string(reqBlk.Head.ParentHash), ShouldEqual, string(blk.Head.Hash()))
			So(reqBlk.Head.Witness, ShouldEqual, "id1")

			ts := consensus_common.GetCurrentTimestamp()
			ts.Add(1)
			len := ts.ToUnixSec() - time.Now().Unix()
			time.Sleep(time.Second * time.Duration(len))
			//blk, msg = generateTestBlockMsg("id2", "SeckeyId2", 3, reqBlk.Head.Hash())
			//blkChan <- msg
			//
			//time.Sleep(time.Second*consensus_common.SlotLength+time.Second*2)
			//fmt.Println("### ")
			//// block 3 by id2
			//So(reqType, ShouldEqual, network.ReqNewBlock)
			//So(bytes.Equal(reqBlk.Head.Hash(), blk.Head.Hash()), ShouldBeTrue)
			//
			//So(pushed, ShouldEqual, 1)

			p.Stop()
		})

		Convey("with fork", func() {
			blk, msg := generateTestBlockMsg("id0", "SeckeyId0", 1, genesis.Head.Hash())
			blkChan <- msg

			var reqType network.ReqType
			var reqBlk block.Block
			mockRouter.EXPECT().Broadcast(Any()).Do(func(req message.Message) {
				reqType = network.ReqType(req.ReqType)
				reqBlk.Decode(req.Body)
			}).AnyTimes()

			var pushed int64
			mockBc.EXPECT().Push(Any()).Do(func(block *block.Block) error {
				pushed = block.Head.Number
				return nil
			}).AnyTimes()
			p.Run()

			time.Sleep(time.Second / 2)
			// block 1 by id0

			blk1 := blk

			time.Sleep(time.Second * consensus_common.SlotLength)
			// block 2 by id1, the node itself
			So(reqBlk.Head.Number, ShouldEqual, 2)
			So(string(reqBlk.Head.ParentHash), ShouldEqual, string(blk.Head.Hash()))
			So(reqBlk.Head.Witness, ShouldEqual, "id1")

			ts := consensus_common.GetCurrentTimestamp()
			ts.Add(1)
			len := ts.ToUnixSec() - time.Now().Unix()
			time.Sleep(time.Second * time.Duration(len))
			blk, msg = generateTestBlockMsg("id2", "SeckeyId2", 2, blk1.Head.Hash())
			blkChan <- msg

			time.Sleep(time.Second / 2)
			// block 2' by id2, is a fork
			//So(bytes.Equal(reqBlk.Head.Hash(), blk.Head.Hash()), ShouldBeTrue)

			ts = consensus_common.GetCurrentTimestamp()
			ts.Add(1)
			len = ts.ToUnixSec() - time.Now().Unix()
			time.Sleep(time.Second * time.Duration(len))
			blk, msg = generateTestBlockMsg("id0", "SeckeyId0", 3, reqBlk.Head.Hash())
			blkChan <- msg
			time.Sleep(time.Second / 2)
			//// block 3 by id0
			//So(bytes.Equal(reqBlk.Head.Hash(), blk.Head.Hash()), ShouldBeTrue)
			//// nothing is pushed until now
			//So(pushed, ShouldEqual, 0)

			time.Sleep(time.Second * consensus_common.SlotLength)
			// block 4 by id1, the node itself
			So(reqBlk.Head.Number, ShouldEqual, 4)
			So(string(reqBlk.Head.ParentHash), ShouldEqual, string(blk.Head.Hash()))
			So(reqBlk.Head.Witness, ShouldEqual, "id1")
			// block 1 and 2 should be pushed
			//So(pushed, ShouldEqual, 2)

			p.Stop()
		})

		Convey("need sync", func() {
			consensus_common.SyncNumber = 2
			p.account.ID = "id3"
			blk1, msg1 := generateTestBlockMsg("id0", "SeckeyId0", 1, genesis.Head.Hash())
			time.Sleep(time.Second * consensus_common.SlotLength)
			blk2, msg2 := generateTestBlockMsg("id1", "SeckeyId1", 2, blk1.Head.Hash())
			time.Sleep(time.Second * consensus_common.SlotLength)
			_, msg3 := generateTestBlockMsg("id2", "SeckeyId2", 3, blk2.Head.Hash())

			blkChan <- msg3

			var bcType network.ReqType
			var bcBlk block.Block
			mockRouter.EXPECT().Broadcast(Any()).Do(func(req message.Message) {
				bcType = network.ReqType(req.ReqType)
				if bcType == network.ReqNewBlock {
					bcBlk.Decode(req.Body)
				}
			}).AnyTimes()

			var pushedBlk *block.Block
			mockBc.EXPECT().Push(Any()).Do(func(block *block.Block) error {
				pushedBlk = block
				return nil
			}).AnyTimes()

			var dlSt, dlEd uint64
			mockRouter.EXPECT().Download(Any(), Any()).Do(func(start, end uint64) error {
				dlSt = start
				dlEd = end
				return nil
			})
			p.Run()

			time.Sleep(time.Second / 2)
			// need sync from 1 to 2
			//So(bcType, ShouldEqual, network.ReqBlockHeight)
			//So(dlSt, ShouldEqual, 1)
			//So(dlEd, ShouldEqual, 3)

			blkChan <- msg2
			time.Sleep(time.Second / 2)

			blkChan <- msg1
			time.Sleep(time.Second / 2)

			// After block1 and block2 received, block 1-3 all set, and block 1 will be pushed
			So(bytes.Equal(pushedBlk.Head.Hash(), blk1.Head.Hash()), ShouldBeTrue)

			p.Stop()
		})
	})
}

func generateTestBlockMsg(witness string, secKeyRaw string, number int64, parentHash []byte) (block.Block, message.Message) {
	blk := block.Block{
		Head: block.BlockHead{
			Number:     number,
			ParentHash: parentHash,
			Witness:    witness,
			Time:       consensus_common.GetCurrentTimestamp().Slot,
		},
	}
	headInfo := generateHeadInfo(blk.Head)
	sig, _ := common.Sign(common.Secp256k1, headInfo, common.Sha256([]byte(secKeyRaw)))
	blk.Head.Signature = sig.Encode()
	msg := message.Message{
		Time:    time.Now().Unix(),
		From:    "",
		To:      "",
		ReqType: int32(network.ReqNewBlock),
		Body:    blk.Encode(),
	}
	return blk, msg
}

//go test -bench=. -benchmem -run=nonce
func BenchmarkAddBlockCache(b *testing.B) { 
	//benchAddBlockCache(b,true)
	benchAddBlockCache(b,false)
}
func BenchmarkGetBlockCache(b *testing.B) { 
//	benchGetBlockCache(b,true) 
	benchGetBlockCache(b,false) 
}
func BenchmarkBlockVerifier(b *testing.B) { benchBlockVerifier(b) }
func BenchmarkTxPool(b *testing.B)        { benchTxPool(b) }
func BenchmarkBlockHead(b *testing.B)     { benchBlockHead(b) }
func BenchmarkGenerateBlock(b *testing.B) { benchGenerateBlock(b) }


func envInit(b *testing.B) (*DPoS,[]account.Account,[]string) {
	var accountList []account.Account

	var witnessList []string

	for i:=0;i<3 ;i++  {
		account, err:=account.NewAccount(nil)
		if err!=nil{
			panic("account.NewAccount error")
		}
		accountList = append(accountList, account)
		witnessList = append(witnessList, account.ID)
	}

	tx.LdbPath = ""

	mockCtr := NewController(b)
	mockRouter := protocol_mock.NewMockRouter(mockCtr)

	network.Route = mockRouter
	//获取router实例
	guard := Patch(network.RouterFactory, func(_ string) (network.Router, error) {
		return mockRouter, nil
	})

	heightChan := make(chan message.Message, 1)
	blkSyncChan := make(chan message.Message, 1)
	mockRouter.EXPECT().FilteredChan(Any()).Return(heightChan, nil)
	mockRouter.EXPECT().FilteredChan(Any()).Return(blkSyncChan, nil)

	//设置第一个通道txchan
	txChan := make(chan message.Message, 1)
	mockRouter.EXPECT().FilteredChan(Any()).Return(txChan, nil)

	//设置第二个通道Blockchan
	blkChan := make(chan message.Message, 1)
	mockRouter.EXPECT().FilteredChan(Any()).Return(blkChan, nil)

	defer guard.Unpatch()

	txDb := tx.TxDbInstance()
	if txDb == nil {
		panic("txDB error")
	}

	blockChain, err := block.Instance()
	if err != nil {
		panic("block.Instance error")
	}

	err = state.PoolInstance()
	if err != nil {
		panic("state.PoolInstance error")
	}

	//verifyFunc := func(blk *block.Block, parent *block.Block, pool state.Pool) (state.Pool, error) {
	//	return pool, nil
	//}

	//blockCache := consensus_common.NewBlockCache(blockChain, state.StdPool, len(witnessList)*2/3)
	//seckey := common.Sha256([]byte("SeckeyId0"))
	//pubkey := common.CalcPubkeyInSecp256k1(seckey)
	p, err := NewDPoS(accountList[0], blockChain, state.StdPool, witnessList)
	if err != nil {
		b.Errorf("NewDPoS error")
	}
	return p,accountList,witnessList
}
// cache中添加block性能测试
func genBlocks(p *DPoS,accountList []account.Account,witnessList []string,n int,continuity bool) (blockPool []*block.Block) {
	confChain := p.blockCache.BlockChain()
	tblock := confChain.Top() //获取创世块

	blockLen := p.blockCache.ConfirmedLength()
	fmt.Println(blockLen)

	//blockNum := 1000
	slot := consensus_common.GetCurrentTimestamp().Slot

	for i := 0; i < n; i++ {
 		var hash []byte
		if len(blockPool) == 0 {
			//用创世块的头hash赋值
			hash =tblock.Head.Hash()
		} else {
			hash = blockPool[len(blockPool)-1].Head.Hash()
		}
		//make every block has no parent
		if continuity==false{
			hash[i%len(hash)]=byte(i%256)
		}
		blk := block.Block{Content: []tx.Tx{}, Head: block.BlockHead{
			Version:    0,
			ParentHash: hash,
			TreeHash:   make([]byte, 0),
			BlockHash:  make([]byte, 0),
			Info:       []byte("test"),
			Number:     int64(i+1),
			Witness:    witnessList[i%3],
			Time:       slot + int64(i),
		}}

		headInfo := generateHeadInfo(blk.Head)
		sig, _ := common.Sign(common.Secp256k1, headInfo, accountList[i%3].Seckey)
		blk.Head.Signature = sig.Encode()

		blockPool = append(blockPool, &blk)
	}
	return
}
func benchAddBlockCache(b *testing.B,continuity bool) {
	
	p,accountList,witnessList:=envInit(b)
	//生成block
	blockPool:=genBlocks(p,accountList,witnessList,b.N,continuity)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, bl := range blockPool {
			p.blockCache.Add(bl, p.blockVerify)
		}
	}

}

// cache中获取block性能测试
func benchGetBlockCache(b *testing.B,continuity bool) {
	p,accountList,witnessList:=envInit(b)
	//生成block
	blockPool:=genBlocks(p,accountList,witnessList,b.N,continuity)
	for i := 0; i < b.N; i++ {
		for _, bl := range blockPool {
			p.blockCache.Add(bl, p.blockVerify)
		}
	}

	//get block
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
	for _, bl := range blockPool {
			p.blockCache.Add(bl, p.blockVerify)
		}
	}
}

// block验证性能测试
func benchBlockVerifier(b *testing.B) {

}

// 交易添加Pool性能测试
func benchTxPool(b *testing.B) {

}

// 生成block head性能测试
func benchBlockHead(b *testing.B) {

}

// 生成块性能测试
func benchGenerateBlock(b *testing.B) {

}
