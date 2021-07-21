// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package crawl

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"strings"

	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/ethereum/go-ethereum/rlp"
)

// startV5 starts an ephemeral discovery v5 node.
func startV5(listenCfg *listenConfig) (*discover.UDPv5, error) {
	ln, config, err := getDiscoveryConfig(listenCfg)
	if err != nil {
		return nil, err
	}
	socket, err := listen(listenCfg)
	if err != nil {
		return nil, err
	}
	disc, err := discover.ListenV5(socket, ln, *config)
	if err != nil {
		return nil, err
	}
	return disc, nil
}

// getDiscoveryConfig returns config for listening v5 node for peer discovery
func getDiscoveryConfig(listenCfg *listenConfig) (*enode.LocalNode, *discover.Config, error) {
	cfg := new(discover.Config)

	cfg.PrivateKey = listenCfg.privateKey
	bootNodes, err := parseBootNodes(listenCfg.bootNodeAddrs)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing bootnodes: %w", err)
	}
	cfg.Bootnodes = bootNodes

	db, err := enode.OpenDB(listenCfg.dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening db: %w", err)
	}
	ln := enode.NewLocalNode(db, cfg.PrivateKey)
	return ln, cfg, nil
}

// listen opens an udp connections on given address
func listen(cfg *listenConfig) (*net.UDPConn, error) {
	udpAddr := &net.UDPAddr{
		IP:   cfg.listenAddress,
		Port: cfg.listenPORT,
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, fmt.Errorf("error listening to udp: %w", err)
	}
	return conn, nil
}

// parseBootNodes parse bootnodes from []string
func parseBootNodes(nodeStr []string) ([]*enode.Node, error) {
	nodes := make([]*enode.Node, len(nodeStr))
	var err error
	for i, record := range nodeStr {
		nodes[i], err = parseNode(record)
		if err != nil {
			return nil, fmt.Errorf("invalid bootstrap node: %w", err)
		}
	}
	return nodes, nil
}

// parseNode parses a node record and verifies its signature.
func parseNode(source string) (*enode.Node, error) {
	if strings.HasPrefix(source, "enode://") {
		return enode.ParseV4(source)
	}
	r, err := parseRecord(source)
	if err != nil {
		return nil, err
	}
	return enode.New(enode.ValidSchemes, r)
}

// parseRecord parses a node record from hex, base64, or raw binary input.
func parseRecord(source string) (*enr.Record, error) {
	bin := []byte(source)
	if d, ok := decodeRecordHex(bytes.TrimSpace(bin)); ok {
		bin = d
	} else if d, ok := decodeRecordBase64(bytes.TrimSpace(bin)); ok {
		bin = d
	}
	var r enr.Record
	err := rlp.DecodeBytes(bin, &r)
	return &r, err
}

func decodeRecordHex(b []byte) ([]byte, bool) {
	if bytes.HasPrefix(b, []byte("0x")) {
		b = b[2:]
	}
	dec := make([]byte, hex.DecodedLen(len(b)))
	_, err := hex.Decode(dec, b)
	return dec, err == nil
}

func decodeRecordBase64(b []byte) ([]byte, bool) {
	if bytes.HasPrefix(b, []byte("enr:")) {
		b = b[4:]
	}
	dec := make([]byte, base64.RawURLEncoding.DecodedLen(len(b)))
	n, err := base64.RawURLEncoding.Decode(dec, b)
	return dec[:n], err == nil
}
