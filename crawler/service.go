// Package crawler holds the whole crawler service. It includes crawler, db component and GraphQL
package crawler

import (
	"eth2-crawler/crawler/crawl"
	"eth2-crawler/store"

	"github.com/ethereum/go-ethereum/params"

	"github.com/ethereum/go-ethereum/log"
)

// Start starts the crawler service
func Start(peerStore store.Provider) {
	h := log.CallerFileHandler(log.StdoutHandler)
	log.Root().SetHandler(h)

	handler := log.MultiHandler(
		log.LvlFilterHandler(log.LvlInfo, h),
	)
	log.Root().SetHandler(handler)

	err := crawl.Initialize(params.V5Bootnodes, peerStore)
	panic(err)
}
