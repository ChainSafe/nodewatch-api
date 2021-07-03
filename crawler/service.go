// Package crawler holds the whole crawler service. It includes crawler, db component and GraphQL
package crawler

import (
	"eth2-crawler/crawler/crawl"

	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/params"
)

// Start starts the crawler service
func Start() {
	h := log.CallerFileHandler(log.StdoutHandler)
	log.Root().SetHandler(h)

	handler := log.MultiHandler(
		log.LvlFilterHandler(log.LvlInfo, h),
	)
	log.Root().SetHandler(handler)

	err := crawl.Initialize(params.V5Bootnodes)
	panic(err)
}
