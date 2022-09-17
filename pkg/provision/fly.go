package provision

import (
	"context"
	"errors"
	"fmt"
	graphqlclt "github.com/Khan/genqlient/graphql"
	"github.com/fly-apps/terraform-provider-fly/graphql"
	"github.com/fly-apps/terraform-provider-fly/pkg/apiv1"
	hreq "github.com/imroc/req/v3"
	cp "github.com/inlets/cloud-provision/provision"
	"github.com/jsiebens/inlets-on-fly/pkg/wg"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type FlyProvisioner struct {
	apiClient  *apiv1.MachineAPI
	graphqlClt graphqlclt.Client
	tunnel     *wg.Tunnel
	orgId      string
}

func NewFlyProvisioner(apiToken string, slug string, region string) (*FlyProvisioner, error) {
	h := http.Client{Timeout: 60 * time.Second, Transport: &Transport{UnderlyingTransport: http.DefaultTransport, Token: apiToken}}
	graphqlClt := graphqlclt.NewClient("https://api.fly.io/graphql", &h)

	orgId, err := getOrgId(context.Background(), graphqlClt, slug)
	if err != nil {
		return nil, err
	}

	tunnel, err := wg.Establish(context.Background(), orgId, region, apiToken, &graphqlClt)
	if err != nil {
		return nil, err
	}

	c := hreq.C()
	c.SetDial(tunnel.DialContext)
	c.SetCommonHeader("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	c.SetTimeout(2 * time.Minute)
	machineApiClient := apiv1.NewMachineAPI(c, "_api.internal:4280")

	return &FlyProvisioner{
		tunnel:     tunnel,
		apiClient:  machineApiClient,
		graphqlClt: graphqlClt,
		orgId:      orgId,
	}, nil
}

func (p *FlyProvisioner) Close() {
	if p.tunnel != nil {
		_ = p.tunnel.Close()
		_, _ = graphql.RemoveWireguardPeer(context.Background(), p.graphqlClt, graphql.RemoveWireGuardPeerInput{
			OrganizationId: p.tunnel.State.Org,
			Name:           p.tunnel.State.Name,
		})
	}
}

func (p *FlyProvisioner) Provision(host cp.BasicHost) (*cp.ProvisionedHost, error) {
	ctx := context.Background()

	tcp := host.Additional["inlets-tcp"] == "true"
	token := host.Additional["inlets-token"]
	version := host.Additional["inlets-version"]
	name := host.Name

	var mode = "http"
	var services = []apiv1.Service{
		{
			Ports: []apiv1.Port{
				{Port: 8123, Handlers: []string{}},
			},
			Protocol:     "tcp",
			InternalPort: 8123,
		},
	}

	if tcp {
		mode = "tcp"
		ports := strings.Split(host.Additional["inlets-ports"], ",")
		for _, p := range ports {
			i, err := strconv.ParseInt(p, 10, 64)
			if err != nil {
				return nil, err
			}
			services = append(services, apiv1.Service{
				Ports: []apiv1.Port{
					{Port: i, Handlers: []string{}},
				},
				Protocol:     "tcp",
				InternalPort: i,
			})
		}
	} else {
		services = append(services, apiv1.Service{
			Ports: []apiv1.Port{
				{Port: 80, Handlers: []string{"http"}},
				{Port: 443, Handlers: []string{"tls", "http"}},
			},
			Protocol:     "tcp",
			InternalPort: 8000,
		})
	}

	if _, err := graphql.CreateAppMutation(ctx, p.graphqlClt, name, p.orgId); err != nil {
		return nil, err
	}

	if _, err := graphql.AllocateIpAddress(ctx, p.graphqlClt, name, "global", graphql.IPAddressTypeV4); err != nil {
		return nil, err
	}
	if _, err := graphql.AllocateIpAddress(ctx, p.graphqlClt, name, "global", graphql.IPAddressTypeV6); err != nil {
		return nil, err
	}

	request := apiv1.MachineCreateOrUpdateRequest{
		Name:   name,
		Region: host.Region,
		Config: apiv1.MachineConfig{
			Image: fmt.Sprintf("ghcr.io/inlets/inlets-pro:%s", version),
			Init: apiv1.InitConfig{
				Entrypoint: []string{"inlets-pro"},
				Cmd: []string{
					mode, "server",
					"--auto-tls-san", fmt.Sprintf("%s.fly.dev", name),
					"--token", token,
				},
			},
			Services: services,
			Guest: apiv1.GuestConfig{
				Cpus:     1,
				MemoryMb: 256,
				CpuType:  "shared",
			},
		},
	}

	machine, err := p.createMachine(request, name)
	if err != nil {
		return nil, err
	}

	return &cp.ProvisionedHost{
		IP:     fmt.Sprintf("%s.fly.dev", name),
		ID:     fmt.Sprintf("%s/%s", name, machine.ID),
		Status: machine.State,
	}, nil
}

func (p *FlyProvisioner) Status(id string) (*cp.ProvisionedHost, error) {
	split := strings.Split(id, "/")
	machine, err := p.readMachine(split[0], split[1])
	if err != nil {
		return nil, err
	}

	status := machine.State
	if status == "started" {
		status = cp.ActiveStatus
	}

	return &cp.ProvisionedHost{
		IP:     fmt.Sprintf("%s.fly.dev", split[0]),
		ID:     id,
		Status: status,
	}, nil
}

func (p *FlyProvisioner) Delete(req cp.HostDeleteRequest) error {
	split := strings.Split(req.ID, "/")
	_, err := graphql.DeleteAppMutation(context.Background(), p.graphqlClt, split[0])
	return err
}

func getOrgId(ctx context.Context, clt graphqlclt.Client, slug string) (string, error) {
	if slug != "" {
		resp, err := graphql.Organization(ctx, clt, slug)
		if err != nil {
			return "", err
		}
		return resp.Organization.Id, nil
	}

	resp, err := graphql.OrgsQuery(ctx, clt)
	if err != nil {
		return "", err
	}

	if len(resp.Organizations.Nodes) > 1 {
		return "", errors.New("organization is ambiguous. Your account has more than one organization, you must specify which to use")
	}

	org := &resp.Organizations.Nodes[0]

	fmt.Println(org.GetId())

	return org.GetId(), nil
}

func (p *FlyProvisioner) createMachine(req apiv1.MachineCreateOrUpdateRequest, app string) (*apiv1.MachineResponse, error) {
	var response apiv1.MachineResponse
	if err := p.apiClient.CreateMachine(req, app, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (p *FlyProvisioner) readMachine(app string, id string) (*apiv1.MachineResponse, error) {
	var response apiv1.MachineResponse
	if _, err := p.apiClient.ReadMachine(app, id, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

type Transport struct {
	UnderlyingTransport http.RoundTripper
	Token               string
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+t.Token)
	return t.UnderlyingTransport.RoundTrip(req)
}
