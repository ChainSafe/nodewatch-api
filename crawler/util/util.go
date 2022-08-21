// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

// Package util holds utility functions for different conversion
package util

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	beacon "github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/ztyp/codec"
)

func AddrsFromEnode(node *enode.Node) (*peer.AddrInfo, error) {
	madds, err := EnodeToMultiAddr(node)
	if err != nil {
		return nil, err
	}

	if len(madds) == 0 {
		return nil, nil
	}

	peerInfo, err := peer.AddrInfoFromP2pAddr(madds[0])
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(madds); i++ {
		transport, _ := peer.SplitAddr(madds[i])
		peerInfo.Addrs = append(peerInfo.Addrs, transport)
	}

	return peerInfo, nil
}

func EnodeToMultiAddr(node *enode.Node) ([]multiaddr.Multiaddr, error) {
	multiAddrs := []multiaddr.Multiaddr{}

	ipScheme := "ip4"
	if len(node.IP()) == net.IPv6len {
		ipScheme = "ip6"
	}
	pubkey, err := crypto.ECDSAPublicKeyFromPubKey(*node.Pubkey())
	if err != nil {
		return nil, err
	}
	peerID, err := peer.IDFromPublicKey(pubkey)
	if err != nil {
		return nil, err
	}
	tcpMultiAddrStr := fmt.Sprintf("/%s/%s/tcp/%d/p2p/%s", ipScheme, node.IP().String(), node.TCP(), peerID)
	tcpMultiAddr, err := multiaddr.NewMultiaddr(tcpMultiAddrStr)
	if err != nil {
		return nil, err
	}
	multiAddrs = append(multiAddrs, tcpMultiAddr)

	udpMultiAddrStr := fmt.Sprintf("/%s/%s/udp/%d/p2p/%s", ipScheme, node.IP().String(), node.UDP(), peerID)
	udpMultiAddr, err := multiaddr.NewMultiaddr(udpMultiAddrStr)
	if err != nil {
		return nil, err
	}
	multiAddrs = append(multiAddrs, udpMultiAddr)

	return multiAddrs, nil
}

type Eth2ENREntry []byte

func (eee Eth2ENREntry) ENRKey() string {
	return "eth2"
}

func (eee Eth2ENREntry) Eth2Data() (*beacon.Eth2Data, error) {
	var dat beacon.Eth2Data
	if err := dat.Deserialize(codec.NewDecodingReader(bytes.NewReader(eee), uint64(len(eee)))); err != nil {
		return nil, err
	}
	return &dat, nil
}

func (eee Eth2ENREntry) String() string {
	dat, err := eee.Eth2Data()
	if err != nil {
		return fmt.Sprintf("invalid eth2 data! Raw: %x", eee[:])
	}
	return fmt.Sprintf("digest: %s, next fork version: %s, next fork epoch: %d",
		dat.ForkDigest, dat.NextForkVersion, dat.NextForkEpoch)
}

func ParseEnrEth2Data(n *enode.Node) (*beacon.Eth2Data, error) {
	var eth2 Eth2ENREntry
	if err := n.Load(&eth2); err != nil {
		return nil, err
	}
	dat, err := eth2.Eth2Data()
	if err != nil {
		return nil, fmt.Errorf("failed parsing eth2 bytes: %w", err)
	}
	return dat, nil
}

func ParseEnrAttnets(n *enode.Node) (*beacon.AttnetBits, error) {
	var attnets AttnetsENREntry
	if err := n.Load(&attnets); err != nil {
		return nil, err
	}
	dat, err := attnets.AttnetBits()
	if err != nil {
		return nil, fmt.Errorf("failed parsing attnets bytes: %w", err)
	}
	return &dat, nil
}

type AttnetsENREntry []byte

func (aee AttnetsENREntry) ENRKey() string {
	return "attnets"
}

func (aee AttnetsENREntry) AttnetBits() (beacon.AttnetBits, error) {
	var dat beacon.AttnetBits
	if err := dat.Deserialize(codec.NewDecodingReader(bytes.NewReader(aee), uint64(len(aee)))); err != nil {
		return beacon.AttnetBits{}, err
	}
	return dat, nil
}

func (aee AttnetsENREntry) String() string {
	return hex.EncodeToString(aee)
}

func getGenesisTime() time.Time {
	t, _ := time.Parse(time.RFC822, "01 Dec 20 12:00 GMT")
	return t
}

func CurrentBlock() int64 {
	duration := time.Since(getGenesisTime())
	return int64((duration / time.Second) / 12)
}
