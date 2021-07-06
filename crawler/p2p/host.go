package p2p

import (
	"context"
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p/p2p/protocol/identify"

	"github.com/ethereum/go-ethereum/p2p/enode"

	"github.com/multiformats/go-multiaddr"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
)

type Client struct {
	host.Host
}

type Host interface {
	ConnectToPair(ctx context.Context, enr string) (peer.ID, error)
	CloseConnection(peerID peer.ID) error
	StorePair(enr string, peerId peer.ID) error
	GetPair(peerId peer.ID) (*enode.Node, error)
	IdentifyRequest(ctx context.Context, peerId peer.ID) error
	GetProtocolVersion(peerId peer.ID) (string, error)
	GetAgentVersion(peerId peer.ID) (string, error)
}

func NewHost() (Host, error) {
	h, err := libp2p.New(context.Background())
	if err != nil {
		return nil, err
	}
	return &Client{Host: h}, nil
}

func (c *Client) ConnectToPair(ctx context.Context, enr string) (peer.ID, error) {
	madd, err := getMultiAddr(enr)
	if err != nil {
		return "", fmt.Errorf("error constructing multiaddress from string: %v", err)
	}
	peerInfo, err := peer.AddrInfoFromP2pAddr(madd)
	if err != nil {
		return "", fmt.Errorf("error getting addressinfo from multiaddress: %v", err)
	}

	err = c.Connect(ctx, *peerInfo)
	if err != nil {
		return "", fmt.Errorf("error connecting to peer: %v", err)
	}
	return peerInfo.ID, nil
}

func (c *Client) CloseConnection(peerID peer.ID) error {
	return c.Host.Network().ClosePeer(peerID)
}

func (c *Client) StorePair(enr string, peerId peer.ID) error {
	key := "eth2-peers"
	return c.Peerstore().Put(peerId, key, enr)
}

func (c *Client) GetPair(peerId peer.ID) (*enode.Node, error) {
	key := "eth2-peers"
	value, err := c.Peerstore().Get(peerId, key)
	if err != nil {
		return nil, fmt.Errorf("error getting pair from pair-store:%v", err)
	}
	enr, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("error converting interface to string")
	}
	// get enode.Node from enr string
	node, err := parseEnrOrEnode(enr)
	if err != nil {
		return nil, fmt.Errorf("error parsing enode from string:%v", err)
	}
	return node, nil
}

func (c *Client) GetProtocolVersion(peerId peer.ID) (string, error) {
	key := "ProtocolVersion"
	value, err := c.Peerstore().Get(peerId, key)
	if err != nil {
		return "", fmt.Errorf("error getting protocal version:%v", err)
	}
	version, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("error converting interface to string")
	}
	return version, nil
}

func (c *Client) GetAgentVersion(peerId peer.ID) (string, error) {
	key := "AgentVersion"
	value, err := c.Peerstore().Get(peerId, key)
	if err != nil {
		return "", fmt.Errorf("error getting protocal version:%v", err)
	}
	version, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("error converting interface to string")
	}
	return version, nil
}

func (c *Client) IdentifyRequest(ctx context.Context, peerId peer.ID) error {
	idService, err := identify.NewIDService(c)
	if err != nil {
		return err
	}

	if conns := c.Network().ConnsToPeer(peerId); len(conns) > 0 {
		select {
		case <-idService.IdentifyWait(conns[0]):
			fmt.Println("completed identification")
		case <-ctx.Done():
			fmt.Println("canceled waiting for identification")
		}
	} else {
		return errors.New("not connected to peer, cannot await connection identify")
	}
	return nil
}

// construct addr info from enr string
func getMultiAddr(v string) (multiaddr.Multiaddr, error) {
	muAddr, err := multiaddr.NewMultiaddr(v)
	if err != nil {
		en, err := parseEnrOrEnode(v)
		if err != nil {
			return nil, err
		}
		muAddr, err = enodeToMultiAddr(en)
		if err != nil {
			return nil, err
		}
	}
	return muAddr, nil
}
