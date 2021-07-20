// Package crawler holds the whole crawler service. It includes crawler, db component and GraphQL
package crawler

import (
	"eth2-crawler/crawler/crawl"
	ipResolver "eth2-crawler/resolver"
	"eth2-crawler/store"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

// Start starts the crawler service
func Start(peerStore store.Provider, ipResolver ipResolver.Provider) {
	h := log.CallerFileHandler(log.StdoutHandler)
	log.Root().SetHandler(h)

	handler := log.MultiHandler(
		log.LvlFilterHandler(log.LvlInfo, h),
	)
	log.Root().SetHandler(handler)

	err := crawl.Initialize(peerStore, ipResolver, params.V5Bootnodes)
	panic(err)
}
