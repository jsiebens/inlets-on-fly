package wg

import (
	"fmt"
	"github.com/fly-apps/terraform-provider-fly/graphql"
	"net"
)

type WireGuardState struct {
	Token        string
	Org          string
	Name         string
	LocalPrivate string                                                          `json:"localpublic"`
	Peer         graphql.AddWireguardPeerAddWireGuardPeerAddWireGuardPeerPayload `json:"peer"`
}

func (s *WireGuardState) TunnelConfig() *Config {
	skey := PrivateKey{}
	if err := skey.UnmarshalText([]byte(s.LocalPrivate)); err != nil {
		panic(fmt.Sprintf("martian local private key: %s", err))
	}

	pkey := PublicKey{}
	if err := pkey.UnmarshalText([]byte(s.Peer.Pubkey)); err != nil {
		panic(fmt.Sprintf("martian local public key: %s", err))
	}

	//fmt.Println(fmt.Sprintf("%s/120", s.Peer.Peerip))
	_, lnet, err := net.ParseCIDR(fmt.Sprintf("%s/120", s.Peer.Peerip))
	if err != nil {
		panic(fmt.Sprintf("martian local public: %s/120: %s", s.Peer.Peerip, err))
	}

	raddr := net.ParseIP(s.Peer.Peerip).To16()
	for i := 6; i < 16; i++ {
		raddr[i] = 0
	}

	_, rnet, _ := net.ParseCIDR(fmt.Sprintf("%s/48", raddr))

	raddr[15] = 3
	dns := net.ParseIP(raddr.String())

	wgl := *lnet
	wgr := *rnet

	return &Config{
		LocalPrivateKey: skey,
		LocalNetwork:    &wgl,
		RemotePublicKey: pkey,
		RemoteNetwork:   &wgr,
		Endpoint:        s.Peer.Endpointip + ":51820",
		DNS:             dns,
	}
}
