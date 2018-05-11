package consensus_common

import (
	"github.com/iost-official/prototype/core/message"
	. "github.com/iost-official/prototype/network"
	"time"
)

var (
	SyncNumber        = 10
	MaxDownloadNumber = 10
)

type Synchronizer interface {
	StartListen() error
	NeedSync(maxHeight uint64) (bool, uint64, uint64)
	SyncBlocks(startNumber uint64, endNumber uint64) error
}

type SyncImpl struct {
	blockCache  BlockCache
	router      Router
	heightChan  chan message.Message
	blkSyncChan chan message.Message
}

func NewSynchronizer(bc BlockCache, router Router) *SyncImpl {
	sync := &SyncImpl{
		blockCache: bc,
		router:     router,
	}
	var err error
	sync.heightChan, err = sync.router.FilteredChan(Filter{
		WhiteList:  []message.Message{},
		BlackList:  []message.Message{},
		RejectType: []ReqType{},
		AcceptType: []ReqType{
			ReqBlockHeight,
		}})
	if err != nil {
		return nil
	}

	sync.blkSyncChan, err = sync.router.FilteredChan(Filter{
		WhiteList:  []message.Message{},
		BlackList:  []message.Message{},
		RejectType: []ReqType{},
		AcceptType: []ReqType{
			ReqDownloadBlock,
		}})
	if err != nil {
		return nil
	}
	return sync
}

//开始监听同步任务
func (sync *SyncImpl) StartListen() error {
	go sync.requestBlockHeightLoop()
	go sync.requestBlockLoop()
	go sync.blockConfirmLoop()

	return nil
}

func (sync *SyncImpl) NeedSync(maxHeight uint64) (bool, uint64, uint64) {
	height := sync.blockCache.ConfirmedLength()
	if height < maxHeight-uint64(SyncNumber) {
		body := message.RequestHeight{
			LocalBlockHeight: height,
			NeedBlockHeight:  maxHeight,
		}
		heightReq := message.Message{
			ReqType: int32(ReqBlockHeight),
			Body:    body.Encode(),
		}
		sync.router.Broadcast(heightReq)
		return true, height + 1, maxHeight
	}
	return false, 0, 0
}

func (sync *SyncImpl) SyncBlocks(startNumber uint64, endNumber uint64) error {
	for endNumber-startNumber > uint64(MaxDownloadNumber) {
		sync.router.Download(startNumber, startNumber+uint64(MaxDownloadNumber))
		//TODO 等待所有区间里的块都收到
		startNumber += uint64(MaxDownloadNumber + 1)
	}
	sync.router.Download(startNumber, endNumber)
	return nil
}

func (sync *SyncImpl) requestBlockHeightLoop() {

	for {
		req, ok := <-sync.heightChan
		if !ok {
			return
		}
		var rh message.RequestHeight
		rh.Decode(req.Body)

		localLength := sync.blockCache.LongestChain().Length()

		//本地链长度小于等于远端，忽略远端的同步链请求
		if localLength <= rh.LocalBlockHeight {
			continue
		}
		//回复当前块的高度
		hr := message.ResponseHeight{BlockHeight: localLength}
		resMsg := message.Message{
			Time:    time.Now().Unix(),
			From:    req.To,
			To:      req.From,
			ReqType: int32(RecvBlockHeight),
			Body:    hr.Encode(),
		}

		sync.router.Send(resMsg)
	}
}

func (sync *SyncImpl) requestBlockLoop() {

	for {
		req, ok := <-sync.blkSyncChan
		if !ok {
			return
		}
		var rh message.RequestBlock
		rh.Decode(req.Body)

		chain := sync.blockCache.LongestChain()

		//todo 需要确定如何获取
		block := chain.GetBlockByNumber(rh.BlockNumber)
		if block == nil {
			continue
		}

		//回复当前块的高度
		resMsg := message.Message{
			Time:    time.Now().Unix(),
			From:    req.To,
			To:      req.From,
			ReqType: int32(ReqNewBlock), //todo 后补类型
			Body:    block.Encode(),
		}

		sync.router.Send(resMsg)
	}
}

func (sync *SyncImpl) blockConfirmLoop() {
	for {
		num, ok := <-sync.blockCache.BlockConfirmChan()
		if !ok {
			return
		}
		sync.router.CancelDownload(num, num)
	}
}
