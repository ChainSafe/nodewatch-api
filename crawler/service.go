// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

// Package crawler holds the whole crawler service. It includes crawler, db component and GraphQL
package crawler

import (
	"eth2-crawler/crawler/crawl"
	ipResolver "eth2-crawler/resolver"
	"eth2-crawler/store/peerstore"
	"eth2-crawler/store/record"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/protolambda/zrnt/eth2/beacon/common"
)

// Start starts the crawler service
func Start(
	peerStore peerstore.Provider,
	historyStore record.Provider,
	ipResolver ipResolver.Provider,
	allowedForkDigest map[common.ForkDigest]struct{},
) {
	h := log.CallerFileHandler(log.StdoutHandler)
	log.Root().SetHandler(h)

	handler := log.MultiHandler(
		log.LvlFilterHandler(log.LvlInfo, h),
	)
	log.Root().SetHandler(handler)

	err := crawl.Initialize(peerStore, historyStore, ipResolver, params.V5Bootnodes, allowedForkDigest)
	if err != nil {
		panic(err)
	}
}
