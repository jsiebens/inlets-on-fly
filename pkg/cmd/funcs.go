package cmd

import (
	"fmt"
	"github.com/jsiebens/inlets-on-fly/pkg/fly"
	"github.com/jsiebens/inlets-on-fly/pkg/templates"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func parsePorts(values []string) ([]Ports, error) {
	var raw = map[int][]int{}

	for _, p := range values {
		parts := strings.Split(p, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid value for port: %s", p)
		}

		internalPort, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}

		externalPort, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		externalPorts, ok := raw[internalPort]
		if !ok {
			raw[internalPort] = []int{externalPort}
		} else {
			raw[internalPort] = append(externalPorts, externalPort)
		}
	}

	var ports = []Ports{}
	for i, e := range raw {
		ports = append(ports, Ports{InternalPort: i, ExternalPorts: e})
	}
	return ports, nil
}

func createApp(app *App) error {
	dname, err := os.MkdirTemp("", "inletsfly-")
	if err != nil {
		return err
	}

	defer os.RemoveAll(dname)

	fmt.Println("Temp dir name:", dname)

	flyFile, err := os.OpenFile(filepath.Join(dname, "fly.toml"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	dockerFile, err := os.OpenFile(filepath.Join(dname, "Dockerfile"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	if err := templates.Render(flyFile, "fly.toml.tpl", &app); err != nil {
		return err
	}

	if err := templates.Render(dockerFile, "Dockerfile.tpl", &app); err != nil {
		return err
	}

	if err := fly.CreateApp(dname, app.Name, app.Org); err != nil {
		return err
	}

	if app.Region != "" {
		if err := fly.SetRegion(dname, app.Region); err != nil {
			return err
		}
	}

	if err := fly.SetSecret(dname, app.Token); err != nil {
		return err
	}

	if err := fly.Deploy(dname); err != nil {
		return err
	}

	if err := templates.Render(os.Stdout, "message.tpl", &app); err != nil {
		return err
	}

	return nil
}
