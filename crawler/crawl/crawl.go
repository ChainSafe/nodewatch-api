package crawl

import (
	"context"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/p2p/enode"
)

type crawler struct {
	disc       resolver
	iter       enode.Iterator
	nodeCh     chan *enode.Node
	privateKey *ecdsa.PrivateKey
}

type resolver interface {
	RequestENR(*enode.Node) (*enode.Node, error)
}

func newCrawler(disc resolver, privateKey *ecdsa.PrivateKey, iter enode.Iterator) *crawler {
	c := &crawler{
		disc:       disc,
		privateKey: privateKey,
		iter:       iter,
		nodeCh:     make(chan *enode.Node),
	}
	return c
}

func (c *crawler) run(ctx context.Context) {
	doneCh := make(chan enode.Iterator)
	go c.runIterator(ctx, doneCh, c.iter)
	for {
		select {
		case n := <-c.nodeCh:
			c.collectNodeInfo(n)
		case <-doneCh:
			// crawling finished
			return
		}
	}
}

func (c *crawler) runIterator(ctx context.Context, doneCh chan enode.Iterator, it enode.Iterator) {
	defer func() { doneCh <- it }()
	for it.Next() {
		select {
		case c.nodeCh <- it.Node():
		case <-ctx.Done():
			return
		}
	}
}

func (c *crawler) collectNodeInfo(node *enode.Node) {
	// TODO check db. if node already exists update node. else insert as new node
	log.Info("found a node")
	// Request the node record to check if the node is active and update the status accordingly.
	nn, err := c.disc.RequestENR(node)
	if err != nil {
		log.Warn("failed on Requesting enr", log.Ctx{"err": err})
		return
	}
	log.Info("retrieved enr successfully", log.Ctx{
		"node_id":   nn.ID(),
		"node_ip ":  nn.IP(),
		"node tcp ": nn.TCP(),
		"node_udp":  nn.UDP(),
	})

	// TODO : rlpxPing is not working as expected. This is just a placeholder for now.
	if nn.TCP() == 0 {
		return
	}
	// get additional info
	h, err := rlpxPing(c.privateKey, nn)
	if err != nil {
		log.Error("failed on rlpx ping", log.Ctx{"err": err})
		return
	}
	log.Info("rlpxPing: ", log.Ctx{"client_name": h.Name, "client_version": h.Version})
}
