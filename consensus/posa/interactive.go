package posa

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

type chainContext struct {
	chainReader consensus.ChainHeaderReader
	engine      consensus.Engine
}

func newChainContext(chainReader consensus.ChainHeaderReader, engine consensus.Engine) *chainContext {
	return &chainContext{
		chainReader: chainReader,
		engine:      engine,
	}
}

// Engine retrieves the chain's consensus engine.
func (cc *chainContext) Engine() consensus.Engine {
	return cc.engine
}

// GetHeader returns the hash corresponding to their hash.
func (cc *chainContext) GetHeader(hash common.Hash, number uint64) *types.Header {
	return cc.chainReader.GetHeader(hash, number)
}

func getInteractiveABIAndAddrs() (map[string]abi.ABI, map[string]common.Address) {
	abiMap := make(map[string]abi.ABI, 0)
	tmpABI, _ := abi.JSON(strings.NewReader(validatorsInteractiveABI))
	abiMap[validatorsContractName] = tmpABI
	tmpABI, _ = abi.JSON(strings.NewReader(punishInteractiveABI))
	abiMap[punishContractName] = tmpABI
	tmpABI, _ = abi.JSON(strings.NewReader(proposalInteractiveABI))
	abiMap[proposalContractName] = tmpABI

	// Add our new abi encoders
	abiMap[IshikariProposalContractName], _ = abi.JSON(strings.NewReader(IshikariProposalABI))
	abiMap[IshikariPunishContractName], _ = abi.JSON(strings.NewReader(IshikariPunishABI))
	abiMap[IshikariReservePoolContractName], _ = abi.JSON(strings.NewReader(IshikariReservePoolABI))
	abiMap[IshikariValidatorsContractName], _ = abi.JSON(strings.NewReader(IshikariValidatorABI))

	// Contract Addresses
	addrs := make(map[string]common.Address, 0)

	// v1 addresses
	addrs[validatorsContractName] = validatorsContractAddr
	addrs[proposalContractName] = proposalAddr
	addrs[punishContractName] = punishContractAddr

	// Ishikari hardfork addresses
	addrs[IshikariValidatorsContractName] = IshikariValidatorsContractAddr
	addrs[IshikariProposalContractName] = IshikariProposalAddr
	addrs[IshikariPunishContractName] = IshikariPunishContractAddr
	addrs[IshikariReservePoolContractName] = IshikariReservePoolAddr

	return abiMap, addrs // TODO
}

// executeMsg executes transaction sent to system contracts.
func executeMsg(msg core.Message, state *state.StateDB, header *types.Header, chainContext core.ChainContext, chainConfig *params.ChainConfig) (ret []byte, err error) {
	// Set gas price to zero
	context := core.NewEVMBlockContext(header, chainContext, nil)
	txContext := core.NewEVMTxContext(msg)
	vmenv := vm.NewEVM(context, txContext, state, chainConfig, vm.Config{})

	msg.GasPrice()

	ret, _, err = vmenv.Call(vm.AccountRef(msg.From()), *msg.To(), msg.Data(), msg.Gas(), msg.Value())

	if err != nil {
		return ret, err
	}

	return ret, nil
}
