package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/celestiaorg/celestia-node/share"

	"github.com/celestiaorg/celestia-node/api/rpc/client"

	"github.com/renaynay/namespace-health/feeder"
	"github.com/renaynay/namespace-health/reader"
	"github.com/renaynay/namespace-health/server"
)

var (
	nodeRPCAddr = "http://151.115.76.162:26658"
	token       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdfQ.TISQ2xg_qX_zTV3qHqIcD3WmO4XE3IUubIMrL6KwuZ8"

	namespaceID share.Namespace
)

func init() {
	nBytes, err := hex.DecodeString("cb02340a89c97f9b94c9")
	if err != nil {
		panic(err)
	}

	namespaceID, err = share.NewBlobNamespaceV0(nBytes)
	if err != nil {
		panic(err)
	}
}

func main() {
	ctx := context.Background() // TODO @renaynay: make this sig-termable

	fmt.Println("starting program....")

	// connect to celestia-node via RPC
	node, err := client.NewClient(ctx, nodeRPCAddr, token)
	if err != nil {
		// TODO @renaynay @ninabarbakadze: write to website endpoint the err
		fmt.Println("failed to connect to client", "err:  ", err)
		panic(err)
	}

	info, err := node.Node.Info(ctx)
	if err != nil {
		fmt.Println("failed to fetch node info", "err:  ", err)
		panic(err)
	}

	fmt.Println("node info", "info:   ", info.Type)

	r := reader.New(node, namespaceID.ToNMT())
	f := feeder.New(node, namespaceID.ToNMT())

	serv := server.New(r, f)
	go serv.Start(ctx)

	r.Read(ctx)
}
