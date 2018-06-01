package dpos

import (
	"bytes"
	"encoding/binary"

	. "github.com/iost-official/prototype/account"
	. "github.com/iost-official/prototype/consensus/common"
	. "github.com/iost-official/prototype/core/tx"
	. "github.com/iost-official/prototype/network"

	"errors"
	"fmt"
	"time"

	"github.com/iost-official/prototype/common"
	"github.com/iost-official/prototype/core/block"
	"github.com/iost-official/prototype/core/message"
	"github.com/iost-official/prototype/core/state"
	"github.com/iost-official/prototype/log"
	"github.com/iost-official/prototype/verifier"
	"github.com/iost-official/prototype/vm"
	"github.com/iost-official/prototype/vm/lua"
)

var TxPerBlk int

type DPoS struct {
	account      Account
	blockCache   BlockCache
	router       Router
	synchronizer Synchronizer
	globalStaticProperty
	globalDynamicProperty

	//测试用，保存投票状态，以及投票消息内容的缓存
	votedStats map[string][]string
	infoCache  [][]byte

	exitSignal chan struct{}
	ChTx       chan message.Message
	chBlock    chan message.Message

	log *log.Logger
}

// NewDPoS: 新建一个DPoS实例
// acc: 节点的Coinbase账户, bc: 基础链(从数据库读取), pool: 基础state池（从数据库读取）, witnessList: 见证节点列表
func NewDPoS(acc Account, bc block.Chain, pool state.Pool, witnessList []string /*, network core.Network*/) (*DPoS, error) {
	TxPerBlk = 3000
	p := DPoS{}
	p.account = acc
	p.blockCache = NewBlockCache(bc, pool, len(witnessList)*2/3)
	if bc.GetBlockByNumber(0) == nil {

		t := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		p.genesis(GetTimestamp(t.Unix()).Slot)
	}

	var err error
	p.router = Route
	if p.router == nil {
		return nil, fmt.Errorf("failed to network.Route is nil")
	}

	p.synchronizer = NewSynchronizer(p.blockCache, p.router)
	if p.synchronizer == nil {
		return nil, err
	}

	//	Tx chan init
	p.ChTx, err = p.router.FilteredChan(Filter{
		AcceptType: []ReqType{
			ReqPublishTx,
			reqTypeVoteTest, // Only for test
		}})
	if err != nil {
		return nil, err
	}

	//	Block chan init
	p.chBlock, err = p.router.FilteredChan(Filter{
		AcceptType: []ReqType{ReqNewBlock}})
	if err != nil {
		return nil, err
	}
	p.exitSignal = make(chan struct{})

	p.log, err = log.NewLogger("consensus.log")
	if err != nil {
		return nil, err
	}

	p.log.NeedPrint = true

	p.initGlobalProperty(p.account, witnessList)

	block := bc.GetBlockByNumber(1)
	if block != nil {
		p.update(&block.Head)
	}

	p.update(&bc.Top().Head)
	return &p, nil
}

func (p *DPoS) initGlobalProperty(acc Account, witnessList []string) {
	p.globalStaticProperty = newGlobalStaticProperty(acc, witnessList)
	p.globalDynamicProperty = newGlobalDynamicProperty()
}

// Run: 运行DPoS实例
func (p *DPoS) Run() {
	p.synchronizer.StartListen()
	go p.txListenLoop()
	go p.blockLoop()
	go p.scheduleLoop()
	//p.genBlock(p.Account, block.Block{})
}

// Stop: 停止DPoS实例
func (p *DPoS) Stop() {
	close(p.ChTx)
	close(p.chBlock)
	close(p.exitSignal)
}

// BlockChain 返回已确认的block chain
func (p *DPoS) BlockChain() block.Chain {
	return p.blockCache.BlockChain()
}

// CachedBlockChain 返回缓存中的最长block chain
func (p *DPoS) CachedBlockChain() block.Chain {
	return p.blockCache.LongestChain()
}

// StatePool 返回已确认的state pool
func (p *DPoS) StatePool() state.Pool {
	return p.blockCache.BasePool()
}

// CacheStatePool 返回缓存中最新的state pool
func (p *DPoS) CachedStatePool() state.Pool {
	return p.blockCache.LongestPool()
}

func (p *DPoS) genesis(initTime int64) error {

	main := lua.NewMethod(vm.Public, "", 0, 0)
	code := `-- @PutHM iost 用户pubkey的base58编码 f10000
@PutHM iost 2BibFrAhc57FAd3sDJFbPqjwskBJb5zPDtecPWVRJ1jxT f100000
@PutHM iost tUFikMypfNGxuJcNbfreh8LM893kAQVNTktVQRsFYuEU f100000
@PutHM iost s1oUQNTcRKL7uqJ1aRqUMzkAkgqJdsBB7uW9xrTd85qB f100000`
	lc := lua.NewContract(vm.ContractInfo{Prefix: "", GasLimit: 0, Price: 0, Publisher: ""}, code, main)

	tx := NewTx(0, &lc)

	genesis := &block.Block{
		Head: block.BlockHead{
			Version: 0,
			Number:  0,
			Time:    initTime,
		},
		Content: make([]Tx, 0),
	}
	genesis.Content = append(genesis.Content, tx)
	stp, err := verifier.ParseGenesis(tx.Contract, p.StatePool())
	if err != nil {
		panic("failed to ParseGenesis")
	}

	err = p.blockCache.SetBasePool(stp)
	if err != nil {
		panic("failed to SetBasePool")
	}

	err = p.blockCache.AddGenesis(genesis)
	if err != nil {
		panic("failed to AddGenesis")
	}
	return nil
}

func (p *DPoS) txListenLoop() {
	p.log.I("Start to listen tx")
	for {
		select {
		case req, ok := <-p.ChTx:
			if !ok {
				return
			}
			if req.ReqType == reqTypeVoteTest {
				p.addWitnessMsg(req)
				continue
			}
			var tx Tx
			tx.Decode(req.Body)
			if VerifyTxSig(tx) {
				p.blockCache.AddTx(&tx)
			}

		case <-p.exitSignal:
			return
		}
	}
}

func (p *DPoS) blockLoop() {

	p.log.I("Start to listen block")
	for {
		select {
		case req, ok := <-p.chBlock:
			if !ok {
				return
			}
			var blk block.Block
			blk.Decode(req.Body)

			/*
				////////////probe//////////////////
				log.Report(&log.MsgBlock{
					SubType:"receive",
					BlockHeadHash:blk.HeadHash(),
					 BlockNum:blk.Head.Number,
				})
				///////////////////////////////////
			*/
			/*
				////////////probe//////////////////
				log.Report(&log.MsgBlock{
					SubType:       "receive",
					BlockHeadHash: blk.HeadHash(),
					BlockNum:      blk.Head.Number,
				})
				///////////////////////////////////
			*/
			p.log.I("Received block:%v , timestamp: %v, Witness: %v, trNum: %v", blk.Head.Number, blk.Head.Time, blk.Head.Witness, len(blk.Content))
			err := p.blockCache.Add(&blk, p.blockVerify)
			if err == nil {
				p.log.I("Link it onto cached chain")
				bc := p.blockCache.LongestChain()
				p.blockCache.UpdateTxPoolOnBC(bc)
			} else {
				p.log.I("Error: %v", err)
			}
			if err != ErrBlock && err != ErrTooOld {
				p.synchronizer.BlockConfirmed(blk.Head.Number)
				if err == nil {
					p.globalDynamicProperty.update(&blk.Head)

					p.blockCache.AddSingles(p.blockVerify)

				} else if err == ErrNotFound {
					// New block is a single block
					need, start, end := p.synchronizer.NeedSync(uint64(blk.Head.Number))
					if need {
						go p.synchronizer.SyncBlocks(start, end)
					}
				}
			}
			/*
								ts := Timestamp{blk.Head.Time}
								if ts.After(p.globalDynamicProperty.NextMaintenanceTime) {
									p.performMaintenance()
				 				}
			*/
		case <-p.exitSignal:
			return
		}
	}
}

func (p *DPoS) scheduleLoop() {
	//通过时间判定是否是本节点的slot，如果是，调用产生块的函数，如果不是，设定一定长的timer睡眠一段时间
	var nextSchedule int64
	p.log.I("Start to schedule")
	for {
		select {
		case <-p.exitSignal:
			return
		case <-time.After(time.Second * time.Duration(nextSchedule)):
			currentTimestamp := GetCurrentTimestamp()
			wid := witnessOfTime(&p.globalStaticProperty, &p.globalDynamicProperty, currentTimestamp)
			p.log.I("currentTimestamp: %v, wid: %v, p.account.ID: %v", currentTimestamp, wid, p.account.ID)
			if wid == p.account.ID {

				//todo test
				bc := p.blockCache.LongestChain()
				iter := bc.Iterator()
				for {
					block := iter.Next()
					if block == nil {
						break
					}
					p.log.I("CBC ConfirmedLength: %v, block Number: %v, witness: %v", p.blockCache.ConfirmedLength(), block.Head.Number, block.Head.Witness)
				}
				// end test

				// TODO 考虑更好的解决方法，因为两次调用之间可能会进入新块影响最长链选择

				pool := p.blockCache.LongestPool()
				blk := p.genBlock(p.account, bc, pool)
				go p.blockCache.ResetTxPoool()
				p.globalDynamicProperty.update(&blk.Head)
				p.log.I("Generating block, current timestamp: %v number: %v", currentTimestamp, blk.Head.Number)

				msg := message.Message{ReqType: int32(ReqNewBlock), Body: blk.Encode()}
				go p.router.Broadcast(msg)
				p.chBlock <- msg
				p.log.I("Broadcasted block, current timestamp: %v number: %v", currentTimestamp, blk.Head.Number)
			}
			nextSchedule = timeUntilNextSchedule(&p.globalStaticProperty, &p.globalDynamicProperty, time.Now().Unix())
			//time.Sleep(time.Second * time.Duration(nextSchedule))
		}
	}
}

func (p *DPoS) genBlock(acc Account, bc block.Chain, pool state.Pool) *block.Block {
	lastBlk := bc.Top()
	blk := block.Block{Content: []Tx{}, Head: block.BlockHead{
		Version:    0,
		ParentHash: lastBlk.Head.Hash(),
		TreeHash:   make([]byte, 0),
		BlockHash:  make([]byte, 0),
		Info:       encodeDPoSInfo(p.infoCache),
		Number:     lastBlk.Head.Number + 1,
		Witness:    acc.ID,
		Time:       GetCurrentTimestamp().Slot,
	}}
	p.infoCache = [][]byte{}
	headInfo := generateHeadInfo(blk.Head)
	sig, _ := common.Sign(common.Secp256k1, headInfo, acc.Seckey)
	blk.Head.Signature = sig.Encode()
	//return &blk
	spool1 := pool.Copy()
	//TODO Content大小控制
	for len(blk.Content) < TxPerBlk {
		tx, err := p.blockCache.GetTx()
		if tx == nil || err != nil {
			break
		}

		if sp, _, err := StdTxsVerifier([]Tx{*tx}, spool1); err == nil {
			blk.Content = append(blk.Content, *tx)
		} else {
			spool1 = sp
		}
	}
	/*
		////////////probe//////////////////
		log.Report(&log.MsgBlock{
			SubType:       "gen",
			BlockHeadHash: blk.HeadHash(),
			BlockNum:      blk.Head.Number,
		})
		///////////////////////////////////
	*/
	return &blk
}

func generateHeadInfo(head block.BlockHead) []byte {
	var info, numberInfo, versionInfo []byte
	info = make([]byte, 8)
	versionInfo = make([]byte, 4)
	numberInfo = make([]byte, 4)
	binary.BigEndian.PutUint64(info, uint64(head.Time))
	binary.BigEndian.PutUint32(versionInfo, uint32(head.Version))
	binary.BigEndian.PutUint32(numberInfo, uint32(head.Number))
	info = append(info, versionInfo...)
	info = append(info, numberInfo...)
	info = append(info, head.ParentHash...)
	info = append(info, head.TreeHash...)
	info = append(info, head.Info...)
	return common.Sha256(info)
}

// 测试函数，用来将info和vote消息进行转换，在块被确认时被调用
// TODO:找到适当的时机调用
func decodeDPoSInfo(info []byte) [][]byte {
	votes := bytes.Split(info, []byte("/"))
	return votes
}

// 测试函数，用来将info和vote消息进行转换，在生成块的时候调用
func encodeDPoSInfo(votes [][]byte) []byte {
	var info []byte
	for _, req := range votes {
		info = append(info, req...)
		info = append(info, byte('/'))
	}
	return info
}

//收到新块，验证新块，如果验证成功，更新DPoS全局动态属性类并将其加入block cache，再广播
func (p *DPoS) blockVerify(blk *block.Block, parent *block.Block, pool state.Pool) (state.Pool, error) {
	/*
		////////////probe//////////////////
		msgBlock := log.MsgBlock{
			SubType:       "verify.fail",
			BlockHeadHash: blk.HeadHash(),
			BlockNum:      blk.Head.Number,
		}
		///////////////////////////////////
	*/
	// verify block head
	if err := VerifyBlockHead(blk, parent); err != nil {
		/*
			////////////probe//////////////////
			log.Report(&msgBlock)
			///////////////////////////////////
		*/
		return nil, err

	}

	// verify block witness
	// TODO currentSlot is negative
	if witnessOfTime(&p.globalStaticProperty, &p.globalDynamicProperty, Timestamp{blk.Head.Time}) != blk.Head.Witness {
		/*
			////////////probe//////////////////
			log.Report(&msgBlock)
			///////////////////////////////////
		*/
		return nil, errors.New("wrong witness")

	}

	headInfo := generateHeadInfo(blk.Head)
	var signature common.Signature
	signature.Decode(blk.Head.Signature)

	// verify block witness signature
	if !common.VerifySignature(headInfo, signature) {
		/*
			////////////probe//////////////////
			log.Report(&msgBlock)
			///////////////////////////////////
		*/
		return nil, errors.New("wrong signature")
	}
	newPool, err := StdBlockVerifier(blk, pool)
	if err != nil {
		/*
			////////////probe//////////////////
			log.Report(&msgBlock)
			///////////////////////////////////
		*/
		return nil, err
	}
	/*
		////////////probe//////////////////
		msgBlock.SubType = "verify.pass"
		log.Report(&msgBlock)
		///////////////////////////////////
	*/
	return newPool, nil
}
