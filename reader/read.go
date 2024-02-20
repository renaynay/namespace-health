package reader

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	"github.com/celestiaorg/celestia-node/api/rpc/client"
	"github.com/celestiaorg/celestia-node/header"
	"github.com/celestiaorg/celestia-node/share"
	"github.com/celestiaorg/nmt/namespace"
)

const (
	// based on a 10s avg block time
	numBlocksInDay = numBlocksInOneMinute * 60 // TODO @renaynay: RETURN TO NORMAL VALUE

	numBlocksInOneMinute = 6
)

type Reader struct {
	node *client.Client
	nID  namespace.ID

	netHead            *header.ExtendedHeader
	lastRollingAverage float64
	networkHeight      uint64

	isLastAverageHealthier atomic.Value
}

func New(node *client.Client, nID namespace.ID) *Reader {
	r := &Reader{
		node: node,
		nID:  nID,
	}
	r.isLastAverageHealthier.Store(&Health{})
	return r
}

func (r *Reader) Read(ctx context.Context) {
	fmt.Println("initialising average...")
	r.initializeAverage(ctx)

	for {
		select {
		case <-ctx.Done():
			panic("CONTEXT CANCELLED!!!")
		default:
			// start a 1m rolling average // TODO @renaynay: MAKE IT 1 MIN!!!!
			time.Sleep(1 * time.Minute)

			fmt.Println("\n\n\nchecking health......")

			netHead, err := r.node.Header.NetworkHead(ctx)
			if err != nil {
				fmt.Println("failed to fetch network head", "err:  ", err)
				return
			}
			r.netHead = netHead

			to := numBlocksInOneMinute + r.netHead.Height()
			fmt.Println("fetching header range from  ", r.netHead.Height(), "  to: ", to)
			headerRange, err := r.node.Header.GetRangeByHeight(ctx, r.netHead, to)
			if err != nil {
				fmt.Println("failed to fetch header range", "err:  ", err)
				return
			}

			total, err := r.getTotalSharesForRange(ctx, headerRange)
			if err != nil {
				fmt.Println("failed to fetch total bytes for header range", "err:  ", err)
				return
			}

			avg := estimateNamespaceHealth(total, len(headerRange))
			fmt.Println("new rolling average: ", avg)

			isHealthier := avg > r.lastRollingAverage
			scale := scaleDistance(r.lastRollingAverage, avg)
			r.isLastAverageHealthier.Store(&Health{
				IsHealthier: isHealthier,
				Scale:       scale,
			})
			r.lastRollingAverage = avg
		}
	}
}

func (r *Reader) Health() *Health {
	h := r.isLastAverageHealthier.Load()
	return h.(*Health)
}

func (r *Reader) initializeAverage(ctx context.Context) {
	// get network height
	netHead, err := r.node.Header.NetworkHead(ctx)
	if err != nil {
		fmt.Println("failed to fetch network head", "err:  ", err)
		return
	}
	currentHeight := netHead.Height()
	fmt.Println("got current height:   ", currentHeight)

	r.netHead = netHead

	// get headers for past day
	startHeight := currentHeight - numBlocksInDay
	start, err := r.node.Header.GetByHeight(ctx, startHeight)
	if err != nil {
		fmt.Println("failed to fetch header for height:   ", startHeight, "  err:  ", err)
		return
	}
	headers, err := r.node.Header.GetRangeByHeight(ctx, start, currentHeight)
	if err != nil {
		fmt.Println("failed to fetch header range:   ", startHeight, "  to   ", currentHeight, "  err:  ", err)
		return
	}

	fmt.Println("got header range for past 24 hours: ", len(headers))

	total, err := r.getTotalSharesForRange(ctx, headers)
	if err != nil {
		fmt.Println("failed to get total shares for range", "err:  ", err)
		return
	}

	// calculate bytes/block avg
	r.lastRollingAverage = estimateNamespaceHealth(total, len(headers))
	r.isLastAverageHealthier.Store(&Health{
		IsHealthier: false,
		Scale:       3, // starting point is neutral
	})
	fmt.Println("INITIALISED AVERAGE:   ", r.lastRollingAverage)
}

func (r *Reader) getTotalSharesForRange(
	ctx context.Context,
	headerRange []*header.ExtendedHeader,
) (uint64, error) {
	total := uint64(0)

	for _, h := range headerRange {
		shrs, err := r.node.Share.GetSharesByNamespace(ctx, h, share.Namespace(r.nID))
		if err != nil {
			fmt.Println("failed to get shares for namespace, at height:  ", h.Height(), "  err:  ", err)
			return 0, err
		}
		numShares := len(shrs.Flatten())
		// TODO @renaynay: remove bytes from NID? shouldn't matter bc it should be constant
		total += uint64(numShares * appconsts.ShareSize)
	}

	return total, nil
}

func estimateNamespaceHealth(totalBytes uint64, numBlocks int) float64 {
	return float64(totalBytes) / float64(numBlocks)
}
