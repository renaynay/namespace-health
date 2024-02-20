package feeder

import (
	"context"
	"fmt"
	"github.com/celestiaorg/celestia-node/blob"
	"github.com/celestiaorg/celestia-node/share"

	"github.com/celestiaorg/celestia-node/api/rpc/client"
	"github.com/celestiaorg/nmt/namespace"
)

type Feeder struct {
	node *client.Client
	nID  namespace.ID
}

func New(node *client.Client, nID namespace.ID) *Feeder {
	return &Feeder{
		node: node,
		nID:  nID,
	}
}

func (f *Feeder) Feed(ctx context.Context) error {
	namespace := share.Namespace(f.nID)

	fmt.Println("\n\n feeding the pooie............")

	b, err := blob.NewBlobV0(namespace, []byte("STUIE GRY????"))
	if err != nil {
		return err
	}

	height, err := f.node.Blob.Submit(ctx, []*blob.Blob{b}, -1)
	if err != nil {
		return err
	}

	fmt.Println("fed stuis at height: ", height)
	return nil
}
