package block

import (
	"github.com/iost-official/prototype/common"
	"github.com/iost-official/prototype/core/tx"
	"github.com/iost-official/prototype/vm"
)

//go:generate gencode go -schema=structs.schema -package=block

func (d *Block) Encode() []byte {
	bin, err := d.Marshal(nil)
	if err != nil {
		panic(err)
	}
	return bin
}

func (d *Block) Decode(bin []byte) error {
	_, err := d.Unmarshal(bin)
	return err
}
func (d *Block) Hash() []byte {
	return common.Sha256(d.Encode())
}

func (d *Block) HeadHash() []byte {
	return d.Head.Hash()
}

func (d *Block) TxGet(x int) tx.Tx {
	return tx.Tx{}
}

func (d *Block) TxLen() int {
	return 0
}

func (d *Block)GetAllContract() []vm.Contract {
	//todo 解析content,获得所有交易

	var allContract []vm.Contract
	//todo 将交易中的contract，添加到contractAll中

	return allContract
}

func (d *BlockHead) Encode() []byte {
	bin, err := d.Marshal(nil)
	if err != nil {
		panic(err)
	}
	return bin
}

func (d *BlockHead) Decode(bin []byte) error {
	_, err := d.Unmarshal(bin)
	return err
}
func (d *BlockHead) Hash() []byte {
	return common.Sha256(d.Encode())
}


