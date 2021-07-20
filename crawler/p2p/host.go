// Package p2p represent p2p host service
package p2p

import (
	"context"
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/identify"
)

// Client represent custom p2p client
type Client struct {
	host.Host
	idSvc idService
}

// Host represent p2p services
type Host interface {
	IdentifyRequest(ctx context.Context, peerInfo *peer.AddrInfo) error
	GetProtocolVersion(peer.ID) (string, error)
	GetAgentVersion(peer.ID) (string, error)
}

type idService interface {
	IdentifyWait(c network.Conn) <-chan struct{}
}

// NewHost initializes custom host
func NewHost(opt ...libp2p.Option) (Host, error) {
	h, err := libp2p.New(context.Background(), opt...)
	if err != nil {
		return nil, err
	}
	idService, err := identify.NewIDService(h)
	if err != nil {
		return nil, err
	}
	return &Client{Host: h, idSvc: idService}, nil
}

// IdentifyRequest performs libp2p identify request after connecting to peer.
// It disconnects to peer after request is done
func (c *Client) IdentifyRequest(ctx context.Context, peerInfo *peer.AddrInfo) error {
	// Connect to peer first
	err := c.Connect(ctx, *peerInfo)
	if err != nil {
		return fmt.Errorf("error connecting to peer: %w", err)
	}
	defer func() {
		_ = c.Network().ClosePeer(peerInfo.ID)
	}()
	if conns := c.Network().ConnsToPeer(peerInfo.ID); len(conns) > 0 {
		select {
		case <-c.idSvc.IdentifyWait(conns[0]):
		case <-ctx.Done():
		}
	} else {
		return errors.New("not connected to peer, cannot await connection identify")
	}
	return nil
}

// GetProtocolVersion returns peer protocol version from peerstore.
// Need to call IdentifyRequest first for a peer.
func (c *Client) GetProtocolVersion(peerID peer.ID) (string, error) {
	key := "ProtocolVersion"
	value, err := c.Peerstore().Get(peerID, key)
	if err != nil {
		return "", fmt.Errorf("error getting protocol version:%w", err)
	}
	version, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("error converting interface to string")
	}
	return version, nil
}

// GetAgentVersion returns peer agent version  from peerstore.
// Need to call IdentifyRequest first for a peer.
func (c *Client) GetAgentVersion(peerID peer.ID) (string, error) {
	key := "AgentVersion"
	value, err := c.Peerstore().Get(peerID, key)
	if err != nil {
		return "", fmt.Errorf("error getting protocol version:%w", err)
	}
	version, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("error converting interface to string")
	}
	return version, nil
}
