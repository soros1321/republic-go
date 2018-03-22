package gateway

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/republic-go/contracts/bindings"
	"github.com/republicprotocol/republic-go/contracts/connection"
)

// Gateway is the interface defining the Republic Protocol's Gateway contract
type Gateway interface {
	Token() (common.Address, error)
	DarkNodeRegistry() (common.Address, error)
	TraderRegistry() (common.Address, error)
	UpdateDarkNodeRegistry(common.Address) (*types.Transaction, error)
	UpdateTraderRegistry(common.Address) (*types.Transaction, error)
}

// EthereumGateway is the Gateway interface
type EthereumGateway struct {
	context        context.Context
	client         *connection.Client
	auth1          *bind.TransactOpts
	auth2          *bind.CallOpts
	binding        *bindings.Gateway
	gatewayAddress common.Address
}

// NewEthereumGateway returns a Gateway
func NewEthereumGateway(context context.Context, clientDetails *connection.ClientDetails, auth1 *bind.TransactOpts, auth2 *bind.CallOpts) (Gateway, error) {
	contract, err := bindings.NewGateway(clientDetails.GatewayAddress, bind.ContractBackend(clientDetails.Client))
	if err != nil {
		return nil, err
	}
	return &EthereumGateway{
		context:        context,
		client:         &clientDetails.Client,
		auth1:          auth1,
		auth2:          auth2,
		binding:        contract,
		gatewayAddress: clientDetails.GatewayAddress,
	}, nil
}

// Token return's the Gateway compatible Republic Token address
func (Gateway *EthereumGateway) Token() (common.Address, error) {
	return Gateway.binding.Token(Gateway.auth2)
}

// DarkNodeRegistry return's the Gateway compatible Dark Node Registry address
func (Gateway *EthereumGateway) DarkNodeRegistry() (common.Address, error) {
	return Gateway.binding.DarkNodeRegistry(Gateway.auth2)
}

// TraderRegistry return's the Gateway compatible Trader Registry address
func (Gateway *EthereumGateway) TraderRegistry() (common.Address, error) {
	return Gateway.binding.TraderRegistry(Gateway.auth2)
}

// UpdateDarkNodeRegistry update's the Dark Node Registry address
func (Gateway *EthereumGateway) UpdateDarkNodeRegistry(newDarkNodeRegistryAddress common.Address) (*types.Transaction, error) {
	return Gateway.binding.UpdateDarkNodeRegistry(Gateway.auth1, newDarkNodeRegistryAddress)
}

// UpdateTraderRegistry update's the Trader Registry address
func (Gateway *EthereumGateway) UpdateTraderRegistry(newTraderRegistryAddress common.Address) (*types.Transaction, error) {
	return Gateway.binding.UpdateTraderRegistry(Gateway.auth1, newTraderRegistryAddress)
}
