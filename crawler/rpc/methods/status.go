// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

// Package methods holds eth2 rpc methods
package methods

import (
	reqresp "eth2-crawler/crawler/rpc/request"

	beacon "github.com/protolambda/zrnt/eth2/beacon/common"
)

var StatusRPCv1 = reqresp.RPCMethod{
	Protocol:                  "/eth2/beacon_chain/req/status/1/ssz",
	RequestCodec:              reqresp.NewSSZCodec(func() reqresp.SerDes { return new(beacon.Status) }, beacon.StatusByteLen, beacon.StatusByteLen),
	ResponseChunkCodec:        reqresp.NewSSZCodec(func() reqresp.SerDes { return new(beacon.Status) }, beacon.StatusByteLen, beacon.StatusByteLen),
	DefaultResponseChunkCount: 1,
}
