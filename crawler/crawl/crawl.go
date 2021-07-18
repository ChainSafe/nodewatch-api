package crawl

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"time"

	"eth2-crawler/crawler/p2p"
	"eth2-crawler/crawler/util"
	"eth2-crawler/models"
	ipResolver "eth2-crawler/resolver"
	"eth2-crawler/store"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

type crawler struct {
	disc       resolver
	peerStore  store.Provider
	ipResolver ipResolver.Provider
	iter       enode.Iterator
	nodeCh     chan *enode.Node
	privateKey *ecdsa.PrivateKey
	host       p2p.Host
}

// resolver holds methods of discovery v5
type resolver interface {
	Ping(n *enode.Node) error
}

// newCrawler inits new crawler service
func newCrawler(disc resolver, peerStore store.Provider, ipResolver ipResolver.Provider, privateKey *ecdsa.PrivateKey, iter enode.Iterator, host p2p.Host) *crawler {
	c := &crawler{
		disc:       disc,
		peerStore:  peerStore,
		ipResolver: ipResolver,
		privateKey: privateKey,
		iter:       iter,
		nodeCh:     make(chan *enode.Node),
		host:       host,
	}
	return c
}

// run runs the crawler
func (c *crawler) run(ctx context.Context) {
	doneCh := make(chan enode.Iterator)
	go c.runIterator(ctx, doneCh, c.iter)
	for {
		select {
		case n := <-c.nodeCh:
			c.collectNodeInfo(n)
		case <-doneCh:
			// crawling finished
			log.Info("finished iterator")
			fmt.Println("finished iterator")
			return
		}
	}
}

// runIterator uses the node iterator and sends node data through channel
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
	// only consider the node having tcp port exported
	if node.TCP() == 0 {
		return
	}
	// filter only eth2 nodes
	eth2Data, err := util.ParseEnrEth2Data(node)
	if err != nil { // not eth2 nodes
		return
	}
	log.Info("found a eth2 node", log.Ctx{"node": node})

	// get basic info
	peer, err := models.NewPeer(node, eth2Data)
	if err != nil {
		return
	}

	go c.collectNodeInfoRetryer(peer)
}

func (c *crawler) collectNodeInfoRetryer(peer *models.Peer) {
	count := 0
	var err error
	for count < 20 {
		time.Sleep(time.Second * 5)
		count++
		ctx := context.Background()
		err = c.host.IdentifyRequest(ctx, peer.GetPeerInfo())
		if err != nil {
			continue
		}
		var ag, pv string
		ag, err = c.host.GetAgentVersion(peer.ID)
		if err != nil {
			continue
		} else {
			peer.SetAgentVersion(ag)
		}
		pv, err = c.host.GetProtocolVersion(peer.ID)
		if err != nil {
			continue
		} else {
			peer.SetProtocolVersion(pv)
		}

		// successfully got all the node info's
		peer.SetConnectionStatus(true)

		// TODO: get the geo location details if we don't already have it in store
		geoLoc, err := c.ipResolver.GetGeoLocation(ctx, peer.IP)
		if err != nil {
			log.Error("unable to get geo information", log.Ctx{
				"error":   err,
				"ip_addr": peer.IP,
			})
			return
		}

		peer.SetGeoLocation(geoLoc)
		log.Info("successfully collected all info", peer.Log())
		err = c.peerStore.Upsert(context.TODO(), peer)
		if err != nil {
			log.Error("unable to write peer info to the store", log.Ctx{
				"error": err,
			})
		}
		return
	}
	// unsuccessful
	log.Error("failed on retryer", log.Ctx{
		"attempt": count,
		"error":   err,
	})
}
