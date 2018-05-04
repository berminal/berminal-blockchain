package verifier

import (
	"fmt"
	"github.com/iost-official/prototype/core/state"
	"github.com/iost-official/prototype/vm"
	"github.com/iost-official/prototype/vm/lua"
	"runtime"
)

const (
	MaxBlockGas uint64 = 1000000
)

//go:generate gencode go -schema=structs.schema -package=verifier

type Verifier struct {
	Pool   state.Pool
	Prefix string
	Vms    map[string]vm.VM
}

func (v *Verifier) StartVm(contract vm.Contract) {
	switch contract.(type) {
	case *lua.Contract:
		var lvm lua.VM
		lvm.Prepare(contract.(*lua.Contract), v.Pool)
		lvm.Start()
		v.Vms[string(contract.Hash())] = &lvm
	}
}
func (v *Verifier) StopVm(contract vm.Contract) {
	v.Vms[string(contract.Hash())].Stop()
	delete(v.Vms, string(contract.Hash()))
}
func (v *Verifier) Stop() {
	for _, vm := range v.Vms {
		vm.Stop()
	}
}

func (v *Verifier) Verify(contract vm.Contract) (state.Pool, uint64, error) {

	vm, ok := v.Vms[string(contract.Hash())]
	if !ok {
		return nil, 0, fmt.Errorf("not prepared")
	}
	_, pool, err := vm.Call("main")

	return pool, vm.PC(), err
}

func (v *Verifier) SetPool(pool state.Pool) {
	v.Pool = pool
	for _, vm := range v.Vms {
		vm.SetPool(pool)
	}
}

type CacheVerifier struct {
	Verifier
}

func (cv *CacheVerifier) VerifyContract(contract vm.Contract, contain bool) (state.Pool, error) {
	sender := contract.Info().Sender
	var balanceOfSender float64
	val0, err := cv.Pool.GetHM("iost", state.Key(sender))
	val, ok := val0.(*state.VFloat)
	if !ok {
		return nil, fmt.Errorf("type error")
	}
	balanceOfSender = val.ToFloat64()

	if balanceOfSender < 1 {
		return nil, fmt.Errorf("balance not enough")
	}

	cv.StartVm(contract)
	pool, gas, err := cv.Verify(contract)
	if err != nil {
		return nil, err
	}
	cv.StopVm(contract)

	if gas > uint64(contract.Info().GasLimit) {
		return nil, fmt.Errorf("gas exceed")
	}

	balanceOfSender -= float64(gas) * contract.Info().Price
	if balanceOfSender < 0 {
		balanceOfSender = 0
		return nil, fmt.Errorf("can not afford gas")
	}
	val1 := state.MakeVFloat(balanceOfSender)
	cv.Pool.PutHM("iost", state.Key(sender), &val1)

	if contain {
		cv.SetPool(pool)
	}

	return cv.Pool, nil
}

func NewCacheVerifier(pool state.Pool) CacheVerifier {
	cv := CacheVerifier{
		Verifier: Verifier{
			Pool:   pool.Copy(),
			Prefix: "cache+",
			Vms:    make(map[string]vm.VM),
		},
	}
	runtime.SetFinalizer(cv, func() {
		cv.Verifier.Stop()
	})
	return cv
}

func VerifyBlock(blockID string, contracts []vm.Contract, pool state.Pool) (state.Pool, error) { // TODO 使用log控制并发
	cv := Verifier{
		Pool:   pool,
		Prefix: blockID,
		Vms:    make(map[string]vm.VM),
	}
	var totalGas uint64
	for _, c := range contracts {

		sender := c.Info().Sender
		var balanceOfSender float64
		val0, err := cv.Pool.GetHM("iost", state.Key(sender))
		val, ok := val0.(*state.VFloat)
		if !ok {
			return nil, fmt.Errorf("type error")
		}
		balanceOfSender = val.ToFloat64()

		if balanceOfSender < 1 {
			return nil, fmt.Errorf("balance not enough")
		}

		cv.StartVm(c)
		_, gas, err := cv.Verify(c)
		if err != nil {
			return nil, err
		}
		if gas > uint64(c.Info().GasLimit) {
			return nil, fmt.Errorf("gas exceed")
		}
		cv.StopVm(c)
		totalGas += gas
		if totalGas > MaxBlockGas {
			return nil, fmt.Errorf("block gas exceed")
		}
		balanceOfSender -= float64(gas) * c.Info().Price
		if balanceOfSender < 0 {
			balanceOfSender = 0
			return nil, fmt.Errorf("can not afford gas")
		}
		val1 := state.MakeVFloat(balanceOfSender)
		cv.Pool.PutHM("iost", state.Key(sender), &val1)

	}
	return cv.Pool, nil
}