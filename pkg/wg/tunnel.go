package wg

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	rawgql "github.com/Khan/genqlient/graphql"
	"github.com/fly-apps/terraform-provider-fly/graphql"
	"golang.org/x/crypto/curve25519"
	"math/rand"
	"net"
	"net/http"
	"net/netip"
	"strconv"
	"time"

	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

type Tunnel struct {
	token  string
	dev    *device.Device
	tun    tun.Device
	net    *netstack.Net
	dnsIP  net.IP
	State  *WireGuardState
	Config *Config

	resolv *net.Resolver
}

func c25519pair() (string, string) {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	var private [32]byte
	_, err := r.Read(private[:])
	if err != nil {
		panic(fmt.Sprintf("reading from random: %s", err))
	}

	public, err := curve25519.X25519(private[:], curve25519.Basepoint)
	if err != nil {
		panic(fmt.Sprintf("can't mult: %s", err))
	}

	return base64.StdEncoding.EncodeToString(public),
		base64.StdEncoding.EncodeToString(private[:])
}

func Establish(ctx context.Context, org string, region string, token string, client *rawgql.Client) (*Tunnel, error) {
	peerName := "inletsctl-" + strconv.FormatInt(time.Now().Unix(), 10)
	public, private := c25519pair()

	peer, err := graphql.AddWireguardPeer(ctx, *client, graphql.AddWireGuardPeerInput{
		OrganizationId: org,
		Region:         region,
		Name:           peerName,
		Pubkey:         public,
	})

	if err != nil {
		return nil, err
	}

	state := WireGuardState{
		LocalPrivate: private,
		Peer:         peer.AddWireGuardPeer,
		Token:        token,
		Org:          org,
		Name:         peerName,
	}

	tunnel, err := doConnect(&state)

	if err != nil {
		return nil, err
	}
	return tunnel, nil
}

func doConnect(state *WireGuardState) (*Tunnel, error) {
	cfg := state.TunnelConfig()
	addr, ok := netip.AddrFromSlice(cfg.LocalNetwork.IP)

	if !ok {
		return nil, fmt.Errorf("could not generate local network addr from IP %s: ", cfg.LocalNetwork.IP)
	}

	localIPs := []netip.Addr{addr}
	dnsIP, ok := netip.AddrFromSlice(cfg.DNS)

	if !ok {
		return nil, fmt.Errorf("could not generate DNS addr from IP %s: ", cfg.DNS)
	}

	mtu := cfg.MTU
	if mtu == 0 {
		mtu = device.DefaultMTU
	}

	tunDev, gNet, err := netstack.CreateNetTUN(localIPs, []netip.Addr{dnsIP}, mtu)
	if err != nil {
		return nil, err
	}

	endpointHost, endpointPort, err := net.SplitHostPort(cfg.Endpoint)
	if err != nil {
		return nil, err
	}

	endpointIPs, err := net.LookupIP(endpointHost)
	if err != nil {
		return nil, err
	}

	endpointIP := endpointIPs[rand.Intn(len(endpointIPs))]
	endpointAddr := net.JoinHostPort(endpointIP.String(), endpointPort)

	wgDev := device.NewDevice(tunDev, conn.NewDefaultBind(), device.NewLogger(cfg.LogLevel, "(fly-ssh) "))

	wgConf := bytes.NewBuffer(nil)
	fmt.Fprintf(wgConf, "private_key=%s\n", cfg.LocalPrivateKey.ToHex())
	fmt.Fprintf(wgConf, "public_key=%s\n", cfg.RemotePublicKey.ToHex())
	fmt.Fprintf(wgConf, "endpoint=%s\n", endpointAddr)
	fmt.Fprintf(wgConf, "allowed_ip=%s\n", cfg.RemoteNetwork)
	fmt.Fprintf(wgConf, "persistent_keepalive_interval=%d\n", cfg.KeepAlive)

	if err := wgDev.IpcSetOperation(bufio.NewReader(wgConf)); err != nil {
		return nil, err
	}
	wgDev.Up()

	return &Tunnel{
		dev:    wgDev,
		tun:    tunDev,
		net:    gNet,
		dnsIP:  cfg.DNS,
		Config: cfg,
		State:  state,

		resolv: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				fmt.Println("resolver.Dial", network, address)
				return gNet.DialContext(ctx, "tcp", net.JoinHostPort(dnsIP.String(), "53"))
			},
		},
	}, nil
}

func (t *Tunnel) Close() error {
	if t.dev != nil {
		t.dev.Close()
	}

	t.dev, t.net, t.tun = nil, nil, nil
	return nil
}

func (t *Tunnel) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	return t.net.DialContext(ctx, network, addr)
}

type Transport struct {
	Dial                *net.Dialer
	underlyingTransport http.RoundTripper
	token               string
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+t.token)
	return t.underlyingTransport.RoundTrip(req)
}

func (t *Tunnel) NewHttpClient() http.Client {
	underlyingTransport := http.Transport{
		DialContext: t.net.DialContext,
	}
	transport := Transport{
		token:               t.State.Token,
		underlyingTransport: &underlyingTransport,
	}
	return http.Client{Transport: &transport}
}
