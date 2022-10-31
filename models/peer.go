// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

// Package models represent the models for the service
package models

import (
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	"eth2-crawler/crawler/util"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"
	ic "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/protolambda/zrnt/eth2/beacon/common"
)

// 256 epochs
const blockIgnoreThreshold = 8192

// ClientName defines the type for eth2 client name
type ClientName string

const (
	PrysmClient      ClientName = "prysm"
	LighthouseClient ClientName = "lighthouse"
	TekuClient       ClientName = "teku"
	CortexClient     ClientName = "cortex"
	LodestarClient   ClientName = "lodestar"
	NimbusClient     ClientName = "nimbus"
	TrinityClient    ClientName = "trinity"
	GrandineClient   ClientName = "grandine"
	OthersClient     ClientName = "others"
)

// clients contains mapping for client name - possible identity names mapping
var clients map[ClientName][]string

func init() {
	clients = map[ClientName][]string{
		PrysmClient:      {"prysm"},
		LighthouseClient: {"lighthouse"},
		TekuClient:       {"teku"},
		CortexClient:     {"cortex"},
		LodestarClient:   {"lodestar", "js-libp2p"},
		NimbusClient:     {"nimbus"},
		TrinityClient:    {"trinity"},
		GrandineClient:   {"grandine", "rust"},
	}
}

// OS defines the type of os of agent

type OS string

const (
	OSLinux   OS = "linux"
	OSMAC     OS = "mac"
	OSWindows OS = "windows"
	OSUnknown OS = "unknown"

	VersionUnknown = "unknown"
	StatusSynced   = "synced"
	StatusUnsynced = "unsynced"
)

// UserAgent holds peer's client related info
type UserAgent struct {
	Name    ClientName `json:"name" bson:"name"`
	Version string     `json:"version" bson:"version"`
	OS      OS         `json:"os" bson:"os"`
}

// UsageType defines the ASN usage type
type UsageType string

const (
	UsageTypeNil            UsageType = ""
	UsageTypeHosting        UsageType = "hosting"
	UsageTypeResidential    UsageType = "residential"
	UsageTypeNonResidential UsageType = "non-residential"
	UsageTypeBusiness       UsageType = "business"
	UsageTypeEducation      UsageType = "education"
	UsageTypeGovernment     UsageType = "government"
	UsageTypeMilitary       UsageType = "military"
)

type Score int

const (
	ScoreGood Score = 3
	ScoreBad  Score = 0
)

// ASN holds the Autonomous system details
type ASN struct {
	ID     string    `json:"id" bson:"id"`
	Name   string    `json:"name" bson:"name"`
	Domain string    `json:"domain" bson:"domain"`
	Route  string    `json:"route" bson:"route"`
	Type   UsageType `json:"type" bson:"type"`
}

// GeoLocation holds peer's geo location related info
type GeoLocation struct {
	ASN       ASN     `json:"asn" nson:"asn"`
	Country   string  `json:"country_name" bson:"country"`
	State     string  `json:"state" bson:"state"`
	City      string  `json:"city" bson:"city"`
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

// Sync holds peer sync related info
type Sync struct {
	Status   bool `json:"status" bson:"status"`
	Distance int  `json:"distance" bson:"distance"` // sync distance in percentage
}

// String returns the sync status
func (s *Sync) String() string {
	if s.Status {
		return StatusSynced
	}
	return StatusUnsynced
}

// Peer holds all information of an eth2 peer
type Peer struct {
	ID     peer.ID `json:"id" bson:"_id"`
	NodeID string  `json:"node_id" bson:"node_id"`
	Pubkey string  `json:"pubkey" bson:"pubkey"`

	IP      string   `json:"ip" bson:"ip"`
	TCPPort int      `json:"tcp_port" bson:"tcp_port"`
	UDPPort int      `json:"udp_port" bson:"udp_port"`
	Addrs   []string `json:"addrs,omitempty" bson:"addrs"`

	Attnets common.AttnetBits `json:"enr_attnets,omitempty" bson:"attnets"`

	ForkDigest      common.ForkDigest `json:"fork_digest" bson:"fork_digest"`
	ForkDigestStr   string            `json:"fork_digest_str" bson:"fork_digest_str"`
	NextForkEpoch   Epoch             `json:"next_fork_epoch" bson:"next_fork_epoch"`
	NextForkVersion common.Version    `json:"next_fork_version" bson:"next_fork_version"`

	ProtocolVersion string       `json:"protocol_version,omitempty" bson:"protocol_version"`
	UserAgent       *UserAgent   `json:"user_agent,omitempty" bson:"user_agent"`
	UserAgentRaw    string       `json:"user_agent_raw" bson:"user_agent_raw"`
	GeoLocation     *GeoLocation `json:"geo_location" bson:"geo_location"`

	Sync  *Sync `json:"sync" bson:"sync"`
	Score Score `json:"score" bson:"score"`

	IsConnectable bool  `json:"is_connectable" bson:"is_connectable"`
	LastConnected int64 `json:"last_connected" bson:"last_connected"`
	LastUpdated   int64 `json:"last_updated" bson:"last_updated"`
}

// NewPeer initializes new peer
func NewPeer(node *enode.Node, eth2Data *common.Eth2Data) (*Peer, error) {
	pk := ic.PubKey((*ic.Secp256k1PublicKey)(node.Pubkey()))
	pkByte, err := pk.Raw()
	if err != nil {
		return nil, err
	}
	addr, err := util.AddrsFromEnode(node)
	if err != nil {
		return nil, err
	}
	addrStr := make([]string, 0)
	for _, madd := range addr.Addrs {
		addrStr = append(addrStr, madd.String())
	}

	attnetsVal := common.AttnetBits{}
	attnets, err := util.ParseEnrAttnets(node)
	if err == nil {
		attnetsVal = *attnets
	}
	return &Peer{
		ID:              addr.ID,
		NodeID:          node.ID().String(),
		Pubkey:          hex.EncodeToString(pkByte),
		IP:              node.IP().String(),
		TCPPort:         node.TCP(),
		UDPPort:         node.UDP(),
		Addrs:           addrStr,
		ForkDigest:      eth2Data.ForkDigest,
		ForkDigestStr:   eth2Data.ForkDigest.String(),
		NextForkVersion: eth2Data.NextForkVersion,
		NextForkEpoch:   Epoch(eth2Data.NextForkEpoch),
		Attnets:         attnetsVal,
		Score:           ScoreGood,
	}, nil
}

// SetProtocolVersion sets peer's protocol version
func (p *Peer) SetProtocolVersion(pv string) {
	p.ProtocolVersion = pv
}

// SetUserAgent sets peer's agent info
func (p *Peer) SetUserAgent(ag string) {
	// split the ag based on this format. might not be identical with each type of node
	// ag = Name/Version/OS(or git commit hash for Prysm)

	userAgent := new(UserAgent)
	parts := strings.Split(ag, "/")

nameChecker:
	for name, identityNames := range clients {
		for i := range identityNames {
			if strings.EqualFold(identityNames[i], parts[0]) {
				userAgent.Name = name
				break nameChecker
			}
		}
	}

	if userAgent.Name == "" {
		userAgent.Name = OthersClient
	}

	var os = ""
	switch userAgent.Name {
	case TekuClient:
		if len(parts) > 2 {
			userAgent.Version = parts[2]
		}
		if len(parts) > 3 {
			os = parts[3]
		}
	case PrysmClient:
		if len(parts) > 1 {
			userAgent.Version = parts[1]
		}
	default:
		if len(parts) > 1 {
			userAgent.Version = parts[1]
		}
		if len(parts) > 2 {
			os = parts[2]
		}
	}
	// update the version and os to standard form
	versions := strings.Split(userAgent.Version, "-")
	versions = strings.Split(versions[0], "+")
	userAgent.Version = versions[0]
	if userAgent.Version == "" {
		userAgent.Version = VersionUnknown
	}

	var validOS = []OS{OSLinux, OSMAC, OSWindows}
	for _, vos := range validOS {
		if strings.Contains(strings.ToLower(os), strings.ToLower(string(vos))) {
			userAgent.OS = vos
		}
	}
	if userAgent.OS == "" {
		userAgent.OS = OSUnknown
	}
	p.UserAgent = userAgent
	p.UserAgentRaw = ag
}

// SetConnectionStatus sets connection status and date
func (p *Peer) SetConnectionStatus(status bool) {
	p.IsConnectable = status
	if status {
		p.LastConnected = time.Now().Unix()
	}
}

// SetSyncStatus sets the sync status of a peer
func (p *Peer) SetSyncStatus(block int64) {
	cb := util.CurrentBlock()
	if cb-block <= blockIgnoreThreshold {
		p.Sync = &Sync{
			Status:   true,
			Distance: 0,
		}
	} else {
		p.Sync = &Sync{
			Status:   false,
			Distance: int(((cb - block) * 100) / cb),
		}
	}
}

// SetGeoLocation sets the geolocation information
func (p *Peer) SetGeoLocation(geoLocation *GeoLocation) {
	p.GeoLocation = geoLocation
}

// GetPeerInfo returns peer's AddrInfo
func (p *Peer) GetPeerInfo() *peer.AddrInfo {
	maddrs := make([]ma.Multiaddr, 0)
	for _, v := range p.Addrs {
		madd, _ := ma.NewMultiaddr(v)
		maddrs = append(maddrs, madd)
	}
	return &peer.AddrInfo{
		ID:    p.ID,
		Addrs: maddrs,
	}
}

// String returns peer object's json form in string
func (p *Peer) String() string {
	if p == nil {
		return "no data available"
	} else {
		dat, err := json.Marshal(p)
		if err != nil {
			return "failed to format peer data"
		}
		return string(dat)
	}
}

// Log returns log ctx from peer
func (p *Peer) Log() log.Ctx {
	dat, err := json.Marshal(p)
	if err != nil {
		return log.Ctx{}
	}
	val := log.Ctx{}
	err = json.Unmarshal(dat, &val)
	if err != nil {
		return log.Ctx{}
	}
	return val
}
