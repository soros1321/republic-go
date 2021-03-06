package dark

import (
	"bytes"
	"fmt"
	"math/big"
	"time"

	"github.com/republicprotocol/go-do"

	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
)

// Ocean of Pools.
type Ocean struct {
	do.GuardedObject

	logger            *logger.Logger
	pools             Pools
	darkNodeRegistrar dnr.DarkNodeRegistrar
}

// NewOcean uses a DarkNodeRegistrar to read all registered nodes and sort them
// into Pools.
func NewOcean(logger *logger.Logger, darkNodeRegistrar dnr.DarkNodeRegistrar) (*Ocean, error) {
	ocean := &Ocean{
		GuardedObject:     do.NewGuardedObject(),
		logger:            logger,
		pools:             Pools{},
		darkNodeRegistrar: darkNodeRegistrar,
	}
	return ocean, ocean.update()
}

// FindPool with the given node ID. Returns the Pool, or nil if no Pool can be
// found.
func (ocean *Ocean) FindPool(id identity.ID) *Pool {
	ocean.EnterReadOnly(nil)
	defer ocean.ExitReadOnly()
	for _, pool := range ocean.pools {
		if pool.Has(id) != nil {
			return pool
		}
	}
	return nil
}

func (ocean *Ocean) Update() error {
	ocean.Enter(nil)
	defer ocean.Exit()
	return ocean.update()
}

func (ocean *Ocean) update() error {
	// TODO: Get these details from the smart contract.
	blockhash := big.NewInt(1234567)
	poolsize := 72

	nodeIDs, err := ocean.darkNodeRegistrar.GetAllNodes()
	if err != nil {
		return err
	}

	// Find the prime smaller or equal to the number of registered nodes
	// Start at +2 because it has to greater than the maximum (x+1)
	previousPrime := big.NewInt(int64(len(nodeIDs) + 2))

	// ProbablyPrime is 100% accurate for inputs less than 2^64.
	// https://golang.org/src/math/big/prime.go
	for !previousPrime.ProbablyPrime(0) {
		previousPrime = previousPrime.Sub(previousPrime, big.NewInt(1))
	}

	// Integer division
	numberOfPools := big.NewInt(0).Div(previousPrime, big.NewInt(int64(poolsize)))
	if numberOfPools.Int64() == 0 {
		numberOfPools = big.NewInt(1)
	}
	pools := make(Pools, numberOfPools.Int64())
	for i := range pools {
		pools[i] = NewPool()
	}

	// Calcualte the pool assignment for each node
	inverse := blockhash.ModInverse(blockhash, previousPrime)
	for n := range nodeIDs {
		nPlusOne := big.NewInt(int64(n + 1))
		i := big.NewInt(0).Mod(big.NewInt(0).Mul(nPlusOne, inverse), previousPrime)

		pool := big.NewInt(0).Mod(i, numberOfPools).Int64()
		pools[pool].Append(NewNode(nodeIDs[n]))
	}

	ocean.pools = pools
	return nil
}

// Watch for changes to the Ocean. This function is a blocking function that
// sleeps and wakes once per period to check for a change in epoch. It accepts
// a channel that is pinged whenever the Ocean changes.
func (ocean *Ocean) Watch(period time.Duration, changes chan struct{}) {
	// Recover from writing to a closed channel
	defer func() { recover() }()

	var currentBlockhash [32]byte
	for {
		epoch, err := ocean.darkNodeRegistrar.CurrentEpoch()
		if err != nil {
			ocean.logger.Error(fmt.Sprintf("cannot update epoch: %s", err.Error()))
			return
		}
		if !bytes.Equal(currentBlockhash[:], epoch.Blockhash[:]) {
			currentBlockhash = epoch.Blockhash
			if err := ocean.Update(); err != nil {
				ocean.logger.Error(fmt.Sprintf("cannot update dark ocean: %s", err.Error()))
				return
			}
			changes <- struct{}{}
		}
		time.Sleep(period)
	}
}
