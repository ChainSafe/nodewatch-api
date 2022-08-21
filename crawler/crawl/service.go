// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

// Package crawl holds the eth2 node discovery utilities
package crawl

import (
	"context"
	"eth2-crawler/store/peerstore"
	"eth2-crawler/store/record"
	"fmt"
	"net"

	"github.com/robfig/cron/v3"

	"eth2-crawler/crawler/p2p"
	ipResolver "eth2-crawler/resolver"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	ma "github.com/multiformats/go-multiaddr"
)

// listenConfig holds configuration for running v5discovry node
type listenConfig struct {
	bootNodeAddrs []string
	listenAddress net.IP
	listenPORT    int
	dbPath        string
}

// Initialize initializes the core crawler component
func Initialize(peerStore peerstore.Provider, historyStore record.Provider, ipResolver ipResolver.Provider, bootNodeAddrs []string) error {
	ctx := context.Background()

	listenCfg := &listenConfig{
		bootNodeAddrs: bootNodeAddrs,
		listenAddress: net.IPv4zero,
		listenPORT:    30304,
		dbPath:        "",
	}
	disc, err := startV5(listenCfg)
	if err != nil {
		return err
	}

	listenAddrs, err := multiAddressBuilder(listenCfg.listenAddress, listenCfg.listenPORT)
	if err != nil {
		return err
	}

	host, err := p2p.NewHost(
		libp2p.ListenAddrs(listenAddrs),
		libp2p.UserAgent("Eth2-Crawler"),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Security(noise.ID, noise.New),
		libp2p.NATPortMap(),
	)
	if err != nil {
		return err
	}

	c := newCrawler(disc, peerStore, historyStore, ipResolver, disc.RandomNodes(), host, 200)
	go c.start(ctx)
	// scheduler for updating peer
	go c.updatePeer(ctx)

	// add scheduler for updating history store
	scheduler := cron.New()
	_, err = scheduler.AddFunc("@daily", c.insertToHistory)
	if err != nil {
		return err
	}
	scheduler.Start()
	return nil
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
