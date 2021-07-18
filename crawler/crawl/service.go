// Package crawl holds the eth2 node discovery utilities
package crawl

import (
	"context"
	"crypto/ecdsa"
	"eth2-crawler/crawler/p2p"
	"eth2-crawler/store"
	"fmt"
	"net"

	noise "github.com/libp2p/go-libp2p-noise"

	ma "github.com/multiformats/go-multiaddr"

	"github.com/libp2p/go-tcp-transport"

	ic "github.com/libp2p/go-libp2p-core/crypto"

	"github.com/libp2p/go-libp2p"

	"github.com/ethereum/go-ethereum/crypto"
)

// listenConfig holds configuration for running v5discovry node
type listenConfig struct {
	bootNodeAddrs []string
	listenAddress net.IP
	listenPORT    int
	dbPath        string
	privateKey    *ecdsa.PrivateKey
}

// Initialize initializes the core crawler component
func Initialize(bootNodeAddrs []string, peerStore store.Provider) error {
	ctx := context.Background()
	pkey, _ := crypto.GenerateKey()
	cfg := &listenConfig{
		bootNodeAddrs: bootNodeAddrs,
		listenAddress: net.IPv4zero,
		listenPORT:    30304,
		dbPath:        "",
		privateKey:    pkey,
	}
	return discv5Crawl(ctx, cfg, peerStore)
}

// discv5Crawl start the crawler
func discv5Crawl(ctx context.Context, listenCfg *listenConfig, peerStore store.Provider) error {
	disc, err := startV5(listenCfg)
	if err != nil {
		return err
	}

	listenAddrs, err := multiAddressBuilder(listenCfg.listenAddress, listenCfg.listenPORT)
	if err != nil {
		return err
	}
	host, err := p2p.NewHost(
		libp2p.Identity(convertToInterfacePrivkey(listenCfg.privateKey)),
		libp2p.ListenAddrs(listenAddrs),
		libp2p.UserAgent("Eth2-Crawler"),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Security(noise.ID, noise.New),
		libp2p.NATPortMap(),
	)
	if err != nil {
		return err
	}

	c := newCrawler(disc, listenCfg.privateKey, disc.RandomNodes(), host, peerStore)
	c.run(ctx)
	return nil
}

func convertToInterfacePrivkey(privkey *ecdsa.PrivateKey) ic.PrivKey {
	typeAssertedKey := ic.PrivKey((*ic.Secp256k1PrivateKey)(privkey))
	return typeAssertedKey
}

func multiAddressBuilder(ipAddr net.IP, port int) (ma.Multiaddr, error) {
	if ipAddr.To4() == nil && ipAddr.To16() == nil {
		return nil, fmt.Errorf("invalid ip address provided: %s", ipAddr)
	}
	if ipAddr.To4() != nil {
		return ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", ipAddr.String(), port))
	}
	return ma.NewMultiaddr(fmt.Sprintf("/ip6/%s/tcp/%d", ipAddr.String(), port))
}
