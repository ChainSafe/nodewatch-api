// Package crawl holds the eth2 node discovery utilities
package crawl

import (
	"context"
	"crypto/ecdsa"
	"eth2-crawler/crawler/p2p"

	"github.com/ethereum/go-ethereum/crypto"
)

type listenConfig struct {
	bootNodeAddrs []string
	listenAddress string
	dbPath        string
	privateKey    *ecdsa.PrivateKey
}

// Initialize initializes the core crawler component
func Initialize(bootNodeAddrs []string) error {
	ctx := context.Background()
	pkey, _ := crypto.GenerateKey()
	cfg := &listenConfig{
		bootNodeAddrs: bootNodeAddrs,
		listenAddress: "0.0.0.0:0",
		dbPath:        "",
		privateKey:    pkey,
	}
	return discv5Crawl(ctx, cfg)
}

func discv5Crawl(ctx context.Context, listenCfg *listenConfig) error {
	disc, err := startV5(listenCfg)
	if err != nil {
		return err
	}
	host, err := p2p.NewHost()
	if err != nil {
		return err
	}

	c := newCrawler(disc, listenCfg.privateKey, disc.RandomNodes(), host)
	c.run(ctx)
	return nil
}
