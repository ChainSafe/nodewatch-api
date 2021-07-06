package crawl

import (
	"context"
	"crypto/ecdsa"
	"eth2-crawler/crawler/p2p"
	"fmt"

	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/p2p/enode"
)

type crawler struct {
	disc       resolver
	iter       enode.Iterator
	nodeCh     chan *enode.Node
	privateKey *ecdsa.PrivateKey
	host       p2p.Host
}

type resolver interface {
	RequestENR(*enode.Node) (*enode.Node, error)
}

func newCrawler(disc resolver, privateKey *ecdsa.PrivateKey, iter enode.Iterator, host p2p.Host) *crawler {
	c := &crawler{
		disc:       disc,
		privateKey: privateKey,
		iter:       iter,
		nodeCh:     make(chan *enode.Node),
		host:       host,
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

var ct = 0
var count = 0
var tcp = 0

func (c *crawler) collectNodeInfo(node *enode.Node) {
	fmt.Println("total:", ct)
	ct++
	// TODO check db. if node already exists update node. else insert as new node
	log.Info("found a node", log.Ctx{
		"node": node,
	})
	log.Info("retrieved enr successfully", log.Ctx{
		"node_id":   node.ID(),
		"node_ip ":  node.IP(),
		"node tcp ": node.TCP(),
		"node_udp":  node.UDP(),
	})
	// only consider the node having tcp port exported
	if node.TCP() == 0 {
		return
	} else {
		fmt.Println(node.String())
	}
	fmt.Println("tcp_count:", tcp)
	tcp++
	ctx, cf := context.WithCancel(context.Background())
	peerID, err := c.host.ConnectToPair(ctx, node.String())
	if err != nil {
		cf()
		fmt.Println(err)
		return
	} else {
		count++
		fmt.Println("connection successful: ", count)

	}

	err = c.host.IdentifyRequest(ctx, peerID)
	if err != nil {
		fmt.Println("error on identify", err)
		err = c.host.CloseConnection(peerID)
		if err != nil {
			fmt.Println("error closing connection", err)
		}
		return
	} else {
		fmt.Println("indentified")
		err = c.host.CloseConnection(peerID)
		if err != nil {
			fmt.Println("error closing connection", err)
		}
	}

	v, err := c.host.GetProtocolVersion(peerID)
	if err != nil {
		fmt.Println("error on p version", err)
		return
	} else {
		fmt.Println("v: ", v)
	}
	a, err := c.host.GetAgentVersion(peerID)
	if err != nil {
		fmt.Println("error on agent", err)
		return
	} else {
		fmt.Println("agent:", a)
	}

}
