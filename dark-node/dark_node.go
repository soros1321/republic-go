package node

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/republic-go/compute"
	"github.com/republicprotocol/republic-go/contracts/connection"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/dark-ocean"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network"
	"github.com/republicprotocol/republic-go/network/dht"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"google.golang.org/grpc"
)

// Prime is the default prime number used to define the finite field.
var Prime, _ = big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)

type DarkNode struct {
	Config

	TestDeltaNotifications chan *compute.Delta

	Logger     *logger.Logger
	ClientPool *rpc.ClientPool
	DHT        *dht.DHT

	DeltaBuilder                      *compute.DeltaBuilder
	DeltaFragmentMatrix               *compute.DeltaFragmentMatrix
	OrderFragmentWorkerQueue          chan *order.Fragment
	OrderFragmentWorker               *OrderFragmentWorker
	DeltaFragmentBroadcastWorkerQueue chan *compute.DeltaFragment
	DeltaFragmentBroadcastWorker      *DeltaFragmentBroadcastWorker
	DeltaFragmentWorkerQueue          chan *compute.DeltaFragment
	DeltaFragmentWorker               *DeltaFragmentWorker
	GossipWorkerQueue                 chan *compute.Delta
	GossipWorker                      *GossipWorker
	FinalizeWorkerQueue               chan *compute.Delta
	FinalizeWorker                    *FinalizeWorker
	ConsensusWorkerQueue              chan *compute.Delta
	ConsensusWorker                   *ConsensusWorker

	Server *grpc.Server
	Swarm  *network.SwarmService
	Dark   *network.DarkService
	Gossip *network.GossipService

	DarkPoolLimit    int64
	DarkPool         *darkocean.DarkPool
	DarkOceanOverlay *darkocean.Overlay
	Registrar        dnr.DarkNodeRegistrar
	EpochBlockhash   [32]byte
}

// NewDarkNode return a DarkNode that adheres to the given Config. The DarkNode
// will configure all of the components that it needs to operate but will not
// start any of them.
func NewDarkNode(config Config, darkNodeRegistrar dnr.DarkNodeRegistrar) (*DarkNode, error) {
	if config.Prime == nil {
		config.Prime = Prime
	}

	// TODO: This should come from the DNR.
	k := int64(5)

	node := &DarkNode{Config: config, TestDeltaNotifications: make(chan *compute.Delta, 100)}

	logger, err := logger.NewLogger(config.LoggerOptions)
	if err != nil {
		return nil, err
	}
	node.Logger = logger
	node.Logger.Start()

	node.ClientPool = rpc.NewClientPool(node.NetworkOptions.MultiAddress)
	node.DHT = dht.NewDHT(node.NetworkOptions.MultiAddress.Address(), node.NetworkOptions.MaxBucketLength)

	node.DeltaBuilder = compute.NewDeltaBuilder(k, node.Prime)
	node.DeltaFragmentMatrix = compute.NewDeltaFragmentMatrix(node.Prime)
	node.OrderFragmentWorkerQueue = make(chan *order.Fragment, 100)
	node.OrderFragmentWorker = NewOrderFragmentWorker(node.Logger, node.DeltaFragmentMatrix, node.OrderFragmentWorkerQueue)
	node.DeltaFragmentBroadcastWorkerQueue = make(chan *compute.DeltaFragment, 100)
	node.DeltaFragmentBroadcastWorker = NewDeltaFragmentBroadcastWorker(node.Logger, node.ClientPool, node.DarkPool, node.DeltaFragmentBroadcastWorkerQueue)
	node.DeltaFragmentWorkerQueue = make(chan *compute.DeltaFragment, 100)
	node.DeltaFragmentWorker = NewDeltaFragmentWorker(node.Logger, node.DeltaBuilder, node.DeltaFragmentWorkerQueue)
	node.GossipWorkerQueue = make(chan *compute.Delta, 100)
	node.GossipWorker = NewGossipWorker(node.Logger, node.ClientPool, node.NetworkOptions.BootstrapMultiAddresses, node.GossipWorkerQueue)
	node.FinalizeWorkerQueue = make(chan *compute.Delta, 100)
	node.FinalizeWorker = NewFinalizeWorker(node.Logger, k, node.FinalizeWorkerQueue)
	node.ConsensusWorkerQueue = make(chan *compute.Delta, 100)
	node.ConsensusWorker = NewConsensusWorker(node.Logger, node.DeltaFragmentMatrix, node.ConsensusWorkerQueue)

	// options := network.Options{}
	node.Server = grpc.NewServer(grpc.ConnectionTimeout(time.Minute))
	node.Swarm = network.NewSwarmService(node, node.NetworkOptions, node.Logger, node.ClientPool, node.DHT)
	node.Dark = network.NewDarkService(node, node.NetworkOptions, node.Logger)
	node.Gossip = network.NewGossipService(node)

	clientDetails, err := connection.FromURI(node.EthereumRPC)
	if err != nil {
		return nil, err
	}
	registrar, err := node.ConnectToRegistrar(clientDetails, config)
	if err != nil {
		return nil, err
	}
	node.Registrar = registrar

	return node, nil
}

// Start the DarkNode.
func (node *DarkNode) Start() {
	// Begin broadcasting CPU/Memory/Network usage
	go func() {
		for {
			node.Usage()
			time.Sleep(20 * time.Second)
		}
	}()
	go node.ServeUI()
	// Wait until the node is registered
	//for isRegistered := node.IsRegistered(); !isRegistered; isRegistered = node.IsRegistered() {
	//	timeout := 60 * time.Second
	//	node.Warn(logger.TagNetwork, fmt.Sprintf("%v not registered. Sleeping for %v seconds.", node.MultiAddress.Address(), timeout.Seconds()))
	//
	//	data := logger.Registration{
	//		NodeID:     "0x" + hex.EncodeToString(node.MultiAddress.ID()),
	//		PublicKey:  "0x" + hex.EncodeToString(append(node.Config.RepublicKeyPair.PublicKey.X.Bytes(), node.Config.RepublicKeyPair.PublicKey.Y.Bytes()...)),
	//		Address:    node.Config.EthereumKey.Address.String(),
	//		RepublicID: node.MultiAddress.ID().String(),
	//	}
	//	dataJson, err := json.Marshal(data)
	//	if err != nil {
	//		node.Error(logger.TagGeneral, err.Error())
	//	}
	//	// Send the info needed for registration as well
	//	err = node.Logger.Info(logger.TagRegister, string(dataJson))
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	time.Sleep(timeout)
	//}
	node.Logger.Info(logger.TagEthereum, "Successfully registered")

	// Start serving the gRPC services
	go func() {
		node.Logger.Info(logger.TagNetwork, fmt.Sprintf("Listening on %s:%s", node.Host, node.Port))
		node.Swarm.Register(node.Server)
		node.Dark.Register(node.Server)
		node.Gossip.Register(node.Server)
		listener, err := net.Listen("tcp", node.Host+":"+node.Port)
		if err != nil {
			node.Logger.Error(logger.TagNetwork, err.Error())
		}
		if err := node.Server.Serve(listener); err != nil {
			node.Logger.Error(logger.TagNetwork, err.Error())
		}
	}()
	time.Sleep(time.Second)

	// Bootstrap into the swarm network
	go node.Swarm.Bootstrap()
	time.Sleep(time.Second)

	// Run the workers
	go node.OrderFragmentWorker.Run(node.DeltaFragmentBroadcastWorkerQueue, node.DeltaFragmentWorkerQueue)
	go node.DeltaFragmentBroadcastWorker.Run()
	go node.DeltaFragmentWorker.Run(node.GossipWorkerQueue, node.TestDeltaNotifications)
	go node.GossipWorker.Run(node.FinalizeWorkerQueue)
	go node.FinalizeWorker.Run(node.ConsensusWorkerQueue)
	go node.ConsensusWorker.Run()

	// oceanChanges := make(chan do.Option)
	// defer close(oceanChanges)
	// go darkocean.WatchForDarkOceanChanges(node.Registrar, oceanChanges)

	// for {
	// 	select {
	// 	case ocean := <-oceanChanges:
	// 		if ocean.Err != nil {
	// 			node.Logger.Error(logger.TagEthereum, ocean.Err.Error())
	// 			continue
	// 		}
	// 		node.AfterEachEpoch()
	// 	}
	// }
	q := make(chan struct{})
	<-q
}

func (node *DarkNode) ServeUI() {
	fs := http.FileServer(http.Dir("darknode-ui"))
	http.Handle("/", fs)
	node.Logger.Info(logger.TagNetwork, "Serving the Dark Node UI")
	err := http.ListenAndServe("0.0.0.0:3000", nil)
	if err != nil {
		node.Logger.Error(logger.TagNetwork, err.Error())
	}
}

// Stop the DarkNode.
func (node *DarkNode) Stop() {
	// Stop serving gRPC services
	node.Server.Stop()
	time.Sleep(time.Second)

	// Stop the workers
	close(node.OrderFragmentWorkerQueue)
	close(node.DeltaFragmentWorkerQueue)
}

// OnOpenOrder writes an order fragment that has been received to the
// OrderFragmentWorkerQueue. This is a potentially blocking operation, however
// this delegate method is called on a dedicated goroutine.
func (node *DarkNode) OnOpenOrder(from identity.MultiAddress, orderFragment *order.Fragment) {
	// Write to a channel that might be closed
	func() {
		defer func() { recover() }()
		node.OrderFragmentWorkerQueue <- orderFragment
	}()
}

// OnBroadcastDeltaFragment writes a delta fragment that has been received to
// the DeltaFragmentWorkerQueue. This is a potentially blocking operation,
// however this delegate method is called on a dedicated goroutine.
func (node *DarkNode) OnBroadcastDeltaFragment(from identity.MultiAddress, deltaFragment *compute.DeltaFragment) {
	// Write to a channel that might be closed
	func() {
		defer func() { recover() }()
		node.DeltaFragmentWorkerQueue <- deltaFragment
	}()
}

func (node *DarkNode) OnGossip(buyOrderID order.ID, sellOrderID order.ID) {
	// Write to a channel that might be closed
	func() {
		defer func() { recover() }()
		node.GossipWorkerQueue <- &compute.Delta{
			ID:          compute.DeltaID(crypto.Keccak256([]byte(buyOrderID), []byte(sellOrderID))),
			BuyOrderID:  buyOrderID,
			SellOrderID: sellOrderID,
		}
	}()
}

func (node *DarkNode) OnFinalize(buyOrderID order.ID, sellOrderID order.ID) {
	// Write to a channel that might be closed
	func() {
		defer func() { recover() }()
		node.FinalizeWorkerQueue <- &compute.Delta{
			ID:          compute.DeltaID(crypto.Keccak256([]byte(buyOrderID), []byte(sellOrderID))),
			BuyOrderID:  buyOrderID,
			SellOrderID: sellOrderID,
		}
	}()
}

// IsRegistered returns true if the dark node is registered for the current epoch
func (node *DarkNode) IsRegistered() bool {
	registered, err := node.Registrar.IsDarkNodeRegistered(node.NetworkOptions.MultiAddress.ID())
	if err != nil {
		return false
	}
	return registered
}

// IsPendingRegistration returns true if the dark node will be registered in the next epoch
func (node *DarkNode) IsPendingRegistration() bool {
	registered, err := node.Registrar.IsDarkNodePendingRegistration(node.NetworkOptions.MultiAddress.ID())
	if err != nil {
		return false
	}
	return registered
}

// Register the node on the registrar smart contract .
func (node *DarkNode) Register() error {
	registered := node.IsRegistered()
	if registered {
		return nil
	}
	publicKey := append(node.Config.RepublicKeyPair.PublicKey.X.Bytes(), node.Config.RepublicKeyPair.PublicKey.Y.Bytes()...)
	_, err := node.Registrar.Register(node.NetworkOptions.MultiAddress.ID(), publicKey)
	if err != nil {
		return err
	}
	err = node.Registrar.WaitTillRegistration(node.NetworkOptions.MultiAddress.ID())
	return err
}

// Deregister the node on the registrar smart contract
func (node *DarkNode) Deregister() error {
	_, err := node.Registrar.Deregister(node.NetworkOptions.MultiAddress.ID())
	if err != nil {
		return err
	}
	return nil
}

// PingDarkPool call rpc.PingTarget on each node in a dark pool
func (node *DarkNode) PingDarkPool(ids darkocean.DarkPoolID) (identity.MultiAddresses, darkocean.DarkPoolID) {

	darkpool := make(identity.MultiAddresses, 0)
	disconnectedDarkPool := make(darkocean.DarkPoolID, 0)

	for _, id := range ids {
		target, err := node.Swarm.FindNode(id)
		if err != nil || target == nil {
			node.Logger.Warn(logger.TagNetwork, fmt.Sprintf("%v couldn't find pool peer %v: %v", node.NetworkOptions.MultiAddress.Address(), id, err))
			disconnectedDarkPool = append(disconnectedDarkPool, id)
			continue
		}

		darkpool = append(darkpool, *target)

		node.ClientPool.Ping(*target)
		if err != nil {
			node.Logger.Warn(logger.TagNetwork, fmt.Sprintf("%v couldn't ping pool peer %v: %v", node.NetworkOptions.MultiAddress.Address(), target, err))
			continue
		}

		err = node.Swarm.DHT.UpdateMultiAddress(*target)
		if err != nil {
			node.Logger.Warn(logger.TagNetwork, fmt.Sprintf("%v coudln't update DHT for pool peer %v: %v", node.NetworkOptions.MultiAddress.Address(), target, err))
			continue
		}
	}
	return darkpool, disconnectedDarkPool
}

// LongPingDarkPool will continually attempt to connect to a set of nodes
// in a darkpool until they are all connected
// Call in a goroutine
func (node *DarkNode) LongPingDarkPool(disconnected darkocean.DarkPoolID) {
	currentBlockhash := node.EpochBlockhash

	for len(disconnected) > 0 {
		if node.EpochBlockhash != currentBlockhash {
			return
		}

		var connected identity.MultiAddresses
		connected, disconnected = node.PingDarkPool(disconnected)

		node.DarkPool.Add(connected...)

		time.Sleep(30 * time.Second)
	}
}

// AfterEachEpoch should be run after each new epoch
func (node *DarkNode) AfterEachEpoch() error {
	node.Logger.Info(logger.TagNetwork, fmt.Sprintf("%v is pinging dark pool\n", node.NetworkOptions.MultiAddress.Address()))

	darkOceanOverlay, err := darkocean.GetDarkPools(node.Registrar)
	if err != nil {
		node.Logger.Error(logger.TagNetwork, fmt.Sprintf("%v couldn't get dark pools: %v", node.NetworkOptions.MultiAddress.Address(), err))
		return err
	}
	node.DarkOceanOverlay = darkOceanOverlay

	poolID := node.DarkOceanOverlay.FindDarkPool(node.NetworkOptions.MultiAddress.ID())
	if poolID == nil {
		return fmt.Errorf("cannot find %s in the dark ocean", node.NetworkOptions.MultiAddress.Address())
	}

	connectedDarkPool, disconnectedDarkPool := node.PingDarkPool(poolID)
	node.DarkPool = darkocean.NewDarkPool(connectedDarkPool)

	node.Logger.Info(logger.TagNetwork, fmt.Sprintf("%v connected to dark pool: %v", node.NetworkOptions.MultiAddress.Address(), node.DarkPool))

	go node.LongPingDarkPool(disconnectedDarkPool)

	return nil
}

// ConnectToRegistrar will connect to the registrar using the given private key to sign transactions
func (node DarkNode) ConnectToRegistrar(clientDetails connection.ClientDetails, config Config) (dnr.DarkNodeRegistrarInterface, error) {
	auth := bind.NewKeyedTransactor(node.Config.EthereumKey.PrivateKey)

	// Gas Price
	auth.GasPrice = big.NewInt(6000000000)

	userConnection := dnr.NewDarkNodeRegistrar(context.Background(), &clientDetails, auth, &bind.CallOpts{})
	return userConnection, nil
}

// Usage logs memory and cpu usage
func (node *DarkNode) Usage() {
	// memory
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		node.Logger.Error(logger.TagUsage, err.Error())
	}
	// cpu - get CPU number of cores and speed
	cpuStat, err := cpu.Info()
	if err != nil {
		node.Logger.Error(logger.TagUsage, err.Error())
	}
	percentage, err := cpu.Percent(0, false)
	if err != nil {
		node.Logger.Error(logger.TagUsage, err.Error())
	}

	node.Logger.Usage(float32(cpuStat[0].Mhz*percentage[0]/100), int32(vmStat.Used/1024/1024), 0)
}
