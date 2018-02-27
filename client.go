package rpc

import (
	"fmt"
	"log"
	"time"

	identity "github.com/republicprotocol/go-identity"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Client struct {
	*grpc.ClientConn

	DarkNode DarkNodeClient
	Options  ClientOptions
	From     *MultiAddress
}

type ClientOptions struct {
	Timeout        time.Duration
	TimeoutBackoff time.Duration
	TimeoutRetries int
}

func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		Timeout:        30 * time.Second,
		TimeoutBackoff: 0 * time.Second,
		TimeoutRetries: 3,
	}
}

func NewClient(to, from identity.MultiAddress) (Client, error) {
	client := Client{
		Options: DefaultClientOptions(),
		From:    SerializeMultiAddress(from),
	}

	host, err := to.ValueForProtocol(identity.IP4Code)
	if err != nil {
		return client, err
	}
	port, err := to.ValueForProtocol(identity.TCPCode)
	if err != nil {
		return client, err
	}

	for i := 0; i < client.Options.TimeoutRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), client.Options.Timeout+(client.Options.TimeoutBackoff*time.Duration(i)))
		defer cancel()
		client.ClientConn, err = grpc.DialContext(ctx, fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
		if err == nil {
			break
		}
		log.Println(err)
	}
	if err != nil {
		return client, err
	}
	client.DarkNode = NewDarkNodeClient(client.ClientConn)

	return client, nil
}

func (client Client) BroadcastDeltaFragment(deltaFragment *DeltaFragment) (*DeltaFragment, error) {
	var resp *DeltaFragment
	var err error

	serializedDeltaFragment := SerializeDeltaFragment(deltaFragment)
	for i := 0; i < client.Options.TimeoutRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), client.Options.Timeout+(client.Options.TimeoutBackoff*time.Duration(i)))
		defer cancel()

		resp, err = client.DarkNode.BroadcastDeltaFragment(ctx, &BroadcastDeltaFragmentRequest{
			From:          client.From,
			DeltaFragment: serializedDeltaFragment,
		}, grpc.FailFast(false))
		if err == nil {
			break
		}
		log.Println(err)
	}
	return resp, err
}
