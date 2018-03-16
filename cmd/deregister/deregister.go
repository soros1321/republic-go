package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	base58 "github.com/jbenet/go-base58"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/contracts/connection"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	node "github.com/republicprotocol/republic-go/dark-node"
)

func main() {

	configFiles := os.Args[1:]
	configs := make([]*node.Config, len(configFiles))

	for file := range configFiles {
		fileName := configFiles[file]
		config, err := node.LoadConfig(fileName)
		if err != nil {
			panic(err)
		}
		configs[file] = config
	}

	DeregisterAll(configs)
}

// DeregisterAll takes a slice of republic private keys and deregisters them
func DeregisterAll(configs []*node.Config) {
	const reset = "\x1b[0m"
	const yellow = "\x1b[33;1m"
	const green = "\x1b[32;1m"
	const red = "\x1b[31;1m"

	do.ForAll(configs, func(i int) {

		keypair := configs[i].RepublicKeyPair

		clientDetails, err := connection.FromURI("https://ropsten.infura.io/", "ropsten")
		if err != nil {
			// TODO: Handler err
			panic(err)
		}

		key := `{"version":3,"id":"7844982f-abe7-4690-8c15-34f75f847c66","address":"db205ea9d35d8c01652263d58351af75cfbcbf07","Crypto":{"ciphertext":"378dce3c1279b36b071e1c7e2540ac1271581bff0bbe36b94f919cb73c491d3a","cipherparams":{"iv":"2eb92da55cc2aa62b7ffddba891f5d35"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"80d3341678f83a14024ba9c3edab072e6bd2eea6aa0fbc9e0a33bae27ffa3d6d","n":8192,"r":8,"p":1},"mac":"3d07502ea6cd6b96a508138d8b8cd2e46c3966240ff276ce288059ba4235cb0d"}}`
		auth, err := bind.NewTransactor(strings.NewReader(key), "password1")
		if err != nil {
			panic(err)
		}

		registrar, err := dnr.NewEthereumDarkNodeRegistrar(context.Background(), &clientDetails, auth, &bind.CallOpts{})
		if err != nil {
			log.Printf("[%v] %sCouldn't connect to registrar%s: %v\n", base58.Encode(keypair.ID()), red, reset, err)
			return
		}

		isRegistered, err := registrar.IsDarkNodeRegistered(keypair.ID())
		if err != nil {
			log.Printf("[%v] %sCouldn't check node's registration%s: %v\n", base58.Encode(keypair.ID()), red, reset, err)
			return
		}

		isPendingRegistration, err := registrar.IsDarkNodePendingRegistration(keypair.ID())
		if err != nil {
			log.Printf("[%v] %sCouldn't check node's registration%s: %v\n", base58.Encode(keypair.ID()), red, reset, err)
			return
		}

		if isRegistered || isPendingRegistration {
			_, err = registrar.Deregister(keypair.ID())
			if err != nil {
				log.Printf("[%v] %sCouldn't deregister node%s: %v\n", base58.Encode(keypair.ID()), yellow, reset, err)
				return
			}
			if isRegistered {
				log.Printf("[%v] %sNode will be deregistered next epoch%s\n", base58.Encode(keypair.ID()), green, reset)
			} else {
				log.Printf("[%v] %sNode deregistered (registration cancelled)%s\n", base58.Encode(keypair.ID()), green, reset)
				registrar.Refund(keypair.ID())
			}
		} else {
			log.Printf("[%v] %sNode already deregistered%s\n", base58.Encode(keypair.ID()), yellow, reset)
		}
	})

}
