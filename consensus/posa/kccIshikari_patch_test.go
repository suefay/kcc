package posa

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
)

func TestVerifyCode(t *testing.T) {

	internalDB := rawdb.NewMemoryDatabase()
	triedb := state.NewDatabaseWithConfig(internalDB, nil)
	stateDB, err := state.New(common.Hash{}, triedb, nil)

	if err != nil {
		t.Fatalf("failed to create statedb: %v", err)
	}

	patch001PunishContract(stateDB)
	patch001ValidatorsContract(stateDB)

}
