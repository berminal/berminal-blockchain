package consensus_common

import (
	"bytes"
	"errors"

	"github.com/iost-official/prototype/core/block"
	"github.com/iost-official/prototype/core/state"
	"github.com/iost-official/prototype/core/tx"
	"github.com/iost-official/prototype/log"
	"github.com/iost-official/prototype/verifier"
)

//go:generate gencode go -schema=structs.schema -package=consensus_common

// VerifyBlockHead 验证块头正确性
// blk 需要验证的块, parentBlk 父块
func VerifyBlockHead(blk *block.Block, parentBlk *block.Block) error {
	bh := blk.Head
	// parent hash
	if !bytes.Equal(bh.ParentHash, parentBlk.Head.Hash()) {
		return errors.New("wrong parent hash")
	}
	// block number
	if bh.Number != parentBlk.Head.Number+1 {
		return errors.New("wrong number")
	}
	//	treeHash := calcTreeHash(DecodeTxs(blk.Content))
	//	// merkle tree hash
	//	if !bytes.Equal(treeHash, bh.TreeHash) {
	//		return false
	//	}
	return nil
}

var ver *verifier.CacheVerifier

// StdBlockVerifier 块内交易的验证函数
func StdBlockVerifier(block *block.Block, pool state.Pool) (state.Pool, error) {
	if ver == nil {
		veri := verifier.NewCacheVerifier(pool)
		ver = &veri
	}
	pool2 := pool.Copy()
	for _, txx := range block.Content { // TODO
		pool3, err := ver.VerifyContractWithPool(txx.Contract, pool2)
		if err != nil {
			return pool3, err
		}
		pool2 = pool3
	}
	return pool2, nil
}

// VerifyTx 验证单个交易的正确性
// 在调用之前需要先调用vm.NewCacheVerifier(pool state.Pool)生成一个cache verifier
// TODO: 考虑自己生成块到达最后一个交易时，直接用返回的state pool更新block cache中的state
func VerifyTx(tx *tx.Tx, txVer *verifier.CacheVerifier) (state.Pool, bool) {
	newPool, err := txVer.VerifyContract(tx.Contract, false)

	////////////probe//////////////////
	var ret string = "pass"
	if err != nil {
		ret = "fail"
	}
	log.Report(&log.MsgTx{
		SubType:   "verify." + ret,
		TxHash:    tx.Hash(),
		Publisher: tx.Publisher.Pubkey,
		Nonce:     tx.Nonce,
	})
	///////////////////////////////////

	return newPool, err == nil
}

// VerifyTxSig 验证交易的签名
func VerifyTxSig(tx tx.Tx) bool {
	err := tx.VerifySelf()
	return err == nil
}
