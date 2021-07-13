// Package util holds utility functions for different conversion
package util

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/multiformats/go-multiaddr"

	"github.com/ethereum/go-ethereum/p2p/enode"
	beacon "github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/ztyp/codec"
)

func AddrsFromEnode(node *enode.Node) (*peer.AddrInfo, error) {
	madd, err := EnodeToMultiAddr(node)
	if err != nil {
		return nil, err
	}
	peerInfo, err := peer.AddrInfoFromP2pAddr(madd)
	if err != nil {
		return nil, err
	}
	return peerInfo, nil
}

func EnodeToMultiAddr(node *enode.Node) (multiaddr.Multiaddr, error) {
	ipScheme := "ip4"
	if len(node.IP()) == net.IPv6len {
		ipScheme = "ip6"
	}
	pubkey := node.Pubkey()
	peerID, err := peer.IDFromPublicKey(crypto.PubKey((*crypto.Secp256k1PublicKey)(pubkey)))
	if err != nil {
		return nil, err
	}
	multiAddrStr := fmt.Sprintf("/%s/%s/tcp/%d/p2p/%s", ipScheme, node.IP().String(), node.TCP(), peerID)
	multiAddr, err := multiaddr.NewMultiaddr(multiAddrStr)
	if err != nil {
		return nil, err
	}
	return multiAddr, nil
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
